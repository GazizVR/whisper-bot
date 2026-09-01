package errors

import "errors"

var (
	ErrEmptyToken        = errors.New("Token is empty")
	ErrBaseURL           = errors.New("BaseURL is empty")
	ErrEmptyWhisperModel = errors.New("Whisper model path is empty")
)
