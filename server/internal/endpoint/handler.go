package endpoint

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if !helper.DecodeJSON(w, r, &req) {
		return
	}
	req.ApplicationID = r.PathValue("appID")

	result, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writeError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.ListByApplication(r.Context(), r.PathValue("appID"))
	if err != nil {
		writeError(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), r.PathValue("appID"), r.PathValue("endpointID")); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidURL), errors.Is(err, ErrInvalidApplicationID), errors.Is(err, ErrInvalidEndpointID):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrApplicationNotFound), errors.Is(err, ErrEndpointNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrEndpointExists):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
