package adapter

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IRabbitMQ interface {
	Publish(queuename string, ctx context.Context, message []byte) error
	Consume(queuename string) (interface{}, error)
	CreateQueue(queuename string) error
}

type rabbitmq struct {
	channel *amqp.Channel
}

func NewRabbitMQChannel(ctx context.Context, uri string) (*amqp.Channel, error) {
	rabbitconn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	rabbitch, err := rabbitconn.Channel()
	if err != nil {
		return nil, err
	}
	return rabbitch, err
}

func NewRabbitMQAdapter(channel *amqp.Channel) *rabbitmq {
	return &rabbitmq{channel}
}

func (r *rabbitmq) CreateQueue(queuename string) error {
	_, err := r.channel.QueueDeclare(
		queuename, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}

func (r *rabbitmq) Publish(queuename string, ctx context.Context, message []byte) error {
	err := r.channel.PublishWithContext(ctx,
		"",        // exchange
		queuename, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})

	return err
}

func (r *rabbitmq) Consume(queuename string) (interface{}, error) {
	msgs, err := r.channel.Consume(
		queuename, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	for d := range msgs {
		result = append(result, d)
	}
	return result, nil
}
