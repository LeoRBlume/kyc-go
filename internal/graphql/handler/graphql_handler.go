package handlergraphql

import (
	"net/http"

	"kyc-sim/internal/domain"
	"kyc-sim/internal/dto/responses"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

type GraphQLHandler struct {
	schema graphql.Schema
}

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func NewGraphQLHandler(schema graphql.Schema) *GraphQLHandler {
	return &GraphQLHandler{schema: schema}
}

func (h *GraphQLHandler) Serve(c *gin.Context) {
	var req graphqlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorEnvelope{
			Error: responses.ErrorBody{
				Code: "BAD_REQUEST", Message: "invalid graphql request payload",
				Details: map[string]any{"cause": err.Error()},
			},
		})
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         h.schema,
		RequestString:  req.Query,
		VariableValues: req.Variables,
		Context:        c.Request.Context(),
	})

	if len(result.Errors) > 0 {
		// opcional: mapear DomainError para extensions
		// aqui a lib já serializa `errors`.
	}

	_ = domain.ErrInternal
	c.JSON(http.StatusOK, result)
}
