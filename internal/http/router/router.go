package router

import (
	handlergraphql "kyc-sim/internal/graphql/handler"
	"kyc-sim/internal/http/handler"

	"github.com/gin-gonic/gin"
)

type Deps struct {
	Health   *handler.HealthHandler
	Customer *handler.CustomerHandler
	Document *handler.DocumentHandler
	Check    *handler.CheckHandler
	Job      *handler.JobHandler
	GraphQL  *handlergraphql.GraphQLHandler
}

func NewRouter(deps Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	//Health
	r.GET("/health", deps.Health.Get)

	v1 := r.Group("/v1/kyc")
	{
		//customer
		v1.POST("/customers", deps.Customer.Create)
		v1.POST("/graphql", deps.GraphQL.Serve)
		v1.PATCH("/customers/:id", deps.Customer.Patch)

		v1.POST("/customers/:id/submit", deps.Customer.Submit)

		//document
		v1.POST("/customers/:id/documents", deps.Document.Create)
		v1.GET("/customers/:id/documents", deps.Document.List)

		//check
		v1.GET("/customers/:id/checks", deps.Check.List)
		v1.POST("/customers/:id/checks/run", deps.Check.Run)

		//job
		v1.GET("/jobs/:jobId", deps.Job.Get)
	}

	return r
}
