package main

import (
	"github.com/streadway/amqp"
	"os"
)

func main() {
	url, _ := os.LookupEnv("RABBITMQ_URL")
	user, _ := os.LookupEnv("USER_NAME")
	pass, _ := os.LookupEnv("PASSWORD")
	v, _ := os.LookupEnv("V")
	target := "amqp://" + user + ":" + pass + "@" + url + "/" + v

	conn, err := amqp.Dial(target)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	body := "Hello World!"
	err = ch.Publish(
		"xch-fanout-webm-uc-pin-level-code", // exchange
		"",                                  // routing key
		false,                               // mandatory
		false,                               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		panic(err)
	}
}
