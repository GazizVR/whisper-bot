package handler

import (
	"log"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
)

type UpdateHandler struct {
	Cl *telegram.Client
}

func NewUpdateHandler(
	client *telegram.Client,
) *UpdateHandler {
	return &UpdateHandler{
		Cl: client,
	}
}

func (h *UpdateHandler) Handle(
	upd telegram.Update,
) {
	if upd.Message != nil {
		if upd.Message.Video != nil {
			filePath, err := h.Cl.DownloadFile(upd.Message.Video.Id, "tmp")
			if err != nil {
				log.Println(err.Error())
				return
			}
			h.Cl.SendMessage(
				upd.Message.Chat.Id,
				*filePath,
				nil,
				nil,
			)
		}
	}
}
