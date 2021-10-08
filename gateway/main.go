package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          string `envconfig:"PORT"`
	OrderUrl      string `envconfig:"ORDER_URL"`
	UserUrl       string `envconfig:"USER_URL"`
	RestaurantUrl string `envconfig:"RESTAURANT_URL"`
}

// @title reserve restaurant
// @version v1.0
// @description foodpanda
// @contact.name  gmail:a885131 at gmail.com
// @contact.url mailto:a885131@gmail.com
// @contact.email a885131@gmail.com
// @host localhost:8080
func main() {
	r := gin.Default()

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("parse env file failed: %v", err)
		return
	}

	// config.OrderUrl = "localhost:3000"
	// config.RestaurantUrl = "localhost:3001"
	// config.UserUrl = "localhost:3002"
	// config.Port = "8080"

	service, err := NewGatewayServer(config.OrderUrl, config.UserUrl, config.RestaurantUrl)
	if err != nil {
		log.Fatalf("init service err: %v", err)
	}

	handler := NewHandler(service)
	handler.Register(r, config.Port)
	r.Run(fmt.Sprintf(":%v", config.Port))
}
