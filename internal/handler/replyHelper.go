package handler

import (
	"context"
	"os"

	telegrambot "github.com/gazizvr/go-telegram-bot"
)

const WaitText = "⏳ Подождите, обрабатывается..."

func (h *UpdateHandler) genAndSendSubtitles(
	ctx context.Context,
	chatID, msgID int64,
	replyMsgID *int64,
	fileID string,
	language *string,
	outFormat string,
) error {
	msg, err := h.bot.EditMessageText(
		ctx,
		telegrambot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: msgID,
			Text:      WaitText,
		},
	)
	if err != nil {
		h.editToErrMsg(ctx, err.Error(), chatID, msgID)
		return err
	}

	file, err := h.bot.GetFile(ctx, telegrambot.GetFileParams{FileID: fileID})
	if err != nil {
		h.sendErrMsg(ctx, err.Error(), chatID, msgID)
		return err
	}
	filepath := file.Path
	defer os.Remove(filepath)

	switch outFormat {
	case "srt":
		fileName, err := h.transcriber.ToSrt(filepath, language)
		if err != nil {
			h.editToErrMsg(ctx, err.Error(), chatID, msg.ID)
			return err
		}
		defer os.Remove(fileName)

		if _, err := h.bot.SendDocument(
			ctx,
			telegrambot.SendDocumentParams{
				ChatID: chatID,
				Document: &telegrambot.LocalSource{
					Filepath: fileName,
				},
				ReplyParams: &telegrambot.ReplyParameters{
					MessageID: replyMsgID,
				},
			},
		); err != nil {
			h.editToErrMsg(ctx, err.Error(), chatID, msg.ID)
			return err
		} else {
			h.bot.DeleteMessage(
				ctx,
				telegrambot.DeleteMessageParams{
					ChatID:    msg.Chat.ID,
					MessageID: msg.ID,
				},
			)
		}
	case "json":
		resp, err := h.transcriber.ToJson(filepath, language)
		if err != nil {
			h.editToErrMsg(ctx, err.Error(), chatID, msg.ID)
			return err
		}
		if _, err := h.bot.EditMessageText(
			ctx,
			telegrambot.EditMessageTextParams{
				ChatID:    chatID,
				MessageID: msgID,
				Text:      resp.Text,
			},
		); err != nil {
			h.editToErrMsg(ctx, err.Error(), chatID, msg.ID)
			return err
		}
	}
	return nil
}
