package handler

import "os"

const WaitText = "⏳ Подождите, обрабатывается..."

func (h *UpdateHandler) genAndSendSubtitles(
	chatId, msgId int64,
	replyMsgId *int64,
	mediaPath, language, outFormat string,
) {
	defer os.Remove(mediaPath)
	msg, err := h.Tg.EditMessageText(
		chatId,
		msgId,
		WaitText,
		nil,
	)
	if err != nil {
		h.editToErrMsg(err.Error(), chatId, msgId)
		return
	}

	filePath, err := h.Trb.ToSrt(mediaPath, &language, outFormat)
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
		replyMsgId,
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
