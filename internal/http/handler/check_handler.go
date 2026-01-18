package handler

import (
	"kyc-sim/internal/domain"
	"net/http"

	"kyc-sim/internal/dto/responses"
	"kyc-sim/internal/http/middleware"
	svcif "kyc-sim/internal/service/interfaces"

	"github.com/gin-gonic/gin"
)

type CheckHandler struct {
	service svcif.CheckService
}

func NewCheckHandler(service svcif.CheckService) *CheckHandler {
	return &CheckHandler{service: service}
}

func (h *CheckHandler) List(c *gin.Context) {
	customerID := c.Param("id")

	checks, err := h.service.ListByCustomer(customerID)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := responses.ListChecksResponse{
		CustomerID: customerID,
		Results:    []responses.CheckResponse{},
		Summary:    responses.ChecksSummary{},
	}

	for _, chk := range checks {
		resp.Results = append(resp.Results, responses.CheckResponse{
			Type:      string(chk.Type),
			Status:    string(chk.Status),
			Score:     chk.Score,
			UpdatedAt: chk.UpdatedAt,
		})

		switch chk.Status {
		case domain.CheckPending:
			resp.Summary.Pending++
		case domain.CheckPass:
			resp.Summary.Pass++
		case domain.CheckFail:
			resp.Summary.Fail++
		case domain.CheckInconclusive:
			resp.Summary.Inconclusive++
		}
	}

	c.JSON(http.StatusOK, resp)
}
