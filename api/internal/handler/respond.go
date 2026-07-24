package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorBody struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func writeJSON(c *gin.Context, status int, v any) {
	c.JSON(status, v)
}

func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, errorBody{Error: msg})
}

func writeFieldErrors(c *gin.Context, fields map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, errorBody{
		Error:  "Please correct the highlighted fields.",
		Fields: fields,
	})
}
