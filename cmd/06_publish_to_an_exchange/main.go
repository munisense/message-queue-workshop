package main

import (
	"context"
	"fmt"
	"github.com/munisense/syntax-workshop-2023/internal/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// Feel free to change this to your own exchange, don't forget to add a binding then!
const exchange = "results"

func main() {
	log := logrus.New()
	log.Level = logrus.DebugLevel

	conf := config.LoadConfig()

	// Connect to RabbitMQ server. The connection abstracts the socket connection, and takes care of protocol version negotiation and authentication and so on for us.
	url := fmt.Sprintf("%s://%s:%s@%s:%d%s", conf.MQProtocol, conf.MQUsername, conf.MQPassword, conf.MQHost, conf.MQPort, conf.MQVHost)
	log.Info("Connecting to the message queue...")
	conn, err := amqp.Dial(url)
	if err != nil {
		log.WithError(err).WithField("url", url).Fatal("failed to connect to the message queue")
	}
	log.Info("Connected to the message queue!")
	defer conn.Close() // "defer" = do this right before the application stops

	// Next we create a channel, which is where most of the API for getting things done resides.
	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).Fatal("failed to open a new channel")
	}

	// Generate a routing-key that we will use for our publishing
	routingKey := fmt.Sprintf("publish-%d", rand.Intn(9999))

	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, true, nil)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"exchange": exchange}).Fatal("failed to declare exchange")
	}

	// Now post some messages!
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // If we can't publis within 30 seconds -> cancel this
	defer cancel()
	go func() {
		for {
			// Generate a message:
			msg := amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
				ContentType:  "text/plain",
				Body:         []byte(fmt.Sprintf("Go Go AMQP! %d", rand.Intn(9999))),
			}
			err := ch.PublishWithContext(ctx, exchange, routingKey, false, false, msg)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"exchange": exchange, "message": msg}).Fatal("failed to publish to exchange")
			}
			log.WithFields(logrus.Fields{"body": string(msg.Body), "routingKey": routingKey}).Debug("published message to exchange")
			// Now wait a moment!
			time.Sleep(1 * time.Second)
		}
	}()

	// Block the main Goroutine, otherwise we would exit
	var forever chan struct{}
	<-forever
}
