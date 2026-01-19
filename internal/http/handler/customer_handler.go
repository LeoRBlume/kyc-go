package handler

import (
	"kyc-sim/internal/dto/requests"
	"kyc-sim/internal/dto/responses"
	"kyc-sim/internal/http/middleware"
	svcif "kyc-sim/internal/service/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service svcif.CustomerService
}

func NewCustomerHandler(service svcif.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req requests.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteBindError(c, err)
		return
	}

	customer, err := h.service.Create(req)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, responses.CustomerCreatedResponse{
		ID:        customer.ID,
		Type:      string(customer.Type),
		Status:    string(customer.Status),
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	})
}

func (h *CustomerHandler) Patch(c *gin.Context) {
	id := c.Param("id")

	var req requests.PatchCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteBindError(c, err)
		return
	}

	customer, err := h.service.Patch(id, req)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.CustomerSummaryResponse{
		ID:        customer.ID,
		Status:    string(customer.Status),
		UpdatedAt: customer.UpdatedAt,
	})
}

func (h *CustomerHandler) Submit(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Submit(id); err != nil {
		middleware.WriteError(c, err)
		return
	}

	customer, _ := h.service.GetByID(id)

	c.JSON(http.StatusOK, responses.SubmitResponse{
		ID:                customer.ID,
		Status:            string(customer.Status),
		SubmittedAt:       customer.UpdatedAt,
		RequiredDocuments: nil, // apenas informativo; regra já validou
		MissingDocuments:  []string{},
	})
}
