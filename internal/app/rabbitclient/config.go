package rabbitclient

// Config stores rabbitmq connection parameters for mail-server
type Config struct {
	Host      string `yaml:"host,omitempty"`
	Port      int    `yaml:"port,omitempty"`
	UserName  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	SendQueue string `yaml:"send_queue,omitempty"`
}

// setDefault is method for adding fedault field values
func (c *Config) SetDefaults() {
	const defaultRabbitMQLoginAndPassword = "guest"
	c.Port = 5672
	c.UserName = defaultRabbitMQLoginAndPassword
	c.Password = defaultRabbitMQLoginAndPassword
	c.SendQueue = "mail"
}
