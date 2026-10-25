package handler

import (
	"fmt"
	"strings"

	telegram "github.com/gazizvr/tg-bot-api/pkg"
	"github.com/gazizvr/whisper-bot/internal/store"
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
	chat := h.store.Get(chatId)
	if chat == nil {
		chat = &store.Request{}
	}

	chat.FileId = &fileId
	chat.ReplyMsgId = &msgId
	h.store.Set(chatId, *chat)

	genLang := func(lang string) string {
		return fmt.Sprintf("%s-%s", LangSelectAction, lang)
	}
	buttons := [][]telegram.InlineButton{
		{
			{Text: "🇺🇸 English", Data: genLang("en")},
			{Text: "🇷🇺 Русский", Data: genLang("ru")},
			{Text: "🇺🇿 Uzbek", Data: genLang("uz")},
		},
		{
			{Text: "🌐 Автоматически", Data: genLang("auto")},
		},
	}
	markup := &telegram.InlineMarkup{
		Keyboard: buttons,
	}
	_, err := h.bot.SendMessage(
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
	h.bot.AnswerCallbackQuery(id)
	if strings.Contains(data, LangSelectAction) {
		request := h.store.Get(chatId)
		if request == nil {
			h.editToErrMsg("Chat is empty handle callback", chatId, msgId)
			return
		}
		if request.FileId == nil {
			h.editToErrMsg("Chat.FileId is empty handle callback", chatId, msgId)
			return
		}
		language := data[strings.LastIndex(data, "-")+1:]
		if language != "auto" {
			request.Language = &language
			h.store.Set(chatId, *request)
		}

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
		_, err := h.bot.EditMessageText(
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
		request := h.store.Get(chatId)
		if request == nil {
			h.editToErrMsg("Chat is empty handle callback", chatId, msgId)
			return
		}
		if request.FileId == nil {
			h.editToErrMsg("Chat.FileId is empty handle callback", chatId, msgId)
			return
		}
		defer h.store.Delete(chatId)
		outFormat := data[strings.LastIndex(data, "-")+1:]
		h.genAndSendSubtitles(
			chatId,
			msgId,
			request.ReplyMsgId,
			*request.FileId,
			request.Language,
			outFormat,
		)
	}
}
