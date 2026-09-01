package handler

import (
	"log"
	"os"
)

const (
	WaitText  = "⏳ Подождите, обрабатывается..."
	ErrorText = "❌ Внутренняя ошибка, попробуйте снова"
)

func (h *UpdateHandler) sendErrorMessage(
	errText string,
	chatId, msgId int64,
) {
	log.Println(errText)
	h.Tg.SendMessage(
		chatId,
		ErrorText,
		nil,
		&msgId,
	)
}

func (h *UpdateHandler) genAndSendSubtitles(
	chatId, msgId int64,
	fileId string,
) {
	msg, err := h.Tg.SendMessage(
		chatId,
		WaitText,
		nil,
		&msgId,
	)
	if err != nil {
		h.sendErrorMessage(err.Error(), chatId, msgId)
		return
	}

	filePath, err := h.Tg.DownloadFile(fileId, "tmp")
	if err != nil {
		h.sendErrorMessage(err.Error(), chatId, msgId)
		return
	}
	defer os.Remove(*filePath)

	file, err := h.Trb.ToSrt(*filePath, "en")
	if err != nil {
		h.sendErrorMessage(err.Error(), chatId, msgId)
		return
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	go h.Tg.DeleteMessage(
		msg.Result.Chat.Id,
		msg.Result.Id,
	)
	if _, err := h.Tg.SendDocument(
		chatId,
		*file,
		&msgId,
	); err != nil {
		h.sendErrorMessage(err.Error(), chatId, msgId)
		return
	}
}
