//go:build windows

package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/term"
)

const (
	enableEchoInput = 0x0004
	enableLineInput = 0x0002
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

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")

	handle := syscall.Handle(syscall.Stdin)

	var oldMode uint32
	r1, _, _ := procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&oldMode)))
	if r1 == 0 {
		fmt.Print("🔑 ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println()
		return string(password), nil
	}

	newMode := oldMode &^ (enableEchoInput | enableLineInput)
	r1, _, _ = procSetConsoleMode.Call(uintptr(handle), uintptr(newMode))
	if r1 == 0 {
		fmt.Print("🔑 ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println()
		return string(password), nil
	}

	defer func() {
		procSetConsoleMode.Call(uintptr(handle), uintptr(oldMode))
	}()

	var password strings.Builder
	var b [1]byte

	for {
		n, err := syscall.Read(handle, b[:])
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue
		}

		char := b[0]
		switch char {
		case 13:
			fmt.Print("\r\n")
			return password.String(), nil
		case 8:
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
