package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Connection struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Connection{
		Conn: conn,
		Ch:   ch,
	}, nil
}

func (c *Connection) Close() {
	if c.Ch != nil {
		c.Ch.Close()
	}

	if c.Conn != nil {
		c.Conn.Close()
	}
}
