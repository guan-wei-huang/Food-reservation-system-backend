package restaurant

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseDsn string `envconfig:"DATABASE_DSN"`
	Port        int    `envconfig:"PORT"`
	ApiKey      string `envconfig:"GOOGLE_API_KEY"`
}

var config Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal("parse env file error: ", err)
			return
		}
	})
	return &config
}
