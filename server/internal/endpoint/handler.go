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

func (h *Handler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	var req CreateEndpoint

	if !helper.DecodeJSON(w, r, &req) {
		return
	}

	appID := r.URL.Query().Get("application_id")
	if appID == "" {
		http.Error(w, "application_id is required", http.StatusBadRequest)
		return
	}
	req.ApplicationID = appID

	result, err := h.svc.CreateEndpoint(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidURL) || errors.Is(err, ErrInvalidApplicationID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, result)
}
