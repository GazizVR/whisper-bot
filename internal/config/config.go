package config

import "os"

type Config struct {
	token          string
	baseURL        string
	ffmpegBinPath  string
	whisperBinPath string
}

func Load() *Config {
	token := os.Getenv("TOKEN")
	baseURL := os.Getenv("BASE_URL")
	ffmpegBinPath := os.Getenv("FFMPEG_BIN_PATH")
	whisperBinPath := os.Getenv("WHISPER_BIN_PATH")
	return &Config{
		token:          token,
		baseURL:        baseURL,
		ffmpegBinPath:  ffmpegBinPath,
		whisperBinPath: whisperBinPath,
	}
}
