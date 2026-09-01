package whisper

type Transcriber struct {
	binPath   string
	modelPath string
}

func NewTranscriber(
	binPath, modelPath string,
) *Transcriber {
	return &Transcriber{
		binPath:   binPath,
		modelPath: modelPath,
	}
}

func (t *Transcriber) Transcribe(
	mediaPath string,
) (*string, error) {

}
