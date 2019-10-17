package api

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLeval    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
	BasePath    string
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8084",
		LogLeval: "debug",
		BasePath: "/api/v1",
	}
}
