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

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	customer, err := h.service.GetByID(id)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.GetCustomerResponse{
		ID:        customer.ID,
		Type:      string(customer.Type),
		Status:    string(customer.Status),
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	})
}
