package main

import (
	"log"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func PublishToRabbit(payload string) {
	url := "amqp://guest:guest@localhost:5672/"
	connection, dialError := amqp.Dial(url)

	closeListener := make(chan *amqp.Error)
	connection.NotifyClose(closeListener)

	go func() {
		closeError := <- closeListener
		if closeError == nil {
			log.Print("No error close \n")
		} else {
			log.Printf("Error close: %v\n", closeError)
		}
	}()

	failOnError(dialError, "Failed to connect to RabbitMQ")
	defer connection.Close()

	channel, channelError := connection.Channel()

	failOnError(channelError, "Failed to open a channel")
	defer channel.Close()

	declareQueue(channel);
	declareExchange(channel);
	declareBinding(channel);

	channel.Publish("incoming_payloads", "", false, false, amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "application/json",
		ContentEncoding: "UTF-8",
		Body:            []byte(payload),
		DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
		Priority:        0,              // 0-9
		// a bunch of application/implementation-specific fields
	})
}

func declareQueue(channel *amqp.Channel) amqp.Queue {
	queue, queueError := channel.QueueDeclarePassive(
		"payloads_for_trips", // name
		true,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(queueError, "Failed to declare a queue")
	return queue
}

func declareExchange(channel *amqp.Channel) {
	exchangeError := channel.ExchangeDeclarePassive(
		"incoming_payloads", // name
		"fanout", //type
		true,   // durable
		false,   // delete when usused
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(exchangeError, "Failed to declare a exchange")
}

func declareBinding(channel *amqp.Channel) {
	bindingError := channel.QueueBind("payloads_for_trips", "", "incoming_payloads", false, nil)
	failOnError(bindingError, "Failed to declare a binding")
}