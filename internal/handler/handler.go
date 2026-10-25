package handler

import (
	"strings"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/store"
	"github.com/gazizvr/whisper-bot/pkg/whisper"
)

type UpdateHandler struct {
	bot         *telegram.Client
	transcriber *whisper.Transcriber
	store       *store.RequestStore
}

func NewUpdateHandler(
	bot *telegram.Client,
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
	upd telegram.Update,
) {
	if upd.Message != nil {
		if upd.Message.Text == "/start" {
			h.bot.SendMessage(
				upd.Message.Chat.Id,
				StartText,
				nil,
				nil,
			)
		}
		if upd.Message.Video != nil {
			h.handleMedia(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Video.Id,
			)
		}
		if upd.Message.Audio != nil {
			h.handleMedia(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Audio.Id,
			)
		}
		if upd.Message.Voice != nil {
			h.handleMedia(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Voice.Id,
			)
		}
		if upd.Message.VideoNote != nil {
			h.handleMedia(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.VideoNote.Id,
			)
		}
		if upd.Message.Document != nil {
			if strings.Contains(upd.Message.Document.MimeType, "video") ||
				strings.Contains(upd.Message.Document.MimeType, "audio") {
				h.handleMedia(
					upd.Message.Chat.Id,
					upd.Message.Id,
					upd.Message.Document.Id,
				)
			}
		}
	}
	if upd.Callback != nil {
		query := upd.Callback
		h.handleCallback(
			query.Id,
			query.Data,
			query.Message.Chat.Id,
			query.Message.Id,
		)
	}
}
