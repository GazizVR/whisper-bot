package errors

import "errors"

var (
	ErrEmptyToken     = errors.New("Token is empty")
	ErrTgBaseURL      = errors.New("Telegram base URL is empty")
	ErrWhisperBaseURL = errors.New("Whisper base URL is empty")
)
