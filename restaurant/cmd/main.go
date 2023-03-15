package main

import (
	"log"
	"reserve_restaurant/restaurant"
)

func main() {
	config := restaurant.GetConfig()

	// config.DatabaseDsn = "postgres://apple:123456@localhost:8889/apple?sslmode=disable"
	// config.Port = 3001

	repo, err := restaurant.NewRestaurantRepository(config.DatabaseDsn)
	if err != nil {
		log.Fatal("new restaurant repo error: ", err)
		return
	}
	defer repo.Close()

	s := restaurant.NewService(repo)
	log.Printf("Listening on port %v...", config.Port)
	log.Fatal(restaurant.ListenGRPC(s, config))
}
