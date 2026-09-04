package handler

import (
	"fmt"
	"strings"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/persist"
)

const (
	LangSelectText      = "🗣 Выберите язык"
	OutFormatSelectText = "📤 Выберите формат вывода"
)

const (
	LangSelectAction      = "ls"
	OutFormatSelectAction = "ofs"
)

func (h *UpdateHandler) handleMedia(
	chatId, msgId int64,
	fileId string,
) {
	chat := h.RM.GetRequest(chatId)
	if chat == nil {
		chat = &persist.Request{}
	}

	mediaPath, err := h.Tg.DownloadFile(fileId, "tmp")
	if err != nil {
		h.sendErrMsg(err.Error(), chatId, msgId)
		return
	}
	chat.MediaPath = mediaPath
	chat.ReplyMsgId = &msgId
	h.RM.PutRequest(chatId, *chat)

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
	_, err = h.Tg.SendMessage(
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

func (h *UpdateHandler) handleCallback(
	id, data string,
	chatId, msgId int64,
) {
	h.Tg.AnswerCallbackQuery(id)
	if strings.Contains(data, LangSelectAction) {
		request := h.RM.GetRequest(chatId)
		if request == nil {
			h.editToErrMsg("Chat is empty handle callback", chatId, msgId)
			return
		}
		if request.MediaPath == nil {
			h.editToErrMsg("Chat.MediaPath is empty handle callback", chatId, msgId)
			return
		}
		lang := data[strings.LastIndex(data, "-")+1:]
		request.Language = &lang
		h.RM.PutRequest(chatId, *request)

		getOutFormat := func(lang string) string {
			return fmt.Sprintf("%s-%s", OutFormatSelectAction, lang)
		}
		buttons := [][]telegram.InlineButton{
			{
				{Text: "📄 .srt файл", Data: getOutFormat("srt")},
				{Text: "📝 текст", Data: getOutFormat("json")},
			},
		}
		markup := &telegram.InlineMarkup{
			Keyboard: buttons,
		}
		_, err := h.Tg.EditMessageText(
			chatId,
			msgId,
			OutFormatSelectText,
			markup,
		)
		if err != nil {
			h.editToErrMsg(err.Error(), chatId, msgId)
			return
		}
	}
	if strings.Contains(data, OutFormatSelectAction) {
		request := h.RM.GetRequest(chatId)
		if request == nil {
			h.editToErrMsg("Chat is empty handle callback", chatId, msgId)
			return
		}
		if request.MediaPath == nil {
			h.editToErrMsg("Chat.MediaPath is empty handle callback", chatId, msgId)
			return
		}
		if request.Language == nil {
			h.editToErrMsg("Chat.Language is empty handle callback", chatId, msgId)
			return
		}
		defer h.RM.RemoveRequest(chatId)
		outFormat := data[strings.LastIndex(data, "-")+1:]
		h.genAndSendSubtitles(
			chatId,
			msgId,
			request.ReplyMsgId,
			*request.MediaPath,
			*request.Language,
			outFormat,
		)
	}
}
