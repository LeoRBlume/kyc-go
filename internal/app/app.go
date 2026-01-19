package app

import (
	"kyc-sim/internal/config"
	"kyc-sim/internal/db"
	handlergraphql "kyc-sim/internal/graphql/handler"
	gqlschema "kyc-sim/internal/graphql/schema"
	"kyc-sim/internal/http/handler"
	"kyc-sim/internal/http/router"
	gormrepo "kyc-sim/internal/repository/gorm"
	svcimpl "kyc-sim/internal/service/impl"
	"kyc-sim/internal/worker"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Config config.Config
	DB     *gorm.DB
	Router *gin.Engine
}

func BuildApp() (*App, error) {
	cfg := config.Load()

	gdb, err := db.NewGormDB(cfg)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(gdb); err != nil {
		return nil, err
	}

	// Repos
	customerRepo := gormrepo.NewCustomerRepo(gdb)
	documentRepo := gormrepo.NewDocumentRepo(gdb)
	checkRepo := gormrepo.NewCheckRepo(gdb)
	jobRepo := gormrepo.NewJobRepo(gdb)

	// Services
	customerService := svcimpl.NewCustomerService(customerRepo, documentRepo)
	documentService := svcimpl.NewDocumentService(customerRepo, documentRepo)
	checkService := svcimpl.NewCheckService(customerRepo, checkRepo)
	kycService := svcimpl.NewKycService(customerRepo, jobRepo, checkRepo)
	jobService := svcimpl.NewJobService(jobRepo)
	readSvc := svcimpl.NewCustomerReadService(customerRepo, documentRepo, checkRepo, jobRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler()
	customerHandler := handler.NewCustomerHandler(customerService)
	documentHandler := handler.NewDocumentHandler(documentService)
	checkHandler := handler.NewCheckHandler(checkService, kycService)
	jobHandler := handler.NewJobHandler(jobService)

	gqlSchema, err := gqlschema.NewSchema(readSvc)
	if err != nil {
		return nil, err
	}

	gqlHandler := handlergraphql.NewGraphQLHandler(gqlSchema)

	// Router
	r := router.NewRouter(router.Deps{
		Health:   healthHandler,
		Customer: customerHandler,
		Document: documentHandler,
		Check:    checkHandler,
		Job:      jobHandler,
		GraphQL:  gqlHandler,
	})

	processor := worker.NewProcessor(jobRepo, checkRepo, customerRepo)
	runner := worker.NewRunner(jobRepo, processor, "worker-1")
	go runner.Start()

	return &App{
		Config: cfg,
		DB:     gdb,
		Router: r,
	}, nil
}
