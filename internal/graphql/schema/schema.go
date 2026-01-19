package gqlschema

import (
	"time"

	"kyc-sim/internal/service/interfaces"

	"github.com/graphql-go/graphql"
)

func mustISOTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func NewSchema(readSvc interfaces.CustomerReadService) (graphql.Schema, error) {
	jobProgressType := graphql.NewObject(graphql.ObjectConfig{
		Name: "JobProgress",
		Fields: graphql.Fields{
			"total":  &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"done":   &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"failed": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		},
	})

	jobType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Job",
		Fields: graphql.Fields{
			"jobId":      &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"kind":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"status":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"customerId": &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"createdAt":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"updatedAt":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"progress":   &graphql.Field{Type: graphql.NewNonNull(jobProgressType)},
		},
	})

	documentType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Document",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"customerId": &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"kind":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"status":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"uploadedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"fileUrl":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	})

	checkType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Check",
		Fields: graphql.Fields{
			"checkType": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"status":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"score":     &graphql.Field{Type: graphql.Int},
			"details":   &graphql.Field{Type: graphql.String},
			"updatedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	})

	checkSummaryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "ChecksSummary",
		Fields: graphql.Fields{
			"pending":      &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"pass":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"fail":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"inconclusive": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		},
	})

	customerType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Customer",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"type":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"status":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"createdAt": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"updatedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},

			"documents":        &graphql.Field{Type: graphql.NewList(documentType)},
			"checks":           &graphql.Field{Type: graphql.NewList(checkType)},
			"checksSummary":    &graphql.Field{Type: graphql.NewNonNull(checkSummaryType)},
			"lastRunChecksJob": &graphql.Field{Type: jobType},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"customer": &graphql.Field{
				Type: customerType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					id, _ := p.Args["id"].(string)
					view, err := readSvc.GetCustomerView(id)
					if err != nil {
						return nil, err
					}

					// Adaptar view para map simples (graphql-go trabalha bem assim)
					docMaps := make([]map[string]any, 0, len(view.Documents))
					for _, d := range view.Documents {
						docMaps = append(docMaps, map[string]any{
							"id": d.ID, "customerId": d.CustomerID, "kind": d.Kind, "status": d.Status,
							"uploadedAt": mustISOTime(d.UploadedAt), "fileUrl": d.FileURL,
						})
					}

					checkMaps := make([]map[string]any, 0, len(view.Checks))
					for _, c := range view.Checks {
						checkMaps = append(checkMaps, map[string]any{
							"checkType": c.CheckType, "status": c.Status, "score": c.Score,
							"details": c.Details, "updatedAt": mustISOTime(c.UpdatedAt),
						})
					}

					var jobMap map[string]any
					if view.LastJob != nil {
						jobMap = map[string]any{
							"jobId":      view.LastJob.JobID,
							"kind":       view.LastJob.Kind,
							"status":     view.LastJob.Status,
							"customerId": view.LastJob.CustomerID,
							"createdAt":  mustISOTime(view.LastJob.CreatedAt),
							"updatedAt":  mustISOTime(view.LastJob.UpdatedAt),
							"progress": map[string]any{
								"total":  view.LastJob.Progress.Total,
								"done":   view.LastJob.Progress.Done,
								"failed": view.LastJob.Progress.Failed,
							},
						}
					}

					return map[string]any{
						"id":        view.ID,
						"type":      view.Type,
						"status":    view.Status,
						"createdAt": mustISOTime(view.CreatedAt),
						"updatedAt": mustISOTime(view.UpdatedAt),
						"documents": docMaps,
						"checks":    checkMaps,
						"checksSummary": map[string]any{
							"pending":      view.Summary.Pending,
							"pass":         view.Summary.Pass,
							"fail":         view.Summary.Fail,
							"inconclusive": view.Summary.Inconclusive,
						},
						"lastRunChecksJob": jobMap,
					}, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery})
}
