package lots

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/apperror"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/client/lot_service"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/jwt"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"net/http"
	"strconv"
	"time"
)

const (
	lotsURL      = "/api/lots"
	lotsOfUser   = "/api/lots/user/:id"
	singleLotURL = "/api/lots/lot/:id"
	weekURL      = "/api/lots/week"
)

type Handler struct {
	Logger     logging.Logger
	LotService lot_service.LotService
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, lotsOfUser, apperror.Middleware(h.GetByUserID))
	router.HandlerFunc(http.MethodGet, lotsURL, apperror.Middleware(h.GetLots))
	router.HandlerFunc(http.MethodPost, lotsURL, jwt.Middleware(apperror.Middleware(h.CreateLot)))
	router.HandlerFunc(http.MethodGet, singleLotURL, apperror.Middleware(h.GetByLotID))
	router.HandlerFunc(http.MethodPatch, singleLotURL, jwt.Middleware(apperror.Middleware(h.UpdateLot)))
	router.HandlerFunc(http.MethodDelete, singleLotURL, jwt.Middleware(apperror.Middleware(h.DeleteLot)))
	router.HandlerFunc(http.MethodGet, weekURL, apperror.Middleware(h.GetLastWeek))
}

// GetLots godoc
//
//	@Summary		Show lots
//	@Description	Get lots with filter from query.
//	@Description	Supported comparisons: eq, neq, lt, lte, gt, gte.
//	@Description	For range use example ?created_by=2022-12-21:2022-12-22
//	@Tags			lots
//	@Produce		json
//	@Param 			estate_type query string false "filter by estate type"
//	@Param 			rooms query string false "filter by rooms quantity"
//	@Param 			district query string false "filter by district"
//	@Param 			price query string false "filter by price"
//	@Param 			created_at query string false "filter by date of creation"
//	@Param 			floor query string false "filter by floor"
//	@Success		200	{object}	lot_service.Lot
//	@Failure		400	{object}	apperror.AppError
//	@Failure		404	{object}	apperror.AppError
//	@Failure		418	{object}	apperror.AppError
//	@Router			/lots [get]
func (h *Handler) GetLots(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("getting raw query from url..")
	rQuery := r.URL.RawQuery
	lots, err := h.LotService.GetWithFilter(r.Context(), rQuery)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lots)

	return nil
}

// GetByLotID godoc
//
//	@Summary		Show lot
//	@Description	get lot by its ID
//	@Tags			lots
//	@Produce		json
//	@Param			id	path		int	true	"Lot ID"
//	@Success		200	{object}	lot_service.Lot
//	@Failure		400	{object}	apperror.AppError
//	@Failure		404	{object}	apperror.AppError
//	@Failure		418	{object}	apperror.AppError
//	@Router			/lots/lot/{id} [get]
func (h *Handler) GetByLotID(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("getting id from context..")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	lotID := params.ByName("id")

	_, err := strconv.Atoi(lotID)
	if err != nil {
		return err
	}

	lot, err := h.LotService.GetByLotID(r.Context(), lotID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lot)

	return nil
}

// GetByUserID godoc
//
//	@Summary		Show lots
//	@Description	get lots created by user
//	@Tags			lots
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{array}		lot_service.Lot
//	@Failure		400	{object}	apperror.AppError
//	@Failure		404	{object}	apperror.AppError
//	@Failure		418	{object}	apperror.AppError
//	@Router			/lots/user/{id} [get]
func (h *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("getting id from context..")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")

	_, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	lots, err := h.LotService.GetByUserID(r.Context(), userID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lots)

	return nil
}

// CreateLot godoc
//
//	@Summary		Create lot
//	@Description	creates lot by user id from JWT
//	@Tags			lots
//	@Accept			json
//	@Produce		json
//	@Param			Token	header		string	true	"JWT token"
//	@Success		201
//	@Header			201 {string} Location "/lots/lot/{created_id}"
//	@Failure		400	{object}	apperror.AppError
//	@Failure		404	{object}	apperror.AppError
//	@Failure		418	{object}	apperror.AppError
//	@Router			/lots [post]
func (h *Handler) CreateLot(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	dto := &lot_service.CreateLotDTO{}
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		return apperror.BadRequestError("failed to decode data", "")
	}

	h.Logger.Info("getting user_id from req.context()..")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		return fmt.Errorf("error with type of req.context value of key 'user_id'")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	dto.CreatedByUserID = uint(id)

	lotID, err := h.LotService.Create(r.Context(), dto)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("%s/%d", singleLotURL, lotID))
	return nil
}

// UpdateLot godoc
//
//	@Summary		Update lot price
//	@Description	Get lots created during last 7 days.
//	@Tags			lots
//	@Accept 		json
//	@Param			price	body		int	true	"new lot price"
//	@Success		204
//	@Failure		400	{object}	apperror.AppError
//	@Failure		404	{object}	apperror.AppError
//	@Failure		418	{object}	apperror.AppError
//	@Router			/lots/lot/{id} [patch]
func (h *Handler) UpdateLot(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("getting id from context..")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	ID := params.ByName("id")

	lotID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	var dto *lot_service.UpdateLotDTO
	if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("failed to decode data", "")
	}

	dto.ID = uint(lotID)

	h.Logger.Info("getting user_id from req.context()")
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		return fmt.Errorf("error with type of req.context value of key 'user_id', must be an unsigned integer")
	}

	dto.CreatedByUserID = userID

	err = h.LotService.Update(r.Context(), dto)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) DeleteLot(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("getting lot_id from query..")
	lotID := r.URL.Query().Get("lot_id")
	if lotID == "" {
		return apperror.BadRequestError("lot_id query parameter is required and must be an unsigned integer", "")
	}

	_, err := strconv.Atoi(lotID)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	h.Logger.Info("getting user_id from req.context()")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		return fmt.Errorf("error with type of req.context value of key 'user_id'")
	}

	err = h.LotService.Delete(r.Context(), lotID, userID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

// GetLastWeek godoc
//
//	@Summary		Show lots
//	@Description	Get lots created during last 7 days.
//	@Tags			lots
//	@Produce		json
//	@Success		200	{array}		lot_service.Lot
//	@Failure		400	{object}	apperror.AppError
//	@Failure		404	{object}	apperror.AppError
//	@Failure		418	{object}	apperror.AppError
//	@Router			/lots/week [get]
func (h *Handler) GetLastWeek(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Info("calculating date range and building query..")
	dateAfter := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	dateBefore := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	rQuery := fmt.Sprintf("created_at=%s:%s", dateAfter, dateBefore)

	lots, err := h.LotService.GetWithFilter(r.Context(), rQuery)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lots)

	return nil
}
