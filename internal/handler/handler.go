package handler

import (
	"strings"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/persist"
	"github.com/gazizvr/whisper-bot/pkg/whisper"
)

type UpdateHandler struct {
	Tg  *telegram.Client
	Trb *whisper.Transcriber
	CM  *persist.ChatManager
}

func NewUpdateHandler(
	client *telegram.Client,
	transcriber *whisper.Transcriber,
	chatManager *persist.ChatManager,
) *UpdateHandler {
	return &UpdateHandler{
		Tg:  client,
		Trb: transcriber,
		CM:  chatManager,
	}
}

const StartText = "🔗 Отправьте медиа файл"

func (h *UpdateHandler) Handle(
	upd telegram.Update,
) {
	if upd.Message != nil {
		if upd.Message.Text == "/start" {
			h.Tg.SendMessage(
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
