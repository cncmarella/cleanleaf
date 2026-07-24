package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (s *Server) health(c *gin.Context) {
	writeJSON(c, http.StatusOK, healthResponse{Status: "ok", Version: s.version})
}
