package delivery

import (
	"errors"
	"net/http"

	"github.com/findardi/forflow-webhooks/server/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	svc *Service
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{svc: NewService(pool)}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Get(r.Context(), r.PathValue("appID"), r.PathValue("deliveryID"))
	if err != nil {
		writeError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) ListAttempts(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.ListAttempts(r.Context(), r.PathValue("appID"), r.PathValue("deliveryID"))
	if err != nil {
		writeError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, result)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidApplicationID), errors.Is(err, ErrInvalidDeliveryID):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrDeliveryNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
