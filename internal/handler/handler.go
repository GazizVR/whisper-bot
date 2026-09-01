package handler

import (
	"log"

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

func (h *UpdateHandler) Handle(
	upd telegram.Update,
) {
	if upd.Message != nil {
		if upd.Message.Video != nil {
			filePath, err := h.Tg.DownloadFile(upd.Message.Video.Id, "tmp")
			if err != nil {
				log.Println(err.Error())
				return
			}
			h.Tg.SendMessage(
				upd.Message.Chat.Id,
				*filePath,
				nil,
				nil,
			)
		}
	}
}
