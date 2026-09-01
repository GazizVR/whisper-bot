package main

import (
	"log"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/config"
	"github.com/gazizvr/whisper-bot/internal/handler"
)

func main() {
	config := config.Load()

	client := telegram.NewClient(
		config.Token,
		config.BaseURL,
	)
	handler := handler.NewUpdateHandler(client)
	if err := client.ListenUpdates(
		handler.Handle,
		[]string{"meessage"},
	); err != nil {
		log.Fatalln(err)
	}
}
