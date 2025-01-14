package rabbitmq

import (
	"awesomeProject1/services/grpc_auth"
	"awesomeProject1/services/rabbitmq/handler"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

// ConsumeMessages sets up a consumer for a given queue
func ConsumeMessages(ch *amqp.Channel, queueName string, wg *sync.WaitGroup, authService *grpc_auth.AuthServiceServer) {
	defer wg.Done()
	msgs, err := ch.Consume(
		queueName, // queue name
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Listening for messages on queue: %s", queueName)
	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
		switch msg.RoutingKey {
		case "auth.register":
			if err := messageHandler(msg.Body, queueName, authService); err != nil {
				log.Printf("Failed to process message: %v", err)
			}
		case "email.welcome":
			{
				if err := handler.HandleEmail(msg.Body, queueName, authService); err != nil {
					log.Printf("Failed to process message: %v", err)
				}
			}
		default:
			log.Println("Unknown routing key or queueName:", msg.RoutingKey)
		}
	}

}
