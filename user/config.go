package user

import (
	"log"
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseDsn string `envconfig:"DATABASE_DSN"`
	Port        int    `envconfig:"PORT"`
	LogFile     string `envconfig:"LOG_FILE"`
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
		if _, err = os.Stat(config.LogFile); os.IsNotExist(err) {
			if _, err := os.Create(config.LogFile); err != nil {
				log.Fatal("create log file error: ", err)
			}
		}
	})
	return &config
}
