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

	repo, err := restaurant.NewRestaurantRepository(config.DatabaseDsn)
	if err != nil {
		log.Fatal("new restaurant repo error: ", err)
		return
	}
	defer repo.Close()

	log.Printf("Listening on port %v...", config.Port)
	s := restaurant.NewService(repo)
	log.Fatal(restaurant.ListenGRPC(s, config.Port))
}
