package main

import (
	"context"
	"log"
	"log/slog"

	telegrambot "github.com/gazizvr/go-telegram-bot"
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

	bot, err := telegrambot.NewBot(
		config.Token,
		telegrambot.WithBaseURL(config.BotAPIBaseURL),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("The server is starting!")
	requestStore := store.NewRequestStore()

	handler := handler.NewUpdateHandler(bot, transcriber, requestStore)

	var lastUpdateID int64
	ctx := context.Background()
	for {
		upds, err := bot.GetUpdates(ctx, telegrambot.GetUpdatesParams{
			Offset:         lastUpdateID + 1,
			Limit:          100,
			Timeout:        45,
			AllowedUpdates: []string{"message"},
		})
		if err != nil {
			log.Fatal(err)
		}
		for _, u := range upds {
			lastUpdateID = u.ID
			go func() {
				if err := handler.Handle(ctx, u); err != nil {
					slog.ErrorContext(ctx, err.Error())
				}
			}()
		}
	}
}
