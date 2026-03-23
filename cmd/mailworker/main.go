package main

import (
	"encoding/json"
	"log"

	"github.com/Cinnamoon-dev/blue-gopher/internal/messaging/events"
	"github.com/Cinnamoon-dev/blue-gopher/internal/messaging/rabbitmq"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	env := config.NewEnv()
	conn, err := rabbitmq.NewConnection(env.RabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ connection: %s", err)
	}
	defer conn.Close()

	q, err := conn.Ch.QueueDeclare(
		"email",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	msgs, err := conn.Ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	go func() {
		authService := services.NewAuthService()
		mailService := services.NewMailService()

		for d := range msgs {
			var event events.EmailVerificationRequested

			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Failed at json unmarshal: %s", err)
				continue
			}

			emailToken, _ := authService.CreateToken(jwt.MapClaims{
				"email": event.Email,
			}, jwt.SigningMethodHS256, []byte(env.JwtKey))
			link := env.BackendUrl + "/mail/" + emailToken

			err = mailService.SendEmail(event.Email, "Email Verification", link)
			if err != nil {
				log.Printf("Failed to send mail: %s", err)
				continue
			}

			log.Printf("Mail sent to %s!", event.Email)
			d.Ack(false)
		}
	}()

	log.Print("Worker running...")
	var forever chan struct{}
	<-forever
}
