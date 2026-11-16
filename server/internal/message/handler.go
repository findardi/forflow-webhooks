package message

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
	return &Handler{
		svc: NewService(pool),
	}
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req CreateMessage

	if !helper.DecodeJSON(w, r, &req) {
		return
	}

	req.ApplicationID = r.PathValue("appID")

	result, err := h.svc.CreateMessage(r.Context(), req)
	if err != nil {
		writeError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, result)
}

func (h *Handler) ListMessage(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.ListByEndpoint(r.Context(), r.PathValue("appID"), r.PathValue("endpointID"))
	if err != nil {
		writeError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, result)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidApplicationID), errors.Is(err, ErrInvalidEndpointID), errors.Is(err, ErrInvalidPayload):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrApplicationNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
