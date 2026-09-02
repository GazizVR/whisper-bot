package config

import (
	"os"
	"strings"

	"github.com/gazizvr/whisper-bot/internal/errors"
)

type Config struct {
	Token          string
	TgBaseURL      string
	WhisperBaseURL string
	FfmpegBinPath  string
}

func Load() (*Config, error) {
	token := os.Getenv("TOKEN")
	if len(strings.TrimSpace(token)) < 1 {
		return nil, errors.ErrEmptyToken
	}
	tgBaseURL := os.Getenv("TG_BASE_URL")
	if len(strings.TrimSpace(tgBaseURL)) < 1 {
		return nil, errors.ErrTgBaseURL
	}
	whisperBaseURL := os.Getenv("WHISPER_BASE_URL")
	if len(strings.TrimSpace(whisperBaseURL)) < 1 {
		return nil, errors.ErrWhisperBaseURL
	}
	ffmpegBinPath := os.Getenv("FFMPEG_BIN_PATH")
	if len(strings.TrimSpace(ffmpegBinPath)) < 1 {
		ffmpegBinPath = "ffmpeg"
	}
	return &Config{
		Token:          token,
		TgBaseURL:      tgBaseURL,
		WhisperBaseURL: whisperBaseURL,
		FfmpegBinPath:  ffmpegBinPath,
	}, nil
}
