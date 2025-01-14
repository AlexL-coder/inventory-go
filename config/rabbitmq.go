package config

import (
	"awesomeProject1/services/rabbitmq"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQConn establishes a connection to RabbitMQ
func rabbitMQConn() (*amqp.Connection, error) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}
	log.Println("Connected to RabbitMQ")
	return conn, nil
}

func InitializeRabbitMQ(maxRetries int) (*amqp.Connection, *amqp.Channel, error) {
	// Retry logic for RabbitMQ connection
	var rabbitConn *amqp.Connection
	var err error
	for i := 0; i < maxRetries; i++ {
		rabbitConn, err = rabbitMQConn()
		if err == nil {
			log.Println("Successfully connected to RabbitMQ!")
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second) // Wait before retrying
	}

	if rabbitConn == nil {
		log.Fatalf("Could not establish RabbitMQ connection after %d attempts: %v", maxRetries, err)
		return nil, nil, err
	}

	// Retry logic for RabbitMQ channel initialization
	var rabbitCh *amqp.Channel
	for i := 0; i < maxRetries; i++ {
		rabbitCh, err = rabbitmq.InitializeRabbitMQ(rabbitConn)
		if err == nil {
			log.Println("RabbitMQ channel initialized successfully!")
			break
		}
		log.Printf("Failed to initialize RabbitMQ channel (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second) // Wait before retrying
	}

	if rabbitCh == nil {
		log.Fatalf("Could not initialize RabbitMQ channel after %d attempts: %v", maxRetries, err)
		return nil, nil, err
	}

	// Your application logic here
	log.Println("RabbitMQ setup is complete!")

	return rabbitConn, rabbitCh, nil
}
