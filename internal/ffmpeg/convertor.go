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

func (c *Convertor) ConvertToFlac(
	mediaPath string,
) (*string, error) {
	binPath, err := exec.LookPath(c.binPath)
	if err != nil {
		return nil, err
	}

	newMediaPath := fmt.Sprint(
		mediaPath[:strings.LastIndex(mediaPath, ".")],
		".flac",
	)
	cmdArgs := []string{
		"-i", mediaPath,
		"-vn",
		newMediaPath,
	}
	cmd := exec.Command(binPath, cmdArgs...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	if _, err = cmd.Output(); err != nil {
		errStr := fmt.Sprintf("Err: %v\nStdErr: %v", err, stdErr)
		return nil, errors.New(errStr)
	}

	return &newMediaPath, nil
}
