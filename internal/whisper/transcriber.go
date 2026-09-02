package whisper

import (
	"bytes"
	"mime/multipart"
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
	wavPath, err := t.convertor.ToWav(mediaPath)
	if err != nil {
		return nil, err
	}

	var byteWriter bytes.Buffer
	writer := multipart.NewWriter(&byteWriter)
	writer.WriteField("response_format", "srt")
	writer.Close()

	return nil, nil
}
