package endpoint

import (
	"errors"
	"net/http"

	"github.com/findardi/forflow-webhooks/server/db/store/sqlc"
	"github.com/findardi/forflow-webhooks/server/internal/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	svc *Service
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		svc: NewService(sqlc.New(pool)),
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

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidURL), errors.Is(err, ErrInvalidApplicationID):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrApplicationNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrEndpointExists):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
