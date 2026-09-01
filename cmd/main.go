package main

import (
	"log"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/config"
	"github.com/gazizvr/whisper-bot/internal/ffmpeg"
	"github.com/gazizvr/whisper-bot/internal/handler"
	"github.com/gazizvr/whisper-bot/internal/whisper"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}
	convertor := ffmpeg.NewConvertor(config.FfmpegBinPath)
	transcriber := whisper.NewTranscriber(
		config.WhisperBinPath,
		config.WhisperModelPath,
		convertor,
	)

	client := telegram.NewClient(
		config.Token,
		config.BaseURL,
	)
	handler := handler.NewUpdateHandler(client, transcriber)
	if err := client.ListenUpdates(
		handler.Handle,
		[]string{"meessage"},
	); err != nil {
		log.Fatalln(err)
	}
}
