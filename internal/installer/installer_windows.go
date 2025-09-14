//go:build windows

package installer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"pgen/internal/i18n"
)

type WindowsInstaller struct {
	installPath string
}

func newWindowsInstaller(messages *i18n.Messages) *WindowsInstaller {
	// Пытаемся установить в Program Files, если есть права
	programFiles := os.Getenv("PROGRAMFILES")
	if programFiles == "" {
		programFiles = "C:\\Program Files"
	}
	installPath := filepath.Join(programFiles, "PGen")

	// Если нет прав администратора, используем локальную папку
	if !IsElevated(messages) {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		installPath = filepath.Join(localAppData, "Programs", "PGen")
	}

	return &WindowsInstaller{
		installPath: installPath,
	}
}

func (w *WindowsInstaller) Install(messages *i18n.Messages) error {
	// Создаем директорию для установки
	if err := os.MkdirAll(w.installPath, 0755); err != nil {
		return fmt.Errorf("%s: %v", messages.InstallErrorCreateDir, err)
	}

	// Получаем путь к текущему исполняемому файлу
	sourcePath, err := GetExecutablePath(messages)
	if err != nil {
		return fmt.Errorf("%s: %v", messages.InstallErrorGetExePath, err)
	}

	// Определяем целевой путь
	targetPath := filepath.Join(w.installPath, "pgen.exe")

	// Копируем файл
	if err := copyFile(sourcePath, targetPath); err != nil {
		return fmt.Errorf("%s: %v", messages.InstallErrorCopyFile, err)
	}

	// Добавляем в PATH
	if err := w.addToPath(messages); err != nil {
		return fmt.Errorf("%s: %v", messages.InstallErrorAddPath, err)
	}

	return nil
}

func (w *WindowsInstaller) IsInstalled() bool {
	targetPath := filepath.Join(w.installPath, "pgen.exe")
	_, err := os.Stat(targetPath)
	return err == nil
}

func (w *WindowsInstaller) GetInstallPath() string {
	return w.installPath
}

func (w *WindowsInstaller) Uninstall(messages *i18n.Messages) error {
	// Удаляем исполняемый файл
	targetPath := filepath.Join(w.installPath, "pgen.exe")

	// Проверяем, существует ли файл
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		// Файл уже удален
		return nil
	}

	// Для Windows мы не можем удалить исполняемый файл, который сейчас выполняется
	// Вместо этого мы создаем скрипт удаления и запускаем его отдельно
	scriptContent := fmt.Sprintf(`
@echo off
echo Waiting for PGen to close...
timeout /t 3 /nobreak >nul
echo Removing PGen files...
del "%s" >nul 2>nul
if exist "%s" (
    echo Failed to remove main executable, trying forced removal...
    takeown /f "%s" >nul 2>nul
    icacls "%s" /grant administrators:F /t >nul 2>nul
    del "%s" >nul 2>nul
)
rmdir "%s" >nul 2>nul
echo PGen uninstallation completed.
`, targetPath, targetPath, targetPath, targetPath, targetPath, w.installPath)

	scriptPath := filepath.Join(os.TempDir(), "pgen_uninstall.bat")
	if err := ioutil.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("%s: %v", messages.UninstallError, err)
	}

	// Запускаем скрипт удаления в фоновом режиме с правами администратора
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Start-Process -FilePath '%s' -Verb RunAs -WindowStyle Hidden", scriptPath))
	if err := cmd.Start(); err != nil {
		// Если не можем запустить с правами администратора, пробуем без них
		cmd = exec.Command("cmd", "/c", scriptPath)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("%s: %v", messages.UninstallError, err)
		}
	}

	// Удаляем из PATH
	if err := w.removeFromPath(messages); err != nil {
		return fmt.Errorf("%s: %v", messages.UninstallError, err)
	}

	return nil
}

func (w *WindowsInstaller) addToPath(messages *i18n.Messages) error {
	if IsElevated(messages) {
		// Системная установка - изменяем системную переменную PATH
		return w.addToSystemPath()
	} else {
		// Пользовательская установка - изменяем пользовательскую переменную PATH
		return w.addToUserPath()
	}
}

func (w *WindowsInstaller) addToSystemPath() error {
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`$env:PATH = [Environment]::GetEnvironmentVariable("PATH", "Machine"); if ($env:PATH -notlike "*%s*") { [Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";%s", "Machine") }`, w.installPath, w.installPath))

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func (w *WindowsInstaller) addToUserPath() error {
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`$env:PATH = [Environment]::GetEnvironmentVariable("PATH", "User"); if ($env:PATH -notlike "*%s*") { [Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";%s", "User") }`, w.installPath, w.installPath))

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func (w *WindowsInstaller) removeFromPath(messages *i18n.Messages) error {
	if IsElevated(messages) {
		// Системная установка - изменяем системную переменную PATH
		return w.removeFromSystemPath(messages)
	} else {
		// Пользовательская установка - изменяем пользовательскую переменную PATH
		return w.removeFromUserPath(messages)
	}
}

func (w *WindowsInstaller) removeFromSystemPath(_ *i18n.Messages) error {
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`$env:PATH = [Environment]::GetEnvironmentVariable("PATH", "Machine"); $newPath = ($env:PATH -split ";" | Where-Object { $_ -ne "%s" }) -join ";"; [Environment]::SetEnvironmentVariable("PATH", $newPath, "Machine")`, w.installPath))

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func (w *WindowsInstaller) removeFromUserPath(_ *i18n.Messages) error {
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`$env:PATH = [Environment]::GetEnvironmentVariable("PATH", "User"); $newPath = ($env:PATH -split ";" | Where-Object { $_ -ne "%s" }) -join ";"; [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")`, w.installPath))

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

// isWindowsElevated упрощенная проверка прав администратора на Windows
func isWindowsElevated(_ *i18n.Messages) bool {
	// Простая проверка через попытку создания файла в системной папке
	programFiles := os.Getenv("PROGRAMFILES")
	if programFiles == "" {
		return false
	}
	testFile := filepath.Join(programFiles, "pgen_admin_test")
	file, err := os.Create(testFile)
	if err != nil {
		return false // Нет прав администратора
	}
	file.Close()
	os.Remove(testFile)
	return true // Есть права администратора
}

// copyFile копирует файл из source в destination
func copyFile(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	return err
}

// Стаб для Unix функций (не вызывается на Windows)
func newUnixInstaller(messages *i18n.Messages) Installer {
	panic(messages.InstallPanicUnixFunc)
}
