package api

import (
	"CoffeeShopOnGo/api/handlers"
	"CoffeeShopOnGo/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(svc *service.Service) *gin.Engine {
	router := gin.Default()

	productHandler := handlers.NewProductHandler(svc)

	api := router.Group("/api")
	{
		products := api.Group("/products")
		{
			products.GET("/", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
			products.POST("/", productHandler.Create)
            products.PATCH("/:id/price", productHandler.UpdatePrice)
		}
	}
	return router
}