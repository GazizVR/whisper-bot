package handler

import (
	"fmt"
	"os"
	"strings"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/persist"
)

const WaitText = "⏳ Подождите, обрабатывается..."

func (h *UpdateHandler) genAndSendSubtitles(
	chatId, msgId int64,
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
		&msg.Result.ReplyTo.Id,
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

const (
	LangSelectText = "🗣 Выберите язык"
)

const (
	LangSelectAction = "ls"
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
		h.CM.PutChat(chatId, *chat)
	}
	if chat.Language == nil {
		genLang := func(lang string) string {
			return fmt.Sprintf("%s-%s", LangSelectAction, lang)
		}
		buttons := [][]telegram.InlineButton{
			{
				{Text: "🇺🇸 English", Data: genLang("en")},
				{Text: "🇷🇺 Русский", Data: genLang("ru")},
				{Text: "🇺🇿 Uzbek", Data: genLang("uz")},
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
}

func (h *UpdateHandler) handleCallback(
	id, data string,
	chatId, msgId int64,
) {
	h.Tg.AnswerCallbackQuery(id)
	if strings.Contains(data, LangSelectAction) {
		chat := h.CM.GetChat(chatId)
		if chat == nil {
			h.editToErrMsg("Chat is empty handle callback", chatId, msgId)
			return
		}
		if chat.MediaPath == nil {
			h.editToErrMsg("Chat.MediaPath is empty handle callback", chatId, msgId)
			return
		}
		defer h.CM.RemoveChat(chatId)
		lang := data[strings.LastIndex(data, "-")+1:]
		h.genAndSendSubtitles(
			chatId,
			msgId,
			*chat.MediaPath,
			lang,
		)
	}
}
