package handler

import "log"

const ErrorText = "❌ Внутренняя ошибка, попробуйте снова"

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

func (h *UpdateHandler) sendErrMsg(
	errStr string,
	chatId, msgId int64,
) {
	log.Println(errStr)
	h.Tg.SendMessage(
		chatId,
		ErrorText,
		nil,
		&msgId,
	)
}
