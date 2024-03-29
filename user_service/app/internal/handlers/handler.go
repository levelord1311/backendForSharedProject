package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/user_service/internal/apperror"
	"github.com/levelord1311/backendForSharedProject/user_service/internal/models"
	"net/http"
	"strconv"
)

const (
	usersURL      = "/api/users"
	singleUserURL = "/api/users/:id"
	authURL       = "/api/users/auth"
)

type Service interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, dto *models.CreateUserDTO) (uint, error)
	SignIn(ctx context.Context, dto *models.SignInUserDTO) (*models.User, error)
	//UpdatePassword(ctx context.Context, dto *models.UpdateUserDTO) error
	//Delete(ctx context.Context, id string) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, singleUserURL, h.GetUser)
	router.HandlerFunc(http.MethodPost, usersURL, h.CreateUser)
	router.HandlerFunc(http.MethodPost, authURL, h.SignIn)
	//router.HandlerFunc(http.MethodPatch, singleUserURL, h.PartiallyUpdateUser)
	//router.HandlerFunc(http.MethodDelete, singleUserURL, h.DeleteUser)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	userID, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		switch err {
		case apperror.ErrNotFound:
			writeError(w, err, http.StatusNotFound)
			return
		default:
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var crUser *models.CreateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	userID, err := h.service.Create(r.Context(), crUser)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%d", usersURL, userID))
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto *models.SignInUserDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := h.service.SignIn(r.Context(), dto)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)

}

//func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) {
//	h.Logger.Info("PARTIALLY UPDATE USER")
//	w.Header().Set("Content-Type", "application/json")
//
//	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
//	userID, err := strconv.Atoi(params.ByName("id"))
//	if err != nil {
//		return err
//	}
//
//	h.Logger.Debug("decode update user dto")
//	var updUser *models.UpdateUserDTO
//	defer r.Body.Close()
//	if err := json.NewDecoder(r.Body).Decode(&updUser); err != nil {
//		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
//	}
//	updUser.ID = uint(userID)
//
//	err = h.service.UpdatePassword(r.Context(), updUser)
//	if err != nil {
//		return err
//	}
//	w.WriteHeader(http.StatusNoContent)
//
//	return nil
//}

// TODO исправить delete - вместо полноценного удаления из БД вешать признак "пометка на удаление"

//func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
//	h.Logger.Info("DELETE USER")
//	w.Header().Set("Content-Type", "application/json")
//
//	h.Logger.Debug("get id from context")
//	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
//	userID := params.ByName("id")
//
//	err := h.service.Delete(r.Context(), userID)
//	if err != nil {
//		return err
//	}
//	w.WriteHeader(http.StatusNoContent)
//
//	return nil
//}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
	return
}
