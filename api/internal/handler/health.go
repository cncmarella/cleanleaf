package handler

import "net/http"

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// handleHealth backs the Fly.io health check. It must stay dependency-free so
// a failing mail provider never takes the deploy down.
func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, s.log, http.StatusOK, healthResponse{Status: "ok", Version: s.version})
}
