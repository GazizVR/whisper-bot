package whisper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gazizvr/whisper-bot/pkg/ffmpeg"
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

func (t *Transcriber) transcribe(
	mediaPath string,
	language *string,
	outFormat string,
) ([]byte, error) {
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
	if err := writer.WriteField("response_format", outFormat); err != nil {
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
	if resp.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("Error code: %d", resp.StatusCode)
		return nil, errors.New(errStr)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (t *Transcriber) ToSrt(
	mediaPath string,
	language *string,
) (*os.File, error) {
	respData, err := t.transcribe(mediaPath, language, "srt")
	if err != nil {
		return nil, err
	}

	fileName := fmt.Sprint(
		string(mediaPath)[:strings.LastIndex(mediaPath, ".")],
		".srt",
	)
	if err := os.WriteFile(fileName, respData, 0644); err != nil {
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (t *Transcriber) ToJson(
	mediaPath string,
	language *string,
) (*InferenceResp, error) {
	respData, err := t.transcribe(mediaPath, language, "json")
	if err != nil {
		return nil, err
	}
	var resp InferenceResp
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
