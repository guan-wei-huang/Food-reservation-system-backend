package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func (handler *Handler) Register(r *gin.Engine) {
	r.POST("/login", handler.UserLogin)
	r.POST("/register", handler.NewUser)

	r.GET("/restaurant/:rid", handler.GetRestaurantMenu)
	r.POST("/restaurant/:rid/food", handler.CreateFood)
	r.POST("/restaurant", handler.CreateRestaurant)

	r.GET("/restaurant/:location", handler.GetNearbyRestaurant)

	r.GET("/order/:oid", handler.GetOrder)
	r.GET("/user/order", handler.GetOrderForUser)
	r.POST("/order", handler.CreateOrder)

	if mode := gin.Mode(); mode == gin.DebugMode {
		url := ginSwagger.URL("http://localhost:7999/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
}
