package main

import (
	"log"
	"reserve_restaurant/user"
)

func main() {
	config := user.GetConfig()

	// config.DatabaseDsn = "postgres://apple:123456@localhost:8890/apple?sslmode=disable"
	// config.Port = 3002

	repo, err := user.NewUserRepository(config.DatabaseDsn)
	if err != nil {
		log.Fatal("new restaurant repo error: ", err)
		return
	}
	defer repo.Close()

	s := user.NewService(repo)
	log.Printf("Listening on port %v...", config.Port)
	log.Fatal(user.ListenGRPC(s, config))
}
