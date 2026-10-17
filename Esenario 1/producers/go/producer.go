package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchange = "anime_exchange"
const genero = "accion"

var titulos = []string{"Attack on Titan", "Jujutsu Kaisen", "Demon Slayer", "Chainsaw Man", "My Hero Academia"}
var eventos = []string{"Nuevo episodio", "Nuevo trailer", "Confirmada nueva temporada", "Estreno de pelicula"}

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
	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		t := titulos[rand.Intn(len(titulos))]
		e := eventos[rand.Intn(len(eventos))]
		body := fmt.Sprintf(`{"genero":"%s","titulo":"%s","evento":"%s","productor":"go","hora":"%s"}`,
			genero, t, e, time.Now().Format("15:04:05"))
		err := ch.PublishWithContext(context.Background(), exchange, "anime."+genero, false, false,
			amqp.Publishing{ContentType: "application/json", Body: []byte(body)})
		if err != nil {
			log.Println("Error:", err)
		}
		log.Printf("[PUBLICADO] anime.%s -> %s | %s", genero, t, e)
		time.Sleep(time.Duration(2+rand.Intn(4)) * time.Second)
	}
}