package handler

import (
	"strings"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/whisper"
)

type UpdateHandler struct {
	Tg  *telegram.Client
	Trb *whisper.Transcriber
}

func NewUpdateHandler(
	client *telegram.Client,
	transcriber *whisper.Transcriber,
) *UpdateHandler {
	return &UpdateHandler{
		Tg:  client,
		Trb: transcriber,
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
			h.genAndSendSubtitles(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Video.Id,
			)
		}
		if upd.Message.Audio != nil {
			h.genAndSendSubtitles(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Audio.Id,
			)
		}
		if upd.Message.Voice != nil {
			h.genAndSendSubtitles(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Voice.Id,
			)
		}
		if upd.Message.VideoNote != nil {
			h.genAndSendSubtitles(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.VideoNote.Id,
			)
		}
		if upd.Message.Document != nil {
			if strings.Contains(upd.Message.Document.MimeType, "video") ||
				strings.Contains(upd.Message.Document.MimeType, "audio") {
				h.genAndSendSubtitles(
					upd.Message.Chat.Id,
					upd.Message.Id,
					upd.Message.Document.Id,
				)
			}
		}
	}
}
