package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

// Exchange Names
const (
	AuthExchange = "auth_exchange"
)

// Queue Names
const (
	AuthRegisterQueue  = "auth_register_queue"
	AuthListUsersQueue = "auth_list_users_queue"
	EmailWelcomeQueue  = "email_welcome_queue"
)

// Routing Keys
const (
	AuthRegisterRoutingKey  = "auth.register"
	AuthListUsersRoutingKey = "auth.list_users"
	EmailWelcomeRoutingKey  = "email.welcome"
)

// InitializeRabbitMQ sets up exchanges and queues
func InitializeRabbitMQ(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare Exchange
	err = ch.ExchangeDeclare(
		AuthExchange, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	// Declare Queues
	_, err = ch.QueueDeclare(
		AuthRegisterQueue, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		AuthListUsersQueue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(
		EmailWelcomeQueue, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}

	// Bind Queues to the Exchange
	err = ch.QueueBind(
		AuthRegisterQueue,      // Queue name
		AuthRegisterRoutingKey, // Routing key (using constant)
		AuthExchange,           // Exchange (using constant)
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		AuthListUsersQueue,      // Queue name
		AuthListUsersRoutingKey, // Routing key (using constant)
		AuthExchange,            // Exchange (using constant)
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		EmailWelcomeQueue,      // Queue name
		EmailWelcomeRoutingKey, // Routing key (using constant)
		AuthExchange,           // Exchange (using constant)
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return nil, err
	}

	log.Println("RabbitMQ setup complete: Exchanges and Queues are ready.")
	return ch, nil
}
