package api

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLeval    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	BasePath    string
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8083",
		LogLeval: "debug",
		BasePath: "/api/v1",
	}
}
