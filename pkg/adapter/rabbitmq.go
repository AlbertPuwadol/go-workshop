package adapter

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IRabbitMQ interface {
	Publish(queuename string, ctx context.Context, message []byte) error
	Consume(queuename string, callback func([]byte))
	CreateQueue(queuename string) error
	Close()
}

type rabbitmq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQChannel(ctx context.Context, uri string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, err
}

func NewRabbitMQAdapter(conn *amqp.Connection, channel *amqp.Channel) *rabbitmq {
	return &rabbitmq{conn, channel}
}

func (r *rabbitmq) Close() {
	r.conn.Close()
	r.channel.Close()
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

func (r *rabbitmq) Consume(queuename string, callback func([]byte)) {
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
		log.Panic(err)
	}

	var forever chan interface{}

	go func() {
		for d := range msgs {
			callback(d.Body)
		}
	}()

	<-forever

}
