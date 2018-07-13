package main

import (
	"context"
	"fmt"
	"log"

	"github.com/olivere/elastic"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func storeMessage(message string) {
	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
	}
	_, err = client.Index().
		Index("message").
		Type("doc").
		//Id("1").
		BodyJson(message).
		Refresh("wait_for").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
}

func readMessages() {
	conn, err := amqp.Dial("amqp://thedude:opinion@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"messages", // name
		false,      // durable
		false,      // delete when usused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			log.Printf("Saving to elasticsearch")
			storeMessage(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main() {
	readMessages()
}
