package api

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLeval    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
}

func git NewConfig() *Config {
	return &Config{
		BindAddr: ":8082",
		LogLeval: "debug",
	}
}
