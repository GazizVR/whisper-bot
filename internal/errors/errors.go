package errors

import "errors"

var (
	ErrEmptyToken             = errors.New("Token is empty")
	ErrEmptyWhisperAPIBaseURL = errors.New("Whisper api base URL is empty")
	ErrFFMPEGNotFound         = errors.New("ffmpeg bin not found")
)
