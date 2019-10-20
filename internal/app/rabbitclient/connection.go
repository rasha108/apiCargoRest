package rabbitclient

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

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
