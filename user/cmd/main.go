package main

import (
	"log"
	"reserve_restaurant/user"

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

	repo, err := user.NewUserRepository(config.DatabaseDsn)
	if err != nil {
		log.Fatal("new restaurant repo error: ", err)
		return
	}
	defer repo.Close()

	log.Printf("Listening on port %v...", config.Port)
	s := user.NewService(repo)
	log.Fatal(user.ListenGRPC(s, config.Port))
}
