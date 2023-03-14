package main

import (
	"fmt"

	_ "reserve_restaurant/gateway/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (handler *Handler) Register(r *gin.Engine, port string) {
	r.POST("/login", handler.UserLogin)
	r.POST("/register", handler.NewUser)

	r.GET("/restaurant/:rid", handler.GetRestaurantMenu)
	r.POST("/restaurant/:rid/food", handler.CreateFood)
	r.POST("/restaurant", handler.CreateRestaurant)

	r.GET("/search/:location", handler.GetNearbyRestaurant)

	orderGroup := r.Group("/order")
	orderGroup.Use(Auth)
	{
		orderGroup.GET("/:oid", handler.GetOrder)
		orderGroup.GET("/", handler.GetOrderForUser)
		orderGroup.POST("/", handler.CreateOrder)
	}

	if mode := gin.Mode(); mode == gin.DebugMode {
		url := ginSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", port))
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))

		handler := promhttp.Handler()
		r.GET("/metrics", func(c *gin.Context) {
			handler.ServeHTTP(c.Writer, c.Request)
		})
	}
}
