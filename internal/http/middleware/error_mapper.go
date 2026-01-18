package middleware

import (
	"kyc-sim/internal/domain"
	"kyc-sim/internal/dto/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Domain errors
	if de, ok := err.(*domain.DomainError); ok {
		status := http.StatusInternalServerError
		switch de.Code {
		case domain.ErrNotFound:
			status = http.StatusNotFound
		case domain.ErrValidation:
			status = http.StatusUnprocessableEntity
		default:
			status = http.StatusInternalServerError
		}

		c.JSON(status, responses.ErrorEnvelope{
			Error: responses.ErrorBody{
				Code:    string(de.Code),
				Message: de.Message,
				Details: de.Details,
			},
		})
		return
	}

	// Fallback
	c.JSON(http.StatusInternalServerError, responses.ErrorEnvelope{
		Error: responses.ErrorBody{
			Code:    "INTERNAL_ERROR",
			Message: "unexpected error",
		},
	})
}
