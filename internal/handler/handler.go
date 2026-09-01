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
	ErrorText = "❌ Внутренняя ошибка, попробуйте снова"
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
			msg, err := h.Tg.SendMessage(
				upd.Message.Chat.Id,
				WaitText,
				nil,
				&upd.Message.Id,
			)
			if err != nil {
				log.Println(err.Error())
				h.Tg.SendMessage(
					upd.Message.Chat.Id,
					ErrorText,
					nil,
					&upd.Message.Id,
				)
				return
			}

			filePath, err := h.Tg.DownloadFile(upd.Message.Video.Id, "tmp")
			if err != nil {
				log.Println(err.Error())
				h.Tg.SendMessage(
					upd.Message.Chat.Id,
					ErrorText,
					nil,
					&upd.Message.Id,
				)
				return
			}

			file, err := h.Trb.ToSrt(*filePath, "ru")
			if err != nil {
				log.Println(err.Error())
				h.Tg.SendMessage(
					upd.Message.Chat.Id,
					ErrorText,
					nil,
					&upd.Message.Id,
				)
				return
			}
			os.Remove(*filePath)

			_, err = h.Tg.DeleteMessage(
				msg.Result.Chat.Id,
				msg.Result.Id,
			)
			_, err = h.Tg.SendDocument(
				upd.Message.Chat.Id,
				*file,
				&upd.Message.Id,
			)

			file.Close()
			os.Remove(file.Name())

			if err != nil {
				log.Println(err.Error())
				h.Tg.SendMessage(
					upd.Message.Chat.Id,
					ErrorText,
					nil,
					&upd.Message.Id,
				)
				return
			}
		}
	}
}
