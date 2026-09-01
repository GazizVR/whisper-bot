package errors

import "errors"

var (
	ErrEmptyToken          = errors.New("Token is empty")
	ErrBaseURL             = errors.New("BaseURL is empty")
	ErrEmptyFfmpegBinPath  = errors.New("ffmpeg bin path is empty")
	ErrEmptyWhisperBinPath = errors.New("Whisper bin path is empty")
	ErrEmptyWhisperModel   = errors.New("Whisper model path is empty")
)
