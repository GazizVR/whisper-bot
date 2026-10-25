package config

import (
	"os"
	"os/exec"
	"strings"

	"github.com/gazizvr/whisper-bot/internal/errors"
)

type Config struct {
	Token             string
	BotAPIBaseURL     string
	WhisperAPIBaseURL string
	FfmpegBinPath     string
}

const defaultBotAPIBaseURL string = "https://api.telegram.org"

func Load() (*Config, error) {
	token := os.Getenv("TOKEN")
	if len(strings.TrimSpace(token)) < 1 {
		return nil, errors.ErrEmptyToken
	}
	botAPIBaseURL := os.Getenv("BOT_API_BASE_URL")
	if len(strings.TrimSpace(botAPIBaseURL)) < 1 {
		botAPIBaseURL = defaultBotAPIBaseURL
	}
	whisperAPIBaseURL := os.Getenv("WHISPER_API_BASE_URL")
	if len(strings.TrimSpace(whisperAPIBaseURL)) < 1 {
		return nil, errors.ErrEmptyWhisperAPIBaseURL
	}
	ffmpegBinPath := os.Getenv("FFMPEG_BIN_PATH")
	if len(strings.TrimSpace(ffmpegBinPath)) < 1 {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			return nil, errors.ErrFFMPEGNotFound
		}
		ffmpegBinPath = "ffmpeg"
	}
	return &Config{
		Token:             token,
		BotAPIBaseURL:     botAPIBaseURL,
		WhisperAPIBaseURL: whisperAPIBaseURL,
		FfmpegBinPath:     ffmpegBinPath,
	}, nil
}
