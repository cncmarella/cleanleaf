package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// errorBody is the single error shape every endpoint returns, so the frontend
// only ever has to parse one thing.
type errorBody struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func writeJSON(w http.ResponseWriter, log *slog.Logger, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Headers are already flushed, so this can only be logged.
		log.Error("write json response", "error", err)
	}
}

func writeError(w http.ResponseWriter, log *slog.Logger, status int, msg string) {
	writeJSON(w, log, status, errorBody{Error: msg})
}

func writeFieldErrors(w http.ResponseWriter, log *slog.Logger, fields map[string]string) {
	writeJSON(w, log, http.StatusUnprocessableEntity, errorBody{
		Error:  "Please correct the highlighted fields.",
		Fields: fields,
	})
}
