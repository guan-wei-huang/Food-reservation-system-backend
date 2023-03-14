package user

type Config struct {
	DatabaseDsn string `envconfig:"DATABASE_DSN"`
	Port        int    `envconfig:"PORT"`
	LogFile     string `envconfig:"LOG_FILE"`
}
