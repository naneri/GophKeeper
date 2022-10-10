package config

// Config keeps all the configuration parameters of the Server app
type Config struct {
	AppEnv           string `env:"APP_ENV" envDefault:"dev"`
	ServerAddress    string `env:"SERVER_ADDRESS,required" envDefault:":8083"`
	RemoteServerPort string `env:"REMOTE_SERVER_PORT,required" envDefault:"http://localhost:8080"`
	PasswordCookie   string
}
