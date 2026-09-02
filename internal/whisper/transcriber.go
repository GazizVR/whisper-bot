package whisper

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
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
	language *string,
) (*os.File, error) {
	wavPath, err := t.convertor.ToWav(mediaPath)
	if err != nil {
		return nil, err
	}
	wavFile, err := os.Open(*wavPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		wavFile.Close()
		os.Remove(*wavPath)
	}()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("response_format", "srt"); err != nil {
		return nil, err
	}
	if language != nil {
		if err := writer.WriteField("language", *language); err != nil {
			return nil, err
		}
	}

	formFile, err := writer.CreateFormFile("file", *wavPath)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(formFile, wavFile); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	u, err := url.Parse(t.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = "inference"
	urlStr := u.String()

	req, err := http.NewRequest(http.MethodPost, urlStr, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

}
