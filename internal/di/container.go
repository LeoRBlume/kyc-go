package di

import (
	"kyc-sim/internal/config"
	"kyc-sim/internal/db"
	"kyc-sim/internal/http/handler"
	"kyc-sim/internal/http/router"
	gormrepo "kyc-sim/internal/repository/gorm"
	svcimpl "kyc-sim/internal/service/impl"

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

	// Services
	customerService := svcimpl.NewCustomerService(customerRepo)
	documentService := svcimpl.NewDocumentService(customerRepo, documentRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler()
	customerHandler := handler.NewCustomerHandler(customerService)
	documentHandler := handler.NewDocumentHandler(documentService)

	// Router
	r := router.NewRouter(router.Deps{
		Health:   healthHandler,
		Customer: customerHandler,
		Document: documentHandler,
	})

	return &App{
		Config: cfg,
		DB:     gdb,
		Router: r,
	}, nil
}
