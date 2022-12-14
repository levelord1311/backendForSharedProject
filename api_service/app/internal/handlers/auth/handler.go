package auth

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	user_service2 "github.com/levelord1311/backendForSharedProject/api_service/internal/client/user_service"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/jwt"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"net/http"
)

const (
	authURL   = "/api/auth"
	signupURL = "/api/signup"
)

type Handler struct {
	Logger      logging.Logger
	UserService user_service2.UserService
	JWTHelper   jwt.Helper
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, signupURL, apperror.Middleware(h.SignUp))
	router.HandlerFunc(http.MethodPost, authURL, apperror.Middleware(h.SignIn))
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	var dto *user_service2.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("failed to decode data")
	}

	u, err := h.UserService.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	token, err := h.JWTHelper.GenerateAccessToken(u)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(token)

	return nil
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {

	var token []byte
	//var err error

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		dto := &user_service2.SignInUserDTO{}
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return apperror.BadRequestError("failed to decode data")
		}
		u, err := h.UserService.SignIn(r.Context(), dto.Login, dto.Password)
		if err != nil {
			return err
		}
		token, err = h.JWTHelper.GenerateAccessToken(u)
		if err != nil {
			return err
		}
		//case http.MethodPut:
		//	defer r.Body.Close()
		//	var rt jwt.RT
		//	if err := json.NewDecoder(r.Body).Decode(&rt); err != nil {
		//		apperror.BadRequestError("failed to decode data")
		//	}
		//	token, err = h.JWTHelper.UpdateRefreshToken(rt)
		//	if err != nil {
		//		return err
		//	}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(token)

	return nil
}
