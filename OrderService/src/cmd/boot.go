package ordercmd

import (
	"CustomerOrderMonoRepo/OrderService/src/internal/handlers"
	"CustomerOrderMonoRepo/OrderService/src/internal/repositories"
	"CustomerOrderMonoRepo/OrderService/src/internal/services"
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

// @title Order Service
// @version 1.0
// @description This is a order server to handle the requests about CRUD ops
// @termsOfService http://swagger.io/terms/

// @contact.name Eren Karataş
// @contact.url https://tesodev.com
// @contact.email bilgi@tesodev.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /orders
// @schemes http

func Boot(HTTPclient *pkg.RestClient, mc *mongo.Client, cfg *config.OrderConfigs) {
	orderCol, err := pkg.GetMongoDbCollection(mc, cfg.DBName, cfg.ColName)
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.Pre(middleware.AddCorrelationID)
	e.Use(middleware.Recovery)
	e.Use(middleware.LoggingMiddleware)

	orderRepo := repositories.NewRepository(orderCol)
	orderGRepo := genericRepository.NewGenericRepository(orderCol)
	orderGService := genericService.NewGenericService(orderGRepo)
	orderGHandler := genericHandler.NewGenericHandler(orderGRepo, orderGService)
	orderService := services.NewService(orderRepo, HTTPclient)
	orderHandler := handlers.NewHandler(orderRepo, orderService, e, orderGHandler)
	e.Validator = pkg.NewValidation()
	orderHandler.InitEndpoints()

	log.Fatal(e.Start(cfg.BaseUrl + cfg.Port))
}
