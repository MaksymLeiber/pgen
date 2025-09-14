package installer

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/MaksymLeiber/pgen/internal/i18n"
)

func TestNewSystemInstaller(t *testing.T) {
	// Тест создания системного установщика
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}

	installer := NewSystemInstaller(messages)
	if installer == nil {
		t.Fatal("NewSystemInstaller не должен возвращать nil")
	}

	// Проверяем, что установщик работает корректно на любой платформе
	path := installer.GetInstallPath()
	if path == "" {
		t.Error("Установщик должен возвращать валидный путь установки")
	}
}

func TestGetExecutablePath(t *testing.T) {
	// Тест получения пути к исполняемому файлу
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}

	path, err := GetExecutablePath(messages)
	if err != nil {
		t.Errorf("GetExecutablePath не должен возвращать ошибку: %v", err)
	}

	if path == "" {
		t.Error("GetExecutablePath должен возвращать непустой путь")
	}

	// Проверяем, что путь существует
	if _, err := os.Stat(path); err != nil {
		t.Errorf("Путь к исполняемому файлу должен существовать: %s, ошибка: %v", path, err)
	}

	// Проверяем, что это абсолютный путь
	if !filepath.IsAbs(path) {
		t.Errorf("GetExecutablePath должен возвращать абсолютный путь, получен: %s", path)
	}
}

func TestGetExecutablePathWithNilMessages(t *testing.T) {
	// Тест с nil сообщениями
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetExecutablePath не должен паниковать с nil messages: %v", r)
		}
	}()

	// Функция может работать с nil messages, но может вернуть ошибку
	_, _ = GetExecutablePath(nil)
}

func TestIsElevated(t *testing.T) {
	// Тест проверки прав администратора
	messages := &i18n.Messages{}

	// Просто проверяем, что функция не паникует
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("IsElevated не должен паниковать: %v", r)
		}
	}()

	result := IsElevated(messages)
	
	// Результат может быть true или false в зависимости от прав
	// Просто проверяем, что функция возвращает булево значение
	if result != true && result != false {
		t.Error("IsElevated должен возвращать булево значение")
	}
}

func TestIsElevatedWithNilMessages(t *testing.T) {
	// Тест с nil сообщениями
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("IsElevated не должен паниковать с nil messages: %v", r)
		}
	}()

	_ = IsElevated(nil)
}

func TestInstallerInterface(t *testing.T) {
	// Тест интерфейса Installer
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}

	installer := NewSystemInstaller(messages)

	// Проверяем, что все методы интерфейса существуют
	t.Run("GetInstallPath", func(t *testing.T) {
		path := installer.GetInstallPath()
		if path == "" {
			t.Error("GetInstallPath не должен возвращать пустую строку")
		}
		
		// Проверяем, что путь содержит ожидаемые компоненты
		if !strings.Contains(strings.ToLower(path), "pgen") {
			t.Errorf("Путь установки должен содержать 'pgen', получен: %s", path)
		}
	})

	t.Run("IsInstalled", func(t *testing.T) {
		// Проверяем, что метод не паникует
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("IsInstalled не должен паниковать: %v", r)
			}
		}()

		result := installer.IsInstalled()
		// Результат может быть true или false
		if result != true && result != false {
			t.Error("IsInstalled должен возвращать булево значение")
		}
	})
}

func TestWindowsInstallerSpecific(t *testing.T) {
	// Тесты специфичные для Windows
	if runtime.GOOS != "windows" {
		t.Skip("Тест только для Windows")
	}

	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}

	installer := NewSystemInstaller(messages)

	t.Run("InstallPath", func(t *testing.T) {
		path := installer.GetInstallPath()
		
		// Путь должен содержать либо Program Files, либо AppData
		if !strings.Contains(path, "Program Files") && !strings.Contains(path, "AppData") {
			t.Errorf("Путь установки Windows должен содержать 'Program Files' или 'AppData', получен: %s", path)
		}
		
		// Путь должен заканчиваться на PGen
		if !strings.HasSuffix(path, "PGen") {
			t.Errorf("Путь установки должен заканчиваться на 'PGen', получен: %s", path)
		}
	})

	t.Run("IsInstalled", func(t *testing.T) {
		// Проверяем логику IsInstalled
		result := installer.IsInstalled()
		
		// Проверяем, что результат соответствует существованию файла
		targetPath := filepath.Join(installer.GetInstallPath(), "pgen.exe")
		_, err := os.Stat(targetPath)
		expectedResult := err == nil
		
		if result != expectedResult {
			t.Errorf("IsInstalled() = %v, но файл %s существует: %v", result, targetPath, expectedResult)
		}
	})
}

func TestUnixInstallerSpecific(t *testing.T) {
	// Тесты специфичные для Unix
	if runtime.GOOS == "windows" {
		t.Skip("Тест только для Unix систем")
	}

	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}

	installer := NewSystemInstaller(messages)

	t.Run("InstallPath", func(t *testing.T) {
		path := installer.GetInstallPath()
		
		// Путь должен быть либо системным, либо пользовательским
		if os.Geteuid() == 0 {
			// Для root пользователя
			if path != "/usr/local/bin" {
				t.Errorf("Для root пользователя путь должен быть '/usr/local/bin', получен: %s", path)
			}
		} else {
			// Для обычного пользователя
			if !strings.Contains(path, ".local/bin") {
				t.Errorf("Для обычного пользователя путь должен содержать '.local/bin', получен: %s", path)
			}
		}
	})

	t.Run("IsInstalled", func(t *testing.T) {
		// Проверяем логику IsInstalled
		result := installer.IsInstalled()
		
		// Проверяем, что результат соответствует существованию файла
		targetPath := filepath.Join(installer.GetInstallPath(), "pgen")
		_, err := os.Stat(targetPath)
		expectedResult := err == nil
		
		if result != expectedResult {
			t.Errorf("IsInstalled() = %v, но файл %s существует: %v", result, targetPath, expectedResult)
		}
	})
}

func TestInstallerWithDifferentMessages(t *testing.T) {
	// Тест с разными сообщениями
	testMessages := []*i18n.Messages{
		{
			InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
			InstallErrorCreateDir:  "Ошибка создания директории",
			InstallErrorCopyFile:   "Ошибка копирования файла",
		},
		{
			InstallErrorGetExePath: "Error getting executable path",
			InstallErrorCreateDir:  "Error creating directory", 
			InstallErrorCopyFile:   "Error copying file",
		},
		{
			InstallErrorGetExePath: "Fehler beim Abrufen des ausführbaren Pfads",
			InstallErrorCreateDir:  "Fehler beim Erstellen des Verzeichnisses",
			InstallErrorCopyFile:   "Fehler beim Kopieren der Datei",
		},
	}

	for i, messages := range testMessages {
		t.Run(strings.Join([]string{"Messages", string(rune('1' + i))}, "_"), func(t *testing.T) {
			installer := NewSystemInstaller(messages)
			if installer == nil {
				t.Error("NewSystemInstaller не должен возвращать nil")
			}

			path := installer.GetInstallPath()
			if path == "" {
				t.Error("GetInstallPath не должен возвращать пустую строку")
			}

			// Проверяем, что IsInstalled не паникует
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("IsInstalled не должен паниковать: %v", r)
				}
			}()
			_ = installer.IsInstalled()
		})
	}
}

func TestInstallerPathValidation(t *testing.T) {
	// Тест валидации путей установки
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}

	installer := NewSystemInstaller(messages)
	path := installer.GetInstallPath()

	// Проверяем, что путь абсолютный
	if !filepath.IsAbs(path) {
		t.Errorf("Путь установки должен быть абсолютным, получен: %s", path)
	}

	// Проверяем, что путь не содержит недопустимые символы (кроме : на Windows)
	invalidChars := []string{"<", ">", "\"", "|", "?", "*"}
	if runtime.GOOS != "windows" {
		invalidChars = append(invalidChars, ":")
	}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			t.Errorf("Путь установки не должен содержать недопустимый символ '%s': %s", char, path)
		}
	}

	// Проверяем, что путь не пустой и не только пробелы
	if strings.TrimSpace(path) == "" {
		t.Error("Путь установки не должен быть пустым или содержать только пробелы")
	}
}

func TestInstallerConsistency(t *testing.T) {
	// Тест консистентности поведения установщика
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}

	// Создаем несколько экземпляров установщика
	installer1 := NewSystemInstaller(messages)
	installer2 := NewSystemInstaller(messages)

	// Проверяем, что они возвращают одинаковые пути
	path1 := installer1.GetInstallPath()
	path2 := installer2.GetInstallPath()

	if path1 != path2 {
		t.Errorf("Разные экземпляры установщика должны возвращать одинаковые пути: %s != %s", path1, path2)
	}

	// Проверяем, что IsInstalled возвращает одинаковые результаты
	result1 := installer1.IsInstalled()
	result2 := installer2.IsInstalled()

	if result1 != result2 {
		t.Errorf("Разные экземпляры установщика должны возвращать одинаковые результаты IsInstalled: %v != %v", result1, result2)
	}
}

func TestInstallerEdgeCases(t *testing.T) {
	// Тест граничных случаев
	t.Run("EmptyMessages", func(t *testing.T) {
		emptyMessages := &i18n.Messages{}
		installer := NewSystemInstaller(emptyMessages)
		
		if installer == nil {
			t.Error("NewSystemInstaller должен работать с пустыми сообщениями")
		}
		
		path := installer.GetInstallPath()
		if path == "" {
			t.Error("GetInstallPath должен возвращать путь даже с пустыми сообщениями")
		}
	})

	t.Run("UnicodeMessages", func(t *testing.T) {
		unicodeMessages := &i18n.Messages{
			InstallErrorGetExePath: "Ошибка получения пути 🔧",
			InstallErrorCreateDir:  "Ошибка создания папки 📁",
			InstallErrorCopyFile:   "Ошибка копирования 📋",
		}
		
		installer := NewSystemInstaller(unicodeMessages)
		if installer == nil {
			t.Error("NewSystemInstaller должен работать с Unicode сообщениями")
		}
	})

	t.Run("LongMessages", func(t *testing.T) {
		longMessage := strings.Repeat("Очень длинное сообщение об ошибке ", 100)
		longMessages := &i18n.Messages{
			InstallErrorGetExePath: longMessage,
			InstallErrorCreateDir:  longMessage,
			InstallErrorCopyFile:   longMessage,
		}
		
		installer := NewSystemInstaller(longMessages)
		if installer == nil {
			t.Error("NewSystemInstaller должен работать с длинными сообщениями")
		}
	})
}

// Бенчмарки для измерения производительности
func BenchmarkNewSystemInstaller(b *testing.B) {
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewSystemInstaller(messages)
	}
}

func BenchmarkGetExecutablePath(b *testing.B) {
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetExecutablePath(messages)
	}
}

func BenchmarkIsElevated(b *testing.B) {
	messages := &i18n.Messages{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsElevated(messages)
	}
}

func BenchmarkInstallerGetInstallPath(b *testing.B) {
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}
	installer := NewSystemInstaller(messages)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = installer.GetInstallPath()
	}
}

func BenchmarkInstallerIsInstalled(b *testing.B) {
	messages := &i18n.Messages{
		InstallErrorGetExePath: "Ошибка получения пути к исполняемому файлу",
	}
	installer := NewSystemInstaller(messages)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = installer.IsInstalled()
	}
}
