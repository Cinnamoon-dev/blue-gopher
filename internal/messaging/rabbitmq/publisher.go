package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	Ch *amqp.Channel
}

func NewRabbitPublisher(conn Connection) *RabbitPublisher {
	return &RabbitPublisher{
		Ch: conn.Ch,
	}
}

func (p *RabbitPublisher) Publish(
	exchange string,
	routingKey string,
	message []byte,
) error {
	return p.Ch.PublishWithContext(
		context.Background(),
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

func (p *RabbitPublisher) PublishEvent(
	exchange string,
	routingKey string,
	event any,
) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.Publish(exchange, routingKey, body)
}
