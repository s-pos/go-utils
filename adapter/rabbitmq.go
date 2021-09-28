package adapter

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"

	_ "github.com/joho/godotenv/autoload"
)

func ConnectionAMQP() *amqp.Channel {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s%s", os.Getenv("RABBIT_USERNAME"), os.Getenv("RABBIT_PASSWORD"), os.Getenv("RABBIT_HOST"), os.Getenv("RABBIT_PORT"), os.Getenv("RABBIT_VH"))

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
