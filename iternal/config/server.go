package config

type Options struct {
	Address  string `env:"ADDRESS"`
	Filename string `env:"FILE_STORAGE_PATH"`
	DBDSN    string `env:"DATABASE_DSN"`
	Key      string `env:"KEY"`
}
