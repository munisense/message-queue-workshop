package main

import (
	"fmt"
	"github.com/munisense/syntax-workshop-2023/internal/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// Name of the queue to get messages from.
// Try changing this name to 'results' and see what happens ;)
const queue = "laeq"
const exchange = "results"

func main() {
	rand.Seed(time.Now().UnixNano()) // If you are using an older golang version <1.20 you need to initialize the random seed generator
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

	// Create (declare) a queue with a semi-random name, in order not to flood the server we set a time-to-live for any messages on the queue to 60 seconds
	exclusiveQueueName := fmt.Sprintf("%s-%d", "syntax-workshop", rand.Intn(9999))
	exclusiveQueue, err := ch.QueueDeclare(exclusiveQueueName, false, true, true, false, amqp.Table{"x-message-ttl": 60000})
	if err != nil {
		log.WithError(err).WithField("exclusiveQueueName", exclusiveQueueName).Fatal("failed to declare a queue")
	}
	log.WithField("exclusiveQueueName", exclusiveQueueName).Debug("exclusive queue declared")

	// Our queue is now created but no data is being sent to the queue yet
	// For that we need to "bind" our queue to an exchange (in this case "results")
	// We can choose what data we want by supplying a routing key:
	routingKey := "*.*.*.*.LAeq" // or use `#` for ALL data
	err = ch.QueueBind(exclusiveQueue.Name, routingKey, exchange, false, nil)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"exclusiveQueueName": exclusiveQueueName, "routingKey": routingKey, "exchange": exchange}).Fatal("failed to bind queue to an exchange")
	}
	log.WithFields(logrus.Fields{"exclusiveQueueName": exclusiveQueueName, "routingKey": routingKey, "exchange": exchange}).Debug("bound queue to exchange")

	// The rest of the code is equal to step 3 other than consuming from the newly created queue

	// Tell the server to deliver us the messages from the queue.
	messages, err := ch.Consume(
		exclusiveQueue.Name,
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
