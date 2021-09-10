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

	r, err := order.NewOrderRepository(config.DatabaseDsn)
	if err != nil {
		log.Fatal("new restaurant repo error: ", err)
		return
	}
	defer r.Close()

	log.Printf("Listening on port %v...", config.Port)
	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, config.Port))
}
