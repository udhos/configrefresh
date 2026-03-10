// Package main implements the tool.
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	version  = "0.0.1"
	cooldown = 5 * time.Second
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("version=%s\n", version)
		return
	}

	amqpURL := envString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
	app := envString("ROUTING_KEY", "springCloudBus")
	destinations := strings.Fields(envString("DESTINATION", "cartao-branco-parameters:** npc-regress-pismo:**"))

	// 3 minutes: must give time to parameters to load toggles from config-server
	interval := envDuration("INTERVAL", 3*time.Minute)

	for {
		send(amqpURL, app, destinations, interval)
		log.Printf("main: send exited, sleeping for %v", cooldown)
		time.Sleep(cooldown)
	}
}

func send(amqpURL, routingKey string, destinations []string, interval time.Duration) {

	var conn *amqp.Connection
	for {
		var err error
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			break
		}
		log.Printf("dial: %v", err)
		log.Printf("dial: sleeping for %v", cooldown)
		time.Sleep(cooldown)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("channel: %v", err)
		return
	}
	defer ch.Close()

	exchange := "springCloudBus"
	exchangeType := "topic"
	//queue := "config-event-queue"

	{
		log.Printf("got Channel, declaring Exchange: exchange=%s exchangeType=%s", exchange, exchangeType)
		err := ch.ExchangeDeclare(
			exchange,     // name of the exchange
			exchangeType, // type
			true,         // durable
			false,        // delete when complete
			false,        // internal
			false,        // noWait
			nil,          // arguments
		)
		if err != nil {
			log.Printf("exchange declare: %v", err)
			return
		}
	}

	timestamp := time.Now().UnixMilli()

	for {
		for _, dst := range destinations {

			id := uuid.New()

			// {"type":"RefreshRemoteApplicationEvent","timestamp":1494514362123,"originService":"config-server:docker:8888","destinationService":"accountservice:**","id":"53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}
			body := fmt.Sprintf(`{"type":"RefreshRemoteApplicationEvent","timestamp":%d,"originService":"config-server:0:4fc186a89bffcc6337d46702011f4c0a","destinationService":"%s","id":"%s"}`,
				timestamp, dst, id)

			log.Printf("publish: body: %s", body)

			err := ch.Publish(
				exchange,   // exchange
				routingKey, // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(body),
				})
			if err != nil {
				log.Printf("publish: %v", err)
				return
			}
		}
		log.Printf("publish: sent, sleeping for %v", interval)
		time.Sleep(interval)
	}
}
