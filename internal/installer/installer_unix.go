//go:build !windows

package installer

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"pgen/internal/i18n"
)

type UnixInstaller struct {
	installPath string
	isSystem    bool
}

func newUnixInstaller(_ *i18n.Messages) *UnixInstaller {
	// Проверяем права root
	if os.Geteuid() == 0 {
		// Системная установка
		return &UnixInstaller{
			installPath: "/usr/local/bin",
			isSystem:    true,
		}
	} else {
		// Пользовательская установка
		homeDir, _ := os.UserHomeDir()
		localBin := filepath.Join(homeDir, ".local", "bin")
		return &UnixInstaller{
			installPath: localBin,
			isSystem:    false,
		}
	}
}

func (u *UnixInstaller) Install(messages *i18n.Messages) error {
	// Создаем директорию для установки
	if err := os.MkdirAll(u.installPath, 0755); err != nil {
		return fmt.Errorf("%s: %v", messages.InstallErrorCreateDir, err)
	}

	// Получаем путь к текущему исполняемому файлу
	sourcePath, err := GetExecutablePath(messages)
	if err != nil {
		return fmt.Errorf("%s: %v", messages.InstallErrorGetExePath, err)
	}

	// Определяем целевой путь
	targetPath := filepath.Join(u.installPath, "pgen")

	// Копируем файл
	if err := copyFile(sourcePath, targetPath); err != nil {
		return fmt.Errorf("%s: %v", messages.InstallErrorCopyFile, err)
	}

	// Устанавливаем права на выполнение
	if err := os.Chmod(targetPath, 0755); err != nil {
		return fmt.Errorf("%s: %v", messages.InstallErrorSetPerms, err)
	}

	// Добавляем в PATH если нужно
	if !u.isSystem {
		if err := u.addToUserPath(messages); err != nil {
			return fmt.Errorf("%s: %v", messages.InstallErrorAddPath, err)
		}
	}

	return nil
}

func (u *UnixInstaller) IsInstalled() bool {
	targetPath := filepath.Join(u.installPath, "pgen")
	_, err := os.Stat(targetPath)
	return err == nil
}

func (u *UnixInstaller) GetInstallPath() string {
	return u.installPath
}

func (u *UnixInstaller) Uninstall(messages *i18n.Messages) error {
	// Удаляем исполняемый файл
	targetPath := filepath.Join(u.installPath, "pgen")

	// Проверяем, существует ли файл
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		// Файл уже удален
		return nil
	}

	// Для Unix систем мы можем попробовать использовать отложенное удаление
	// Создаем скрипт удаления
	scriptContent := fmt.Sprintf(`#!/bin/sh
sleep 1
rm -f "%s"
rmdir "%s" 2>/dev/null || true
`, targetPath, u.installPath)

	scriptPath := filepath.Join("/tmp", "pgen_uninstall.sh")
	if err := ioutil.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		// Если не можем создать скрипт, пробуем удалить напрямую
		if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("%s: %v", messages.UninstallError, err)
		}
	} else {
		// Запускаем скрипт в фоновом режиме
		cmd := exec.Command("/bin/sh", scriptPath)
		if err := cmd.Start(); err != nil {
			// Если не можем запустить скрипт, пробуем удалить напрямую
			if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("%s: %v", messages.UninstallError, err)
			}
		}
		// Удаляем скрипт после использования
		defer os.Remove(scriptPath)
	}

	// Удаляем записи из профилей shell
	if !u.isSystem {
		if err := u.removeFromUserPath(messages); err != nil {
			return fmt.Errorf("%s: %v", messages.UninstallError, err)
		}
	}

	return nil
}

func (u *UnixInstaller) addToUserPath(messages *i18n.Messages) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Определяем файлы профиля для обновления
	profileFiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".profile"),
	}

	pathLine := fmt.Sprintf(`export PATH="$PATH:%s"`, u.installPath)

	for _, profileFile := range profileFiles {
		if _, err := os.Stat(profileFile); err == nil {
			if err := u.addPathToFile(profileFile, pathLine, messages); err != nil {
				continue // Пробуем следующий файл
			}
		}
	}

	return nil
}

func (u *UnixInstaller) removeFromUserPath(messages *i18n.Messages) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Определяем файлы профиля для обновления
	profileFiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".profile"),
	}

	pathLine := fmt.Sprintf(`export PATH="$PATH:%s"`, u.installPath)

	for _, profileFile := range profileFiles {
		if _, err := os.Stat(profileFile); err == nil {
			if err := u.removePathFromFile(profileFile, pathLine, messages); err != nil {
				continue // Пробуем следующий файл
			}
		}
	}

	return nil
}

func (u *UnixInstaller) addPathToFile(filename, pathLine string, messages *i18n.Messages) error {
	// Проверяем, есть ли уже эта строка в файле
	if u.pathExistsInFile(filename, u.installPath) {
		return nil // Путь уже добавлен
	}

	// Добавляем строку в конец файла
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\n# %s\n%s\n", messages.InstallProfileComment, pathLine))
	return err
}

func (u *UnixInstaller) removePathFromFile(filename, pathLine string, messages *i18n.Messages) error {
	// Читаем содержимое файла
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Преобразуем в строку и разбиваем на строки
	lines := strings.Split(string(content), "\n")

	// Создаем новый срез для строк без удаляемых
	var newLines []string
	profileComment := fmt.Sprintf("# %s", messages.InstallProfileComment)

	for _, line := range lines {
		// Пропускаем строки, связанные с установкой PGen
		if line == profileComment || line == pathLine || strings.Contains(line, u.installPath) {
			continue
		}
		newLines = append(newLines, line)
	}

	// Записываем обновленное содержимое обратно в файл
	return ioutil.WriteFile(filename, []byte(strings.Join(newLines, "\n")), 0644)
}

func (u *UnixInstaller) pathExistsInFile(filename, path string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, path) && strings.Contains(line, "PATH") {
			return true
		}
	}

	return false
}

// copyFile копирует файл из source в destination
func copyFile(source, destination string) error {
	// Используем системную команду cp для сохранения прав доступа
	cmd := exec.Command("cp", source, destination)
	return cmd.Run()
}

// Стабы для Windows функций (не вызываются на Unix)
func newWindowsInstaller(messages *i18n.Messages) Installer {
	panic(messages.InstallPanicWindowsFunc)
}

func isWindowsElevated(messages *i18n.Messages) bool {
	panic(messages.InstallPanicUnixFunc)
}
