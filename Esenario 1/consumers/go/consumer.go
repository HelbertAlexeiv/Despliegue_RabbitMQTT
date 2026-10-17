package main

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", os.Getenv("RABBIT_USER"), os.Getenv("RABBIT_PASS"), os.Getenv("RABBIT_HOST"))

	var conn *amqp.Connection
	var err error
	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Println("Esperando RabbitMQ...")
		time.Sleep(3 * time.Second)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	err = ch.ExchangeDeclare("anime_exchange", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	q, err := ch.QueueDeclare("cola_go", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	claves := []string{"anime.accion", "anime.romance"}
	for _, k := range claves {
		if err := ch.QueueBind(q.Name, k, "anime_exchange", false, nil); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Go escuchando:", claves)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for m := range msgs {
		log.Printf("[GO] (%s) %s", m.RoutingKey, string(m.Body))
	}
}