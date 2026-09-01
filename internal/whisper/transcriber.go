package whisper

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/gazizvr/whisper-bot/internal/ffmpeg"
)

type Transcriber struct {
	binPath   string
	modelPath string
	convertor *ffmpeg.Convertor
}

func NewTranscriber(
	binPath, modelPath string,
	convertor *ffmpeg.Convertor,
) *Transcriber {
	return &Transcriber{
		binPath:   binPath,
		modelPath: modelPath,
		convertor: convertor,
	}
}

func (t *Transcriber) ToSrt(
	mediaPath string,
	language string,
) (*os.File, error) {
	binPath, err := exec.LookPath(t.binPath)
	if err != nil {
		return nil, err
	}

	cnvMediaPath, err := t.convertor.ToFlac(mediaPath)
	if err != nil {
		return nil, err
	}
	defer os.Remove(*cnvMediaPath)

	cmdArgs := []string{
		"-m", t.modelPath,
		"-f", *cnvMediaPath,
		"-osrt",
	}
	if len(language) > 0 {
		cmdArgs = append(cmdArgs, "-l", language)
	}
	cmd := exec.Command(binPath, cmdArgs...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	if _, err = cmd.Output(); err != nil {
		errStr := fmt.Sprintf("Err: %v\nStdErr: %v", err, stdErr)
		return nil, errors.New(errStr)
	}

	outFilePath := fmt.Sprint(*cnvMediaPath, ".srt")
	outFile, err := os.Open(outFilePath)
	if err != nil {
		return nil, err
	}

	return outFile, nil
}
