package handler

import (
	"net/http"

	"kyc-sim/internal/dto/requests"
	"kyc-sim/internal/dto/responses"
	"kyc-sim/internal/http/middleware"
	svcif "kyc-sim/internal/service/interfaces"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	service svcif.DocumentService
}

func NewDocumentHandler(service svcif.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) Create(c *gin.Context) {
	customerID := c.Param("id")

	var req requests.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.WriteBindError(c, err)
		return
	}

	doc, err := h.service.Create(customerID, req)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, responses.DocumentResponse{
		ID:         doc.ID,
		CustomerID: doc.CustomerID,
		Kind:       string(doc.Kind),
		Status:     string(doc.Status),
		UploadedAt: doc.UploadedAt,
	})
}

func (h *DocumentHandler) List(c *gin.Context) {
	customerID := c.Param("id")

	docs, err := h.service.List(customerID)
	if err != nil {
		middleware.WriteError(c, err)
		return
	}

	resp := make([]responses.DocumentResponse, 0, len(docs))
	for _, d := range docs {
		resp = append(resp, responses.DocumentResponse{
			ID:         d.ID,
			CustomerID: d.CustomerID,
			Kind:       string(d.Kind),
			Status:     string(d.Status),
			UploadedAt: d.UploadedAt,
		})
	}

	c.JSON(http.StatusOK, responses.ListDocumentsResponse{
		CustomerID: customerID,
		Documents:  resp,
	})
}
