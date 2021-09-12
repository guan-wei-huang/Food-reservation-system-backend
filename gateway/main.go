package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title reserve restaurant
// @version v1.0
// @description foodpanda
// @contact.name  gmail:a885131 at gmail.com
// @contact.url mailto:a885131@gmail.com
// @contact.email a885131@gmail.com
// @host localhost:7999
func main() {
	r := gin.Default()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("load .env failed: %v", err)
	}
	service, err := NewGatewayServer(os.Getenv("ORDER_URL"), os.Getenv("USER_URL"), os.Getenv("RESTAURANT_URL"))
	if err != nil {
		log.Fatalf("init service err: %v", err)
	}

	handler := NewHandler(service)
	handler.Register(r)
	r.Run()
}
