package main

import (
	"fmt"
	"github.com/munisense/message-queue-workshop/internal/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// Name of the queue to get messages from.
// Try changing this name to 'results' and see what happens ;)
const queue = "laeq"

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

	// Tell the server to deliver us the messages from the queue.
	messages, err := ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.WithError(err).Fatal("failed to create a consumer")
	}

	// Fire off a Goroutine that will wait and receive messages from the queue!
	go func() {
		for d := range messages {
			log.WithFields(logrus.Fields{
				"RoutingKey": d.RoutingKey,
				"Body":       string(d.Body),
			}).Debug("Got a message from the queue")
		}
	}()

	// Block the main Goroutine, otherwise we would exit
	var forever chan struct{}
	<-forever
}
