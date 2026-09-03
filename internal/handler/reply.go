package handler

import (
	"log"
	"os"
)

const (
	WaitText  = "⏳ Подождите, обрабатывается..."
	ErrorText = "❌ Внутренняя ошибка, попробуйте снова"
)

func (h *UpdateHandler) editToErrMsg(
	errText string,
	chatId, msgId int64,
) {
	log.Println(errText)
	h.Tg.EditMessageText(
		chatId,
		msgId,
		ErrorText,
		nil,
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
		log.Println(err.Error())
		h.Tg.SendMessage(
			chatId,
			ErrorText,
			nil,
			&msgId,
		)
		return
	}

	mediaPath, err := h.Tg.DownloadFile(fileId, "tmp")
	if err != nil {
		h.editToErrMsg(err.Error(), chatId, msg.Result.Id)
		return
	}
	defer os.Remove(*mediaPath)

	language := "en"
	filePath, err := h.Trb.ToSrt(*mediaPath, &language)
	if err != nil {
		h.editToErrMsg(err.Error(), chatId, msg.Result.Id)
		return
	}
	file, err := os.Open(*filePath)
	if err != nil {
		h.editToErrMsg(err.Error(), chatId, msg.Result.Id)
		return
	}
	defer func() {
		file.Close()
		os.Remove(*filePath)
	}()

	if _, err := h.Tg.SendDocument(
		chatId,
		*file,
		&msgId,
	); err != nil {
		h.editToErrMsg(err.Error(), chatId, msg.Result.Id)
		return
	} else {
		h.Tg.DeleteMessage(
			msg.Result.Chat.Id,
			msg.Result.Id,
		)
	}
}
