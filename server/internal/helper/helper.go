package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	maxJSONBody = 4 << 10
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBody)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil && !errors.Is(err, io.EOF) {
		http.Error(w, "invalid JSON Body", http.StatusBadRequest)
		return false
	}

	return true
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func ToUUID(u string) (pgtype.UUID, error) {
	var id pgtype.UUID

	if u == "" {
		return id, nil
	}

	if err := id.Scan(u); err != nil {
		return id, fmt.Errorf("parse to uuid: %w", err)
	}

	return id, nil
}
