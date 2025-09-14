//go:build !windows

package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func readPasswordWithStars(messages *InputMessages) (string, error) {
	fd := int(syscall.Stdin)

	if !term.IsTerminal(fd) {
		reader := bufio.NewReader(os.Stdin)
		password, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(password), nil
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Print("🔑 ")
		password, err := term.ReadPassword(fd)
		if err != nil {
			return "", err
		}
		fmt.Println()
		return string(password), nil
	}
	defer term.Restore(fd, oldState)

	var password strings.Builder
	var b [1]byte

	for {
		n, err := os.Stdin.Read(b[:])
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue
		}

		char := b[0]
		switch char {
		case 10, 13:
			fmt.Print("\r\n")
			return password.String(), nil
		case 127, 8:
			if password.Len() > 0 {
				str := password.String()
				password.Reset()
				password.WriteString(str[:len(str)-1])
				fmt.Print("\b \b")
			}
		case 3:
			fmt.Print("\r\n")
			return "", fmt.Errorf("%s", messages.UserCanceled)
		case 27:
			fmt.Print("\r\n")
			return "", fmt.Errorf("%s", messages.InputCanceled)
		default:
			if char >= 32 && char <= 126 {
				password.WriteByte(char)
				fmt.Print("*")
			}
		}
	}
}
