package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic"
	"github.com/streadway/amqp"
)

type Configuration struct {
	RabbitmqConnectionString      string
	RabbitmqExchangeName          string
	RabbitmqQueueName             string
	ElasticsearchConnectionString string
}

func createConfiguration() Configuration {

	file, _ := os.Open("config.json")

	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(configuration)

	return configuration
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func buildChannel(exchangeName, rabbitmqConnectionString string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(rabbitmqConnectionString)
	if err != nil {
		return nil, err
	}
	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = amqpChan.ExchangeDeclare(exchangeName,
		"fanout",
		true,
		false,
		false,
		false, nil)
	if err != nil {
		return nil, err
	}

	// Clear amqp channel if connection to server is lost
	amqpErrorChan := make(chan *amqp.Error)
	amqpChan.NotifyClose(amqpErrorChan)
	go func(ec chan *amqp.Error) {
		for msg := range ec {
			log.Fatalf("Channel Cleanup %s\n", msg)
		}
	}(amqpErrorChan)

	return amqpChan, err
}

func CreateHash(message string) string {
	hasher := sha256.New()
	hasher.Write([]byte(message))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return hash
}

func storeMessage(message, elasticsearchConnectionString string) {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(elasticsearchConnectionString), elastic.SetSniff(false))
	if err != nil {
		// Handle error
	}
	hash := CreateHash(message)
	currentDate := time.Now()

	//format := "2015/01/01 12:10:30"
	log.Printf("Current date %s", currentDate.Format(time.UnixDate))

	_, err = client.Index().
		Index("message").
		Type("doc").
		Id(string(hash)).
		BodyJson(message).
		Refresh("wait_for").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
}

func readMessages(config Configuration) {
	log.Printf("Connection to %s", config.RabbitmqConnectionString)
	conn, err := amqp.Dial(config.RabbitmqConnectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	log.Printf("Sending to exchange: %s", config.RabbitmqExchangeName)
	ch, err := buildChannel(config.RabbitmqExchangeName, config.RabbitmqConnectionString)
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		config.RabbitmqQueueName, // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		config.RabbitmqExchangeName,
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue")

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
			storeMessage(string(d.Body), config.ElasticsearchConnectionString)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main() {
	config := createConfiguration()
	fmt.Println(config)

	readMessages(config)
}
