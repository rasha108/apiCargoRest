package rabbitclient

// Config stores rabbitmq connection parameters for mail-server
type Config struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	UserName  string `yaml:"username"`
	Password  string `yaml:"password"`
	SendQueue string `yaml:"send_queue"`
}

// setDefault is method for adding fedault field values
func (c *Config) SetDefaults() {
	const defaultRabbitMQLoginAndPassword = "guest"
	c.Port = 5672
	c.UserName = defaultRabbitMQLoginAndPassword
	c.Password = defaultRabbitMQLoginAndPassword
	c.SendQueue = "mail"
}
