package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Convertor struct {
	binPath string
}

func NewConvertor(
	binPath string,
) *Convertor {
	return &Convertor{
		binPath: binPath,
	}
}

func (c *Convertor) ToWav(
	mediaPath string,
) (*string, error) {
	binPath, err := exec.LookPath(c.binPath)
	if err != nil {
		return nil, err
	}

	newMediaPath := fmt.Sprint(
		mediaPath[:strings.LastIndex(mediaPath, ".")],
		".wav",
	)
	cmdArgs := []string{
		"-i", mediaPath,
		"-vn",
		"-ac 1",
		"-ar 16000",
		"-c:a pcm_s16le",
		newMediaPath,
	}
	cmd := exec.Command(binPath, cmdArgs...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	if _, err = cmd.Output(); err != nil {
		errStr := fmt.Sprintf(
			"Err: %s\nStdErr: %s",
			err.Error(),
			stdErr.String(),
		)
		return nil, errors.New(errStr)
	}

	return &newMediaPath, nil
}
