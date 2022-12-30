package auth

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/apperror"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/client/user_service"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/jwt"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"net/http"
)

const (
	authURL   = "/api/auth"
	signupURL = "/api/signup"
)

type Handler struct {
	Logger      logging.Logger
	UserService user_service.UserService
	JWTHelper   jwt.Helper
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, signupURL, apperror.Middleware(h.SignUp))
	router.HandlerFunc(http.MethodPost, authURL, apperror.Middleware(h.SignIn))
}

// SignUp godoc
//
//	@Summary Create user
//	@Description Creates User & returns JWT
//	@Tags user
//	@Accept json
//	@Param DTO body user_service.CreateUserDTO true "user data"
//	@Produce json
//	@Success 201 {string} string "jwt.token.string"
//	@Failure 400 {object}	apperror.AppError
//	@Failure 418 {object}	apperror.AppError
//	@Router /signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	var dto *user_service.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("failed to decode data", "")
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

// SignIn godoc
//
//	@Summary Authenticate user
//	@Description authenticates user and returns JWT
//	@Tags user
//	@Accept json
//	@Param DTO body user_service.SignInUserDTO true "user data"
//	@Produce json
//	@Success 200 {string} jwt.token.string
//	@Failure 400 {object}	apperror.AppError
//	@Failure 418 {object}	apperror.AppError
//	@Router /auth [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {

	var token []byte

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		var dto *user_service.SignInUserDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return apperror.BadRequestError("failed to decode data", "")
		}
		u, err := h.UserService.SignIn(r.Context(), dto)
		if err != nil {
			return err
		}
		h.Logger.Debugf("user:%v", u)
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

	w.WriteHeader(http.StatusOK)
	w.Write(token)

	return nil
}
