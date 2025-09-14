package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"pgen/internal/i18n"
)

// Installer интерфейс для установки приложения
type Installer interface {
	Install(messages *i18n.Messages) error
	IsInstalled() bool
	GetInstallPath() string
	Uninstall(messages *i18n.Messages) error
}

// SystemInstaller создает инсталлятор для текущей платформы
func NewSystemInstaller(messages *i18n.Messages) Installer {
	switch runtime.GOOS {
	case "windows":
		return newWindowsInstaller(messages)
	default:
		return newUnixInstaller(messages)
	}
}

// GetExecutablePath возвращает путь к текущему исполняемому файлу
func GetExecutablePath(messages *i18n.Messages) (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("%s: %v", messages.InstallErrorGetExePath, err)
	}

	// Разрешаем символические ссылки
	resolved, err := filepath.EvalSymlinks(executable)
	if err != nil {
		return executable, nil // Возвращаем оригинальный путь если не удается разрешить
	}

	return resolved, nil
}

// IsElevated проверяет, запущено ли приложение с правами администратора
func IsElevated(messages *i18n.Messages) bool {
	switch runtime.GOOS {
	case "windows":
		return isWindowsElevated(messages)
	default:
		return os.Geteuid() == 0
	}
}
