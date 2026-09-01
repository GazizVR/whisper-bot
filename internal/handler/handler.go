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

	if _, err := h.Tg.DeleteMessage(
		msg.Result.Chat.Id,
		msg.Result.Id,
	); err != nil {
		h.sendErrorMessage(err.Error(), chatId, msgId)
		return
	}
	if _, err := h.Tg.SendDocument(
		chatId,
		*file,
		&msgId,
	); err != nil {
		h.sendErrorMessage(err.Error(), chatId, msgId)
		return
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
			h.genAndSendSubtitles(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Video.Id,
			)
		}
		if upd.Message.Audio != nil {
			h.genAndSendSubtitles(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Audio.Id,
			)
		}
		if upd.Message.Voice != nil {
			h.genAndSendSubtitles(
				upd.Message.Chat.Id,
				upd.Message.Id,
				upd.Message.Voice.Id,
			)
		}
	}
}
