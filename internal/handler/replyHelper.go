package handler

import "os"

const WaitText = "⏳ Подождите, обрабатывается..."

func (h *UpdateHandler) genAndSendSubtitles(
	chatId, msgId int64,
	replyMsgId *int64,
	fileId string,
	language *string,
	outFormat string,
) {
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

	mediaPath, err := h.Tg.DownloadFile(fileId, "tmp")
	if err != nil {
		h.sendErrMsg(err.Error(), chatId, msgId)
		return
	}
	defer os.Remove(*mediaPath)

	switch outFormat {
	case "srt":
		file, err := h.Trb.ToSrt(*mediaPath, language)
		if err != nil {
			h.editToErrMsg(err.Error(), chatId, msg.Result.Id)
			return
		}
		defer func() {
			file.Close()
			os.Remove(file.Name())
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
	case "json":
		resp, err := h.Trb.ToJson(*mediaPath, language)
		if err != nil {
			h.editToErrMsg(err.Error(), chatId, msg.Result.Id)
			return
		}
		h.Tg.EditMessageText(
			chatId,
			msgId,
			resp.Text,
			nil,
		)
	}
}
