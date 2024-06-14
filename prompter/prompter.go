package prompter

import (
	"fmt"
	"io"
	"os"
)

type Prompter struct {
}

func New() *Prompter {
	return &Prompter{}
}

func (p *Prompter) Get(name string) (string, error) {
	filePath := fmt.Sprintf("prompter/prompts/%s.txt", name)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
