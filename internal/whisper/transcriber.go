package whisper

import (
	"os"

	"github.com/gazizvr/whisper-bot/internal/ffmpeg"
)

type Transcriber struct {
	baseURL   string
	convertor *ffmpeg.Convertor
}

func NewTranscriber(
	baseURL string,
	convertor *ffmpeg.Convertor,
) *Transcriber {
	return &Transcriber{
		baseURL:   baseURL,
		convertor: convertor,
	}
}

func (t *Transcriber) ToSrt(
	mediaPath string,
	language string,
) (*os.File, error) {
	return nil, nil
}
