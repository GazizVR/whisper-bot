package main

import (
	"log"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/config"
	"github.com/gazizvr/whisper-bot/internal/handler"
	"github.com/gazizvr/whisper-bot/internal/store"
	"github.com/gazizvr/whisper-bot/pkg/ffmpeg"
	"github.com/gazizvr/whisper-bot/pkg/whisper"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}
	convertor := ffmpeg.NewConvertor(config.FfmpegBinPath)
	transcriber := whisper.NewTranscriber(
		config.WhisperAPIBaseURL,
		convertor,
	)

	client := telegram.NewClient(
		config.Token,
		config.BotAPIBaseURL,
	)
	log.Println("The server is starting!")
	requestStore := store.NewRequestStore()

	handler := handler.NewUpdateHandler(client, transcriber, requestStore)
	if err := client.ListenUpdates(
		handler.Handle,
		[]string{"message"},
	); err != nil {
		log.Fatalln(err)
	}
}
