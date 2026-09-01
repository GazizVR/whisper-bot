package config

import (
	"os"
	"strings"

	"github.com/gazizvr/whisper-bot/internal/errors"
)

type Config struct {
	Token            string
	BaseURL          string
	FfmpegBinPath    string
	WhisperBinPath   string
	WhisperModelPath string
}

func Load() (*Config, error) {
	token := os.Getenv("TOKEN")
	if len(strings.TrimSpace(token)) < 1 {
		return nil, errors.ErrEmptyToken
	}
	baseURL := os.Getenv("BASE_URL")
	if len(strings.TrimSpace(baseURL)) < 1 {
		return nil, errors.ErrBaseURL
	}
	ffmpegBinPath := os.Getenv("FFMPEG_BIN_PATH")
	if len(strings.TrimSpace(ffmpegBinPath)) < 1 {
		ffmpegBinPath = "ffmpeg"
	}
	whisperBinPath := os.Getenv("WHISPER_BIN_PATH")
	if len(strings.TrimSpace(whisperBinPath)) < 1 {
		whisperBinPath = "whisper-cli"
	}
	whisperModelPath := os.Getenv("WHISPER_MODEL_PATH")
	if len(strings.TrimSpace(whisperModelPath)) < 1 {
		return nil, errors.ErrEmptyWhisperModel
	}
	return &Config{
		Token:            token,
		BaseURL:          baseURL,
		FfmpegBinPath:    ffmpegBinPath,
		WhisperBinPath:   whisperBinPath,
		WhisperModelPath: whisperModelPath,
	}, nil
}
