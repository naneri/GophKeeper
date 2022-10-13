package config

// Config keeps all the configuration parameters of the Server app
type Config struct {
	AppEnv        string `env:"APP_ENV,required"`
	DatabaseDsn   string `env:"DATABASE_DSN,required"`
	ServerAddress string `env:"SERVER_ADDRESS,required"`
	AppKey        string `env:"APP_KEY,required"`
}
