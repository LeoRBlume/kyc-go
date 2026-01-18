package router

import (
	"kyc-sim/internal/http/handler"

	"github.com/gin-gonic/gin"
)

type Deps struct {
	Health   *handler.HealthHandler
	Customer *handler.CustomerHandler
}

func NewRouter(deps Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", deps.Health.Get)

	v1 := r.Group("/v1/kyc")
	{
		v1.POST("/customers", deps.Customer.Create)
		v1.GET("/customers/:id", deps.Customer.GetByID)
		v1.PATCH("/customers/:id", deps.Customer.Patch)

	}

	return r
}
