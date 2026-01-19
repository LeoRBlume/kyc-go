package handler

import (
	"net/http"

	"kyc-sim/internal/http/middleware"
	svcif "kyc-sim/internal/service/interfaces"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	service svcif.JobService
}

func NewJobHandler(service svcif.JobService) *JobHandler {
	return &JobHandler{service: service}
}

func (h *JobHandler) Get(c *gin.Context) {
	jobID := c.Param("jobId")

	resp, err := h.service.Get(jobID)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
