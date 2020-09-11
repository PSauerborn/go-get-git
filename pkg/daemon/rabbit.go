package daemon

import (
	"fmt"
	"encoding/json"
	"github.com/streadway/amqp"
	log "github.com/sirupsen/logrus"
	events "github.com/PSauerborn/go-get-git/pkg/events"
)

// define function used to send message over rabbitmq server
func sendRabbitPayload(event events.Event) error {
	conn, err := amqp.Dial(RabbitQueueUrl)
	if err != nil {
		log.Error(fmt.Errorf("unable to connect to rabbitmq server: %s", err))
		return err
	}
	defer conn.Close()

	// create channel on rabbitmq server
	channel, err := conn.Channel()
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbitmq channel: %s", err))
		return err
	}
	// declare events exchange with fanout type
	err = channel.ExchangeDeclare("events", "fanout", false, true, false, false, nil)
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbitmq exchange: %s", err))
		return err
	}
	// construct payload and send over rabbit server
	body, _ := json.Marshal(&event)
	payload := amqp.Publishing{ ContentType: "application/json", Body: []byte(body) }
	err = channel.Publish("events", "", false, false, payload)
	if err != nil {
		log.Error(fmt.Errorf("unable to send payload over rabbitmq server: %s", err))
		return err
	}
	log.Info(fmt.Sprintf("successfully sent payload %+v over rabbitMQ exchange", event))
	return nil
}

func RabbitListener(handler func(payload []byte)) error {
	log.Info(fmt.Sprintf("connecting to rabbitmq server at %s", RabbitQueueUrl))
	// connect to rabbitmq server using queue url
	conn, err := amqp.Dial(RabbitQueueUrl)
	if err != nil {
		log.Error(fmt.Errorf("unable to connect to rabbitmq server: %s", err))
		return err
	}
	defer conn.Close()

	// create channel on rabbitmq server
	channel, err := conn.Channel()
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbitmq channel: %s", err))
		return err
	}
	// declare queue
	queue, err := channel.QueueDeclare(QueueName, false, false, true, false, nil)
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbit queue: %v", err))
		return err
	}
	// bind queue to exchange
	err = channel.QueueBind(queue.Name, "", EventExchangeName, false, nil)
	if err != nil {
		log.Error(fmt.Errorf("unable to bind to exchange: %v", err))
		return err
	}
	// consume to messages on channel
	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbit queue: %v", err))
		return err
	}
	// start goroutine to handle messages
	forever :=make(chan bool)
	go func() { for d := range messages { handler(d.Body) }}()
	<-forever

	return nil
}