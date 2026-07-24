package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cncmarella/cleanleaf/api/internal/service"
)

// maxContactBody caps the request body well above a legitimate enquiry.
const maxContactBody = 16 << 10 // 16 KiB

type contactPayload struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Source  string `json:"source"`
	Website string `json:"website"`
}

type contactResponse struct {
	Message string `json:"message"`
}

// contactHandler adapts HTTP requests to the contact service.
type contactHandler struct {
	svc *service.Contact
}

func (h *contactHandler) submit(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxContactBody)

	var p contactPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(c, http.StatusRequestEntityTooLarge, "That message is too long.")
			return
		}
		writeError(c, http.StatusBadRequest, "We could not read that request.")
		return
	}

	// contactPayload and service.ContactInput have identical fields, so the
	// conversion is a compiler-checked copy with no field-by-field mapping.
	err := h.svc.Submit(c.Request.Context(), service.ContactInput(p))

	var validation *service.ValidationError
	var delivery *service.DeliveryError
	switch {
	case err == nil:
		writeJSON(c, http.StatusAccepted, contactResponse{Message: "Thanks for reaching out. We'll be in touch shortly."})
	case errors.As(err, &validation):
		writeFieldErrors(c, validation.Fields)
	case errors.As(err, &delivery):
		writeError(c, http.StatusBadGateway, "We could not send your message right now. Please call us on 8341099962.")
	default:
		writeError(c, http.StatusInternalServerError, "Something went wrong on our side.")
	}
}
