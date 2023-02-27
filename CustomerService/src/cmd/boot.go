package customercmd

import (
	"CustomerOrderMonoRepo/CustomerService/src/internal/handlers"
	"CustomerOrderMonoRepo/CustomerService/src/internal/repositories"
	"CustomerOrderMonoRepo/CustomerService/src/internal/services"
	"CustomerOrderMonoRepo/config"
	"CustomerOrderMonoRepo/shared/genericEndpoint/genericHandler"
	"CustomerOrderMonoRepo/shared/genericEndpoint/genericRepository"
	"CustomerOrderMonoRepo/shared/genericEndpoint/genericService"
	"github.com/erenkaratas99/COApiCore/pkg"
	"github.com/erenkaratas99/COApiCore/pkg/middleware"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// @title Customer Service
// @version 1.0
// @description This is a customer server to handle the requests about CRUD ops
// @termsOfService http://swagger.io/terms/

// @contact.name Eren Karataş
// @contact.url https://tesodev.com
// @contact.email bilgi@tesodev.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /customers
// @schemes http

func Boot(HTTPclient *pkg.RestClient, mc *mongo.Client, cfg *config.CustomerConfigs) {
	customerCol, err := pkg.GetMongoDbCollection(mc, cfg.DBName, cfg.ColName)
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.Pre(middleware.AddCorrelationID)
	e.Use(middleware.Recovery)
	e.Use(middleware.LoggingMiddleware)
	customerRepo := repositories.NewRepository(customerCol)
	customerGRepo := genericRepository.NewGenericRepository(customerCol)
	customerGService := genericService.NewGenericService(customerGRepo)
	customerGHandler := genericHandler.NewGenericHandler(customerGRepo, customerGService)
	customerService := services.NewService(customerRepo, HTTPclient)
	customerHandler := handlers.NewHandler(customerRepo, customerService, e, customerGHandler)
	e.Validator = pkg.NewValidation()
	customerHandler.InitEndpoints()
	log.Fatal(e.Start(cfg.BaseUrl + cfg.Port))
}
