package main

import (
	"log"
	"reserve_restaurant/order"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseDsn string `envconfig:"DATABASE_DSN"`
	Port        int    `envconfig:"PORT"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal("parse env file error: ", err)
		return
	}

	// config.DatabaseDsn = "postgres://apple:123456@localhost:8892/apple?sslmode=disable"
	// config.Port = 3000

	r, err := order.NewOrderRepository(config.DatabaseDsn)
	if err != nil {
		log.Fatal("new order repo error: ", err)
		return
	}
	defer r.Close()

	s := order.NewService(r)
	log.Printf("Listening on port %v...", config.Port)
	log.Fatal(order.ListenGRPC(s, config.Port))
}
