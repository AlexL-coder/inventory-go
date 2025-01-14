package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

// PublishMessage sends a message to the specified exchange with a routing key
func PublishMessage(ch *amqp.Channel, exchange, routingKey string, message interface{}) error {

	// Serialize the email details into JSON
	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal email details: %v", err)
		return err
	}

	err = ch.Publish(
		exchange,   // exchange name
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Message published to exchange '%s' with routing key '%s': %s", exchange, routingKey, body)
	return nil
}
