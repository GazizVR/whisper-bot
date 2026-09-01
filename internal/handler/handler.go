package handler

import (
	"log"
	"os"

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

const (
	StartText = "🔗 Отправьте медиа файл"
	WaitText  = "⏳ Подождите, обрабатывается..."
)

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
			msg, _ := h.Tg.SendMessage(
				upd.Message.Chat.Id,
				WaitText,
				nil,
				&upd.Message.Id,
			)
			filePath, err := h.Tg.DownloadFile(upd.Message.Video.Id, "tmp")
			if err != nil {
				log.Println(err.Error())
				return
			}

			file, err := h.Trb.ToSrt(*filePath, "ru")
			if err != nil {
				log.Println(err.Error())
				return
			}

			h.Tg.EditMessageText(
				msg.Result.Chat.Id,
				msg.Result.Id,
				file.Name(),
				nil,
			)
			file.Close()
			os.Remove(file.Name())
		}
	}
}
