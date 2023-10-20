// Ok I lied, they are not queues but channels :)
package main

import (
	"encoding/json"
	"fmt"
	"github.com/munisense/syntax-workshop-2023/internal/pkg/config"
	"github.com/munisense/syntax-workshop-2023/internal/pkg/message"
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
	ch.Qos(1, 0, false) // Only allow 1 "un-acked" message
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

	// Tell the server to deliver us the messages from the queue.
	messages, err := ch.Consume(
		exclusiveQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.WithError(err).Fatal("failed to create a consumer")
	}

	// Fire off a Goroutine that will wait and receive messages from the queue!
	// This channel will hold a maximum of 10 entries, because our parsing (below) is really slow this channel will get full
	// Once that happens the code will be blocked and you can see the messages filling the rabbitmq queue
	channel := make(chan string, 10)
	go func() {
		for d := range messages {
			log.WithFields(logrus.Fields{
				"RoutingKey": d.RoutingKey,
				"Body":       string(d.Body),
			}).Debug("Got a message from the queue")
			// Also send the body (string) to our internal channel
			channel <- string(d.Body)
			d.Ack(false)
		}
	}()

	// But: in this case we also want to parse the messages. Perhaps that is an API call, or it takes some time. You don't want to handle this in the main thread
	// however how to you get your data in the parser? We can use golang channels (sort of queues but without routing keys) for that
	// see line 79 for the channel creation and 87 for the posting (publishing) on that channel

	go func() {
		for body := range channel {
			// Let's assume that this is an expensive operation (but here we just sleep)
			time.Sleep(10 * time.Second)
			m := message.Message{}
			err := json.Unmarshal([]byte(body), &m)
			if err != nil {
				log.WithError(err).WithField("message", m).Fatal("failed to unmarshall message")
			}
			log.WithFields(logrus.Fields{"metadata": m.Meta, "results": m.Results, "amountOfResults": len(m.Results)}).Debug("parsed message")
			for _, r := range m.Results {
				log.WithFields(logrus.Fields{"timestamp": r.Timestamp, "value": r.Value}).Debug("parsed message results")
			}
		}
	}()

	// Block the main Goroutine, otherwise we would exit
	var forever chan struct{}
	<-forever
}
