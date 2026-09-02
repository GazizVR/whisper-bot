package errors

import "errors"

var (
	ErrEmptyToken          = errors.New("Token is empty")
	ErrEmptyTgBaseURL      = errors.New("Telegram base URL is empty")
	ErrEmptyWhisperBaseURL = errors.New("Whisper base URL is empty")
	ErrFFMPEGNotFound      = errors.New("ffmpeg bin not found")
)
