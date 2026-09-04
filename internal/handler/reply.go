package handler

import (
	"log"
	"os"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/persist"
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

func (h *UpdateHandler) genAndSendSubtitles(
	chatId, msgId, replyMsgId int64,
	mediaPath, language string,
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

	filePath, err := h.Trb.ToSrt(mediaPath, &language)
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
		&replyMsgId,
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

const LangSelectText = "🗣 Выберите язык"

const (
	LangSelect = "ls"
)

func (h *UpdateHandler) handleMedia(
	chatId, msgId int64,
	fileId string,
) {
	chat := h.CM.GetChat(chatId)
	if chat == nil {
		chat = &persist.Chat{}
	}

	if chat.MediaPath == nil {
		mediaPath, err := h.Tg.DownloadFile(fileId, "tmp")
		if err != nil {
			h.sendErrMsg(err.Error(), chatId, msgId)
			return
		}
		chat.MediaPath = mediaPath
	}
	if chat.Language == nil {
		buttons := [][]telegram.InlineButton{
			{
				{Text: "🇺🇸 English", Data: "en"},
				{Text: "🇷🇺 Русский", Data: "ru"},
			},
		}
		markup := &telegram.InlineMarkup{
			Keyboard: buttons,
		}
		_, err := h.Tg.SendMessage(
			chatId,
			LangSelectText,
			markup,
			&msgId,
		)
		if err != nil {
			h.sendErrMsg(err.Error(), chatId, msgId)
			return
		}
	}
	// h.genAndSendSubtitles(
	// 	chatId,
	// 	msgId,
	// 	fileId,
	// )
}
