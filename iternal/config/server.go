package config

type Options struct {
	Host string `env:"SERV_HOST" envDefault:"localhost"`
	Port string `env:"SERV_PORT" envDefault:"8080"`
}
