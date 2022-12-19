package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/user_service/internal/user"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/user_service/pkg/logging"
	"net/http"
	"strconv"
)

const (
	authURL       = "/api/users/auth"
	usersURL      = "/api/users"
	singleUserURL = "/api/users/:id"
)

type Handler struct {
	Logger      logging.Logger
	UserService user.Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, authURL, apperror.Middleware(h.SignIn))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, singleUserURL, apperror.Middleware(h.GetUser))
	router.HandlerFunc(http.MethodPatch, singleUserURL, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, singleUserURL, apperror.Middleware(h.DeleteUser))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("getting id from context..")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")

	user, err := h.UserService.GetByID(r.Context(), userID)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshalling user..")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decoding data to create user dto..")
	var crUser *user.CreateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	defer r.Body.Close()

	userID, err := h.UserService.Create(r.Context(), crUser)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%d", usersURL, userID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("SIGN IN USER WITH LOGIN AND PASSWORD")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decoding data to create user dto..")
	var signUser *user.SignInUserDTO

	if err := json.NewDecoder(r.Body).Decode(&signUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	defer r.Body.Close()

	h.Logger.Debugf("decoded DTO:%v", signUser)
	u, err := h.UserService.SignIn(r.Context(), signUser)
	if err != nil {
		return err
	}
	h.Logger.Debugf("recieved user:%v", u)
	h.Logger.Debug("marshalling user..")
	userBytes, err := json.Marshal(&u)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)

	return nil
}

func (h *Handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE USER")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		return err
	}

	h.Logger.Debug("decode update user dto")
	var updUser *user.UpdateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	updUser.ID = uint(userID)

	err = h.UserService.UpdatePassword(r.Context(), updUser)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// TODO исправить delete - вместо полноценного удаления из БД вешать признак "пометка на удаление"

//func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
//	h.Logger.Info("DELETE USER")
//	w.Header().Set("Content-Type", "application/json")
//
//	h.Logger.Debug("get id from context")
//	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
//	userID := params.ByName("id")
//
//	err := h.UserService.Delete(r.Context(), userID)
//	if err != nil {
//		return err
//	}
//	w.WriteHeader(http.StatusNoContent)
//
//	return nil
//}
