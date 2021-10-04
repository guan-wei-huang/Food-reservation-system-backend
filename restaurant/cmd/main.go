package main

import (
	"log"
	"reserve_restaurant/restaurant"

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

	config.DatabaseDsn = "postgres://apple:123456@localhost:8889/apple?sslmode=disable"
	config.Port = 3001

	repo, err := restaurant.NewRestaurantRepository(config.DatabaseDsn)
	if err != nil {
		log.Fatal("new restaurant repo error: ", err)
		return
	}
	defer repo.Close()

	s := restaurant.NewService(repo)
	log.Printf("Listening on port %v...", config.Port)
	log.Fatal(restaurant.ListenGRPC(s, config.Port))
}
