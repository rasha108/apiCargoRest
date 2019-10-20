package rabbitclient

// Config stores rabbitmq connection parameters for mail-server
type Config struct {
	Host      string `toml:"host"`
	Port      int    `toml:"port"`
	UserName  string `toml:"username"`
	Password  string `toml:"password"`
	SendQueue string `toml:"send_queue"`
}

// setDefault is method for adding fedault field values
func (c *Config) SetDefaults() {
	const defaultRabbitMQLoginAndPassword = "guest"
	c.Port = 5672
	c.UserName = defaultRabbitMQLoginAndPassword
	c.Password = defaultRabbitMQLoginAndPassword
	c.SendQueue = "mail"
}
