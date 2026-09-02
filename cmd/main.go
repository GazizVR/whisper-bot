package main

import (
	"log"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/config"
	"github.com/gazizvr/whisper-bot/internal/handler"
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
		config.WhisperBaseURL,
		convertor,
	)

	client := telegram.NewClient(
		config.Token,
		config.TgBaseURL,
	)
	log.Println("The server is starting!")
	handler := handler.NewUpdateHandler(client, transcriber)
	if err := client.ListenUpdates(
		handler.Handle,
		[]string{"meessage"},
	); err != nil {
		log.Fatalln(err)
	}
}
