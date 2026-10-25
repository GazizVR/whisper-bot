package handler

import (
	"context"
	"errors"
	"strings"

	telegrambot "github.com/gazizvr/go-telegram-bot"
	"github.com/gazizvr/whisper-bot/internal/store"
	"github.com/gazizvr/whisper-bot/pkg/whisper"
)

type UpdateHandler struct {
	bot         *telegrambot.Bot
	transcriber *whisper.Transcriber
	store       *store.RequestStore
}

func NewUpdateHandler(
	bot *telegrambot.Bot,
	transcriber *whisper.Transcriber,
	store *store.RequestStore,
) *UpdateHandler {
	return &UpdateHandler{
		bot:         bot,
		transcriber: transcriber,
		store:       store,
	}
}

const StartText = "🔗 Отправьте медиа файл"

func (h *UpdateHandler) Handle(
	ctx context.Context,
	upd telegrambot.Update,
) error {
	if upd.Message != nil {
		if upd.Message.Text == "/start" {
			if _, err := h.bot.SendMessage(
				ctx,
				telegrambot.SendMessageParams{
					ChatID: upd.Message.Chat.ID,
					Text:   StartText,
				},
			); err != nil {
				return err
			}
		}
		if upd.Message.Video != nil {
			if err := h.handleMedia(
				ctx,
				upd.Message.Chat.ID,
				upd.Message.ID,
				upd.Message.Video.FileID,
			); err != nil {
				return err
			}
		}
		if upd.Message.Audio != nil {
			if err := h.handleMedia(
				ctx,
				upd.Message.Chat.ID,
				upd.Message.ID,
				upd.Message.Audio.FileID,
			); err != nil {
				return err
			}
		}
		if upd.Message.Voice != nil {
			if err := h.handleMedia(
				ctx,
				upd.Message.Chat.ID,
				upd.Message.ID,
				upd.Message.Voice.FileID,
			); err != nil {
				return err
			}
		}
		if upd.Message.VideoNote != nil {
			if err := h.handleMedia(
				ctx,
				upd.Message.Chat.ID,
				upd.Message.ID,
				upd.Message.VideoNote.FileID,
			); err != nil {
				return err
			}
		}
		if upd.Message.Document != nil {
			if strings.Contains(upd.Message.Document.MimeType, "video") ||
				strings.Contains(upd.Message.Document.MimeType, "audio") {
				if err := h.handleMedia(
					ctx,
					upd.Message.Chat.ID,
					upd.Message.ID,
					upd.Message.Document.FileID,
				); err != nil {
					return err
				}
			}
		}
	}
	if upd.CallbackQuery != nil {
		query := upd.CallbackQuery
		var message telegrambot.InaccessibleMessage
		switch v := query.Message.(type) {
		case telegrambot.InaccessibleMessage:
			message = v
		case telegrambot.Message:
			message.ID = v.ID
			message.Chat = v.Chat
		case nil:
			return errors.New("message is nil")
		default:
			return errors.New("unsupported message type")
		}
		if err := h.handleCallback(
			ctx,
			query.ID,
			query.Data,
			message.Chat.ID,
			message.ID,
		); err != nil {
			return err
		}
	}
	return nil
}
