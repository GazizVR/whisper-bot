package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	telegrambot "github.com/gazizvr/go-telegram-bot"
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
	ctx context.Context,
	chatId, msgId int64,
	fileId string,
) error {
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
	buttons := [][]telegrambot.InlineKeyboardButton{
		{
			{Text: "🇺🇸 English", CallbackData: genLang("en")},
			{Text: "🇷🇺 Русский", CallbackData: genLang("ru")},
			{Text: "🇺🇿 Uzbek", CallbackData: genLang("uz")},
		},
		{
			{Text: "🌐 Автоматически", CallbackData: genLang("auto")},
		},
	}
	markup := &telegrambot.InlineKeyboardMarkup{
		Keyboard: buttons,
	}
	_, err := h.bot.SendMessage(
		ctx,
		telegrambot.SendMessageParams{
			ChatID:      chatId,
			Text:        LangSelectText,
			ReplyMarkup: markup,
			ReplyParams: &telegrambot.ReplyParameters{
				MessageID: &msgId,
			},
		},
	)
	if err != nil {
		h.sendErrMsg(ctx, err.Error(), chatId, msgId)
		return err
	}
	return nil
}

func (h *UpdateHandler) handleCallback(
	ctx context.Context,
	id, data string,
	chatId, msgId int64,
) error {
	h.bot.AnswerCallbackQuery(ctx, telegrambot.AnswerCallbackQueryParams{CallbackQueryID: id})
	if strings.Contains(data, LangSelectAction) {
		request := h.store.Get(chatId)
		if request == nil {
			errText := "Chat is empty handle callback"
			h.editToErrMsg(ctx, errText, chatId, msgId)
			return errors.New(errText)
		}
		if request.FileId == nil {
			errText := "Chat.FileId is empty handle callback"
			h.editToErrMsg(ctx, errText, chatId, msgId)
			return errors.New(errText)
		}
		language := data[strings.LastIndex(data, "-")+1:]
		if language != "auto" {
			request.Language = &language
			h.store.Set(chatId, *request)
		}

		getOutFormat := func(lang string) string {
			return fmt.Sprintf("%s-%s", OutFormatSelectAction, lang)
		}
		buttons := [][]telegrambot.InlineKeyboardButton{
			{
				{Text: "📄 .srt файл", CallbackData: getOutFormat("srt")},
				{Text: "📝 текст", CallbackData: getOutFormat("json")},
			},
		}
		markup := &telegrambot.InlineKeyboardMarkup{
			Keyboard: buttons,
		}
		_, err := h.bot.EditMessageText(
			ctx,
			telegrambot.EditMessageTextParams{
				ChatID:      chatId,
				MessageID:   msgId,
				Text:        OutFormatSelectText,
				ReplyMarkup: markup,
			},
		)
		if err != nil {
			h.editToErrMsg(ctx, err.Error(), chatId, msgId)
			return err
		}
	}
	if strings.Contains(data, OutFormatSelectAction) {
		request := h.store.Get(chatId)
		if request == nil {
			errText := "Chat is empty handle callback"
			h.editToErrMsg(ctx, errText, chatId, msgId)
			return errors.New(errText)
		}
		if request.FileId == nil {
			errText := "Chat.FileId is empty handle callback"
			h.editToErrMsg(ctx, errText, chatId, msgId)
			return errors.New(errText)
		}
		defer h.store.Delete(chatId)
		outFormat := data[strings.LastIndex(data, "-")+1:]
		if err := h.genAndSendSubtitles(
			ctx,
			chatId,
			msgId,
			request.ReplyMsgId,
			*request.FileId,
			request.Language,
			outFormat,
		); err != nil {
			return err
		}
	}
	return nil
}
