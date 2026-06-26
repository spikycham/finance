package business

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/spikycham/finance/network"
)

type Handler struct {
	s *Service
}

type CreateItemRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Type   string    `json:"type" validate:"required,oneof=food supply transport treatment study job wear workout entertainment house income others"`
	Amount float64   `json:"amount" validate:"required,gt=0"`
	Note   string    `json:"note" validate:"max=100"`
}

type GetItemsResponse struct {
	Items []Item `json:"items"`
}

func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

var validate = validator.New()

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	body, err := network.ReadBody[CreateItemRequest](r)
	if err != nil {
		network.ResponseError(w, http.StatusBadRequest, network.ErrInvalidJSON.Error())
		return
	}

	if err := validate.Struct(body); err != nil {
		network.ResponseError(w, http.StatusBadRequest, network.ErrMissingFields.Error())
		return
	}

	if err := h.s.CreateRecordItem(r.Context(), &Item{
		ID:        0,
		UserID:    body.UserID,
		Type:      body.Type,
		Amount:    body.Amount,
		Note:      body.Note,
		CreatedAt: time.Now(),
	}); err != nil {
		network.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	network.ResponseMessage(w, http.StatusCreated, "success to create a new financial record")
}

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	userID, err := uuid.Parse(q.Get("user_id"))
	if err != nil {
		network.ResponseError(w, http.StatusBadRequest, network.ErrMissingFields.Error())
		return
	}

	startTime, err := strconv.Atoi(q.Get("start_time"))
	if err != nil {
		network.ResponseError(w, http.StatusBadRequest, network.ErrMissingFields.Error())
		return
	}

	items, err := h.s.GetYearItems(r.Context(), userID, time.UnixMilli(int64(startTime)))
	if err != nil {
		network.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	network.ResponseJSON(w, http.StatusOK, &GetItemsResponse{
		Items: items,
	})
}
