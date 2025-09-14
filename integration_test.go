package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestMainIntegration тестирует основные сценарии использования CLI
func TestMainIntegration(t *testing.T) {
	// Собираем бинарный файл для тестирования
	binaryName := "pgen_test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// Создаем временную директорию для сборки
	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, binaryName)

	// Собираем приложение
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Не удалось собрать приложение: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
		expectedError  bool
		timeout        time.Duration
	}{
		{
			name:           "Показать версию",
			args:           []string{"--version"},
			expectedOutput: []string{"0.0.8"},
			expectedError:  false,
			timeout:        5 * time.Second,
		},
		{
			name:           "Показать помощь",
			args:           []string{"--help"},
			expectedOutput: []string{"pgen", "Использование:", "Флаги:"},
			expectedError:  false,
			timeout:        5 * time.Second,
		},
		{
			name:           "Показать помощь на русском",
			args:           []string{"--lang", "ru", "--help"},
			expectedOutput: []string{"pgen", "Использование:", "Флаги:"},
			expectedError:  false,
			timeout:        5 * time.Second,
		},
		{
			name:           "Показать информацию о программе",
			args:           []string{"--about"},
			expectedOutput: []string{"PGen", "безопасных паролей", "детерминированных"},
			expectedError:  false,
			timeout:        5 * time.Second,
		},
		{
			name:           "Показать информацию о программе на русском",
			args:           []string{"--lang", "ru", "--about"},
			expectedOutput: []string{"PGen", "безопасных паролей", "детерминированных"},
			expectedError:  false,
			timeout:        5 * time.Second,
		},
		{
			name:           "Неизвестный флаг",
			args:           []string{"--unknown-flag"},
			expectedOutput: []string{},
			expectedError:  true,
			timeout:        5 * time.Second,
		},
		{
			name:           "Неверная длина пароля",
			args:           []string{"--length", "0"},
			expectedOutput: []string{},
			expectedError:  true,
			timeout:        5 * time.Second,
		},
		{
			name:           "Слишком большая длина пароля",
			args:           []string{"--length", "200"},
			expectedOutput: []string{},
			expectedError:  true,
			timeout:        5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем команду с таймаутом
			cmd := exec.Command(binaryPath, tt.args...)
			
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Запускаем команду с таймаутом
			done := make(chan error, 1)
			go func() {
				done <- cmd.Run()
			}()

			select {
			case err := <-done:
				// Команда завершилась
				if tt.expectedError && err == nil {
					t.Errorf("Ожидалась ошибка, но команда выполнилась успешно")
				}
				if !tt.expectedError && err != nil {
					t.Errorf("Неожиданная ошибка: %v\nStderr: %s", err, stderr.String())
				}

				// Проверяем ожидаемый вывод
				output := stdout.String() + stderr.String()
				for _, expected := range tt.expectedOutput {
					if !strings.Contains(output, expected) {
						t.Errorf("Ожидаемый текст '%s' не найден в выводе:\n%s", expected, output)
					}
				}

			case <-time.After(tt.timeout):
				// Таймаут - убиваем процесс
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
				t.Errorf("Команда превысила таймаут %v", tt.timeout)
			}
		})
	}
}

// TestConfigIntegration тестирует работу с конфигурацией
func TestConfigIntegration(t *testing.T) {
	// Собираем бинарный файл
	binaryName := "pgen_test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Не удалось собрать приложение: %v", err)
	}

	// Устанавливаем временную домашнюю директорию
	originalHome := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		originalHome = os.Getenv("USERPROFILE")
	}
	
	testHome := t.TempDir()
	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", testHome)
	} else {
		os.Setenv("HOME", testHome)
	}
	
	defer func() {
		if runtime.GOOS == "windows" {
			os.Setenv("USERPROFILE", originalHome)
		} else {
			os.Setenv("HOME", originalHome)
		}
	}()

	tests := []struct {
		name        string
		commands    [][]string
		expectError bool
	}{
		{
			name: "Показать конфигурацию",
			commands: [][]string{
				{"config", "show"},
			},
			expectError: false,
		},
		{
			name: "Установить и получить значение конфигурации",
			commands: [][]string{
				{"config", "set", "default_length", "20"},
				{"config", "get", "default_length"},
			},
			expectError: false,
		},
		{
			name: "Установить язык",
			commands: [][]string{
				{"config", "set", "default_language", "ru"},
				{"config", "get", "default_language"},
			},
			expectError: false,
		},
		{
			name: "Неверное значение конфигурации",
			commands: [][]string{
				{"config", "set", "default_length", "invalid"},
			},
			expectError: true,
		},
		{
			name: "Неизвестный ключ конфигурации",
			commands: [][]string{
				{"config", "set", "unknown_key", "value"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, cmdArgs := range tt.commands {
				cmd := exec.Command(binaryPath, cmdArgs...)
				
				var stdout, stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err := cmd.Run()
				
				// Для последней команды в тесте проверяем ожидание ошибки
				if i == len(tt.commands)-1 {
					if tt.expectError && err == nil {
						t.Errorf("Ожидалась ошибка для команды %v", cmdArgs)
					}
					if !tt.expectError && err != nil {
						t.Errorf("Неожиданная ошибка для команды %v: %v\nStderr: %s", 
							cmdArgs, err, stderr.String())
					}
				}
			}
		})
	}
}

// TestLanguageIntegration тестирует поддержку разных языков
func TestLanguageIntegration(t *testing.T) {
	binaryName := "pgen_test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Не удалось собрать приложение: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedTexts  []string
		unexpectedTexts []string
	}{
		{
			name:            "Английский язык через флаг",
			args:            []string{"--lang", "en", "--help"},
			expectedTexts:   []string{"Usage:", "Flags:", "Generate"},
			unexpectedTexts: []string{},
		},
		{
			name:            "Русский язык через флаг",
			args:            []string{"--lang", "ru", "--help"},
			expectedTexts:   []string{"Использование:", "Флаги:", "генерации"},
			unexpectedTexts: []string{},
		},
		{
			name:            "Русский язык - краткий флаг",
			args:            []string{"-l", "ru", "--help"},
			expectedTexts:   []string{"Использование:", "Флаги:"},
			unexpectedTexts: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Logf("Команда завершилась с ошибкой (ожидаемо для --help): %v", err)
			}

			output := stdout.String() + stderr.String()

			// Проверяем ожидаемые тексты
			for _, expected := range tt.expectedTexts {
				if !strings.Contains(output, expected) {
					t.Errorf("Ожидаемый текст '%s' не найден в выводе:\n%s", expected, output)
				}
			}

			// Проверяем отсутствие неожиданных текстов
			for _, unexpected := range tt.unexpectedTexts {
				if strings.Contains(output, unexpected) {
					t.Errorf("Неожиданный текст '%s' найден в выводе:\n%s", unexpected, output)
				}
			}
		})
	}
}

// TestEnvironmentVariables тестирует работу с переменными окружения
func TestEnvironmentVariables(t *testing.T) {
	binaryName := "pgen_test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Не удалось собрать приложение: %v", err)
	}

	tests := []struct {
		name        string
		envVars     map[string]string
		args        []string
		expectedTexts []string
	}{
		{
			name: "Русский язык через LANG",
			envVars: map[string]string{
				"LANG": "ru_RU.UTF-8",
			},
			args:          []string{"--help"},
			expectedTexts: []string{"Использование:", "Флаги:"},
		},
		{
			name: "Русский язык через LC_ALL",
			envVars: map[string]string{
				"LC_ALL": "ru_RU.UTF-8",
			},
			args:          []string{"--help"},
			expectedTexts: []string{"Использование:", "Флаги:"},
		},
		{
			name: "Английский язык принудительно",
			envVars: map[string]string{
				"LANG": "C",
			},
			args:          []string{"--lang", "en", "--help"},
			expectedTexts: []string{"Usage:", "Flags:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			
			// Устанавливаем переменные окружения
			env := os.Environ()
			for key, value := range tt.envVars {
				env = append(env, key+"="+value)
			}
			cmd.Env = env
			
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Logf("Команда завершилась с ошибкой (ожидаемо для --help): %v", err)
			}

			output := stdout.String() + stderr.String()

			// Проверяем ожидаемые тексты
			for _, expected := range tt.expectedTexts {
				if !strings.Contains(output, expected) {
					t.Errorf("Ожидаемый текст '%s' не найден в выводе:\n%s", expected, output)
				}
			}
		})
	}
}

// BenchmarkCLIStartup измеряет время запуска CLI
func BenchmarkCLIStartup(b *testing.B) {
	binaryName := "pgen_test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	tempDir := b.TempDir()
	binaryPath := filepath.Join(tempDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		b.Fatalf("Не удалось собрать приложение: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command(binaryPath, "--version")
		err := cmd.Run()
		if err != nil {
			b.Errorf("Ошибка выполнения команды: %v", err)
		}
	}
}

// BenchmarkCLIHelp измеряет время показа справки
func BenchmarkCLIHelp(b *testing.B) {
	binaryName := "pgen_test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	tempDir := b.TempDir()
	binaryPath := filepath.Join(tempDir, binaryName)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		b.Fatalf("Не удалось собрать приложение: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command(binaryPath, "--help")
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		_ = cmd.Run() // Игнорируем ошибку, так как --help возвращает код 0 или 1
	}
}
