package handler

import (
	"kyc-sim/internal/domain"
	"kyc-sim/internal/dto/requests"
	"net/http"

	"kyc-sim/internal/dto/responses"
	"kyc-sim/internal/http/middleware"
	svcif "kyc-sim/internal/service/interfaces"

	"github.com/gin-gonic/gin"
)

type CheckHandler struct {
	service    svcif.CheckService
	kycService svcif.KycService
}

func NewCheckHandler(service svcif.CheckService, kycService svcif.KycService) *CheckHandler {
	return &CheckHandler{service: service, kycService: kycService}
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

func (h *CheckHandler) Run(c *gin.Context) {
	customerID := c.Param("id")

	var req requests.RunChecksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteBindError(c, err)
		return
	}

	job, reused, err := h.kycService.EnqueueRunChecks(customerID, req)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	status := "QUEUED"
	if reused {
		status = "ALREADY_RUNNING"
	}

	c.JSON(http.StatusAccepted, responses.EnqueueChecksResponse{
		JobID:      job.ID,
		CustomerID: job.CustomerID,
		Status:     status,
		QueuedAt:   job.CreatedAt,
	})
}
