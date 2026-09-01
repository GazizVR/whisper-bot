package main

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

func (h *UpdateHandler) handle(
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

func main() {
	client := telegram.NewClient(
		"8863928065:AAF7s1uvrAOgymwvievxYxNQI7cmFcp9laE",
		// "http://localhost:8081",
		"https://api.telegram.org",
	)
	handler := NewUpdateHandler(client)
	if err := client.ListenUpdates(
		handler.handle,
		[]string{"meessage"},
	); err != nil {
		log.Fatalln(err)
	}
}
