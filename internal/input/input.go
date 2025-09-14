package input

import (
	"bufio"
	"os"
	"strings"
)

type InputMessages struct {
	UserCanceled  string
	InputCanceled string
}

func ReadPasswordWithStarsAndMessages(messages *InputMessages) (string, error) {
	return readPasswordWithStars(messages)
}

func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
