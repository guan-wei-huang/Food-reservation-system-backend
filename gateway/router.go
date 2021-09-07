package main

import "github.com/gin-gonic/gin"

func (handler *Handler) Register(r *gin.Engine) {
	r.POST("/login", handler.UserLogin)
	r.POST("/register", handler.NewUser)

	r.GET("/restaurant/:rid", handler.GetRestaurantMenu)
	r.POST("/restaurant/:rid/food", handler.CreateFood)
	r.POST("/restaurant", handler.CreateRestaurant)

	r.GET("/restaurants/:location", handler.GetNearbyRestaurant)

	r.GET("/order/:oid", handler.GetOrder)
	r.GET("/user/order", handler.GetOrderForUser)
	r.POST("/order", handler.CreateOrder)

}
