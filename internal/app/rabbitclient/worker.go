package rabbitclient

import (
	"encoding/json"
	"fmt"
	"net/mail"

	"github.com/streadway/amqp"
)

type Mail struct {
	From        *mail.Address   `json:"from,omitempty"`
	To          []*mail.Address `json:"to,omitempty"`
	Cc          []*mail.Address `json:"cc,omitempty"`
	Bcc         []*mail.Address `json:"bcc,omitempty"`
	Subject     string          `json:"subject,omitempty"`
	Text        string          `json:"text,omitempty"`
	Attachments []Attachment    `json:"attachments,omitempty"`
}

type Attachment struct {
	Name    string `json:"name,omitempty"`    // File name
	Content string `json:"content,omitempty"` // Base64 encoded []byte
}

// Config stores rabbitmq connection parameters for mail-server
type Config struct {
	Host      string `yaml:"host,omitempty"`
	Port      int    `yaml:"port,omitempty"`
	UserName  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	SendQueue string `yaml:"send_queue,omitempty"`
}

// Connection is store for rabbitmq conncetions
type Connection struct {
	queue   string
	channel *amqp.Channel
}

func NewConnection(config Config) (*Connection, error) {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.UserName, config.Password, config.Host, config.Port)
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	mailServer := &Connection{
		queue:   config.SendQueue,
		channel: ch,
	}

	return mailServer, err
}

func (ms *Connection) Send(mail *Mail) error {
	body, err := json.Marshal(mail)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{ContentType: "application/json", Body: body}
	err = ms.channel.Publish("", ms.queue, false, false, msg)
	if err != nil {
		return err
	}
	return nil
}

func (ms *Connection) Close() error {
	return ms.channel.Close()
}
