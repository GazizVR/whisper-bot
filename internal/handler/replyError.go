package handler

import (
	"context"
	"log"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

const ErrorText = "❌ Внутренняя ошибка, попробуйте снова"

func (h *UpdateHandler) editToErrMsg(
	ctx context.Context,
	errText string,
	chatId, msgId int64,
) {
	log.Println(errText)
	h.bot.EditMessageText(
		ctx,
		telegrambot.EditMessageTextParams{
			ChatID:    chatId,
			MessageID: msgId,
			Text:      ErrorText,
		},
	)
}

func (h *UpdateHandler) sendErrMsg(
	ctx context.Context,
	errStr string,
	chatId, msgId int64,
) {
	log.Println(errStr)
	h.bot.SendMessage(
		ctx,
		telegrambot.SendMessageParams{
			ChatID: chatId,
			Text:   ErrorText,
			ReplyParams: &telegrambot.ReplyParameters{
				MessageID: &msgId,
			},
		},
	)
}
