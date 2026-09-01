package config

import "os"

type Config struct {
	Token          string
	BaseURL        string
	FfmpegBinPath  string
	WhisperBinPath string
}

func Load() *Config {
	token := os.Getenv("TOKEN")
	baseURL := os.Getenv("BASE_URL")
	ffmpegBinPath := os.Getenv("FFMPEG_BIN_PATH")
	whisperBinPath := os.Getenv("WHISPER_BIN_PATH")
	return &Config{
		Token:          token,
		BaseURL:        baseURL,
		FfmpegBinPath:  ffmpegBinPath,
		WhisperBinPath: whisperBinPath,
	}
}
