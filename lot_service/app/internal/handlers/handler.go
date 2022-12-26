package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/apperror"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/lot"
	"github.com/levelord1311/backendForSharedProject/lot_service/internal/lot/service"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/api/sort"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/logging"
	"net/http"
	"strconv"
)

const (
	lotsURL      = "/api/lots"
	lotsOfUser   = "/api/lots/user/:id"
	singleLotURL = "/api/lots/lot/:id"
)

type Handler struct {
	Logger     logging.Logger
	LotService service.Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, lotsURL, apperror.Middleware(h.CreateLot))
	router.HandlerFunc(http.MethodGet, singleLotURL, apperror.Middleware(h.GetLot))
	router.HandlerFunc(http.MethodGet, lotsURL, sort.Middleware(apperror.Middleware(h.GetLots)))
	router.HandlerFunc(http.MethodGet, lotsOfUser, apperror.Middleware(h.GetLotsByUser))
	router.HandlerFunc(http.MethodPatch, singleLotURL, apperror.Middleware(h.UpdateLotPrice))
	//	router.HandlerFunc(http.MethodDelete, singleLotURL, apperror.Middleware(h.DeleteLot))
}

func (h *Handler) CreateLot(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE LOT")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decoding r.body into create lot dto..")
	dto := &lot.CreateLotDTO{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		return apperror.BadRequestError("invalid data", "")
	}

	lotID, err := h.LotService.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%d", lotsURL, lotID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) GetLot(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET LOT")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("getting id from context..")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	lotID := params.ByName("id")

	l, err := h.LotService.GetByLotID(r.Context(), lotID)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshalling lot..")
	lotBytes, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("failed to marshall lot. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lotBytes)

	return nil
}

func (h *Handler) GetLotsByUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET LOTS BY USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("getting id from context..")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")

	lots, err := h.LotService.GetByUserID(r.Context(), userID)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshalling lots..")
	lotsBytes, err := json.Marshal(lots)
	if err != nil {
		return fmt.Errorf("failed to marshall lots. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lotsBytes)
	return nil
}

func (h *Handler) GetLots(w http.ResponseWriter, r *http.Request) error {

	h.Logger.Info("GET LOTS")
	w.Header().Set("Content-Type", "application/json")

	lots, err := h.LotService.GetLotsWithFilter(r.Context(), r.URL.Query())
	if err != nil {
		return err
	}

	h.Logger.Debug("marshalling lots..")
	lotsBytes, err := json.Marshal(lots)
	if err != nil {
		return fmt.Errorf("failed to marshall lots. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lotsBytes)

	return nil
}

func (h *Handler) UpdateLotPrice(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("UPDATE LOT PRICE")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("getting id from context..")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	lotID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	h.Logger.Debug("decoding r.body into update lot dto..")
	var dto *lot.UpdateLotDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid data", "")
	}

	dto.ID = uint(lotID)

	err = h.LotService.Update(r.Context(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// TODO исправить delete - вместо полноценного удаления из БД вешать признак "пометка на удаление"

//func (h *Handler) DeleteLot(w http.ResponseWriter, r *http.Request) error {
//	h.Logger.Info("DELETE LOT")
//	w.Header().Set("Content-Type", "application/json")
//
//	h.Logger.Debug("getting lot_id from URL..")
//	lotIDStr := r.URL.Query().Get("lot_id")
//	if lotIDStr == "" {
//		return apperror.BadRequestError("lot_id query parameter is required and must be an unsigned integer")
//	}
//
//	lotID, err := strconv.Atoi(lotIDStr)
//	if err != nil {
//		return err
//	}
//	h.Logger.Debug("searching for `user_id` in header...")
//	if r.Header["user_id"] == nil {
//		err = errors.New(" `user_id` field has not been found in header")
//		return err
//	}
//
//	userID, err := strconv.Atoi(r.Header["user_id"][0])
//	if err != nil {
//		return err
//	}
//
//	err = h.LotService.Delete(r.Context(), uint(lotID), uint(userID))
//	if err != nil {
//		return err
//	}
//	w.WriteHeader(http.StatusNoContent)
//
//	return nil
//}
