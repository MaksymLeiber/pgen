package input

import (
	"bufio"
	"os"
	"strings"

	"github.com/MaksymLeiber/pgen/internal/security"
)

type InputMessages struct {
	UserCanceled  string
	InputCanceled string
}

func ReadPasswordWithStarsAndMessages(messages *InputMessages) (*security.SecureString, error) {
	password, err := readPasswordWithStars(messages)
	if err != nil {
		return nil, err
	}
	
	// Создаем SecureString из введенного пароля
	securePassword := security.NewSecureString(password)
	
	// Очищаем обычную строку из памяти
	security.SecureWipe([]byte(password))
	
	return securePassword, nil
}

func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
