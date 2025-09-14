package colors

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestColorVariables(t *testing.T) {
	// Тест на существование всех цветовых переменных
	tests := []struct {
		name     string
		colorVar *color.Color
	}{
		{"Success", Success},
		{"Error", Error},
		{"Info", Info},
		{"Prompt", Prompt},
		{"Generated", Generated},
		{"Title", Title},
		{"Subtle", Subtle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.colorVar == nil {
				t.Errorf("Цветовая переменная %s не должна быть nil", tt.name)
			}
		})
	}
}

func TestSuccessMsg(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Простое сообщение",
			input:    "Операция выполнена успешно",
			expected: "Операция выполнена успешно",
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "Сообщение с символами",
			input:    "Пароль сгенерирован: MyPass123!",
			expected: "Пароль сгенерирован: MyPass123!",
		},
		{
			name:     "Unicode сообщение",
			input:    "Успех! ✅ Готово",
			expected: "Успех! ✅ Готово",
		},
		{
			name:     "Многострочное сообщение",
			input:    "Строка 1\nСтрока 2",
			expected: "Строка 1\nСтрока 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SuccessMsg(tt.input)
			
			// Проверяем, что результат содержит исходный текст
			if !strings.Contains(result, tt.expected) {
				t.Errorf("SuccessMsg() = %q, должно содержать %q", result, tt.expected)
			}

			// Проверяем, что результат содержит ANSI коды для цвета (только для непустых строк)
			if tt.input != "" && !containsANSICodes(result) {
				t.Errorf("SuccessMsg() должно содержать ANSI коды цвета, получено: %q", result)
			}
		})
	}
}

func TestErrorMsg(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Сообщение об ошибке",
			input:    "Произошла ошибка",
			expected: "Произошла ошибка",
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "Техническая ошибка",
			input:    "Error: file not found",
			expected: "Error: file not found",
		},
		{
			name:     "Ошибка с эмодзи",
			input:    "Ошибка! ❌ Не удалось",
			expected: "Ошибка! ❌ Не удалось",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ErrorMsg(tt.input)
			
			if !strings.Contains(result, tt.expected) {
				t.Errorf("ErrorMsg() = %q, должно содержать %q", result, tt.expected)
			}

			if tt.input != "" && !containsANSICodes(result) {
				t.Errorf("ErrorMsg() должно содержать ANSI коды цвета, получено: %q", result)
			}
		})
	}
}

func TestInfoMsg(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Информационное сообщение",
			input:    "Информация о системе",
			expected: "Информация о системе",
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "Подсказка",
			input:    "Используйте --help для справки",
			expected: "Используйте --help для справки",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InfoMsg(tt.input)
			
			if !strings.Contains(result, tt.expected) {
				t.Errorf("InfoMsg() = %q, должно содержать %q", result, tt.expected)
			}

			if tt.input != "" && !containsANSICodes(result) {
				t.Errorf("InfoMsg() должно содержать ANSI коды цвета, получено: %q", result)
			}
		})
	}
}

func TestPromptMsg(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Запрос ввода",
			input:    "Введите мастер-пароль:",
			expected: "Введите мастер-пароль:",
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "Вопрос пользователю",
			input:    "Продолжить? (y/n)",
			expected: "Продолжить? (y/n)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PromptMsg(tt.input)
			
			if !strings.Contains(result, tt.expected) {
				t.Errorf("PromptMsg() = %q, должно содержать %q", result, tt.expected)
			}

			if tt.input != "" && !containsANSICodes(result) {
				t.Errorf("PromptMsg() должно содержать ANSI коды цвета, получено: %q", result)
			}
		})
	}
}

func TestGeneratedMsg(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Сгенерированный пароль",
			input:    "MySecurePass123!",
			expected: "MySecurePass123!",
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "Длинный пароль",
			input:    "VeryLongGeneratedPasswordWith123!@#$%",
			expected: "VeryLongGeneratedPasswordWith123!@#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeneratedMsg(tt.input)
			
			if !strings.Contains(result, tt.expected) {
				t.Errorf("GeneratedMsg() = %q, должно содержать %q", result, tt.expected)
			}

			if tt.input != "" && !containsANSICodes(result) {
				t.Errorf("GeneratedMsg() должно содержать ANSI коды цвета, получено: %q", result)
			}
		})
	}
}

func TestTitleMsg(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Заголовок приложения",
			input:    "PGen - Генератор паролей",
			expected: "PGen - Генератор паролей",
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "Заголовок раздела",
			input:    "Настройки безопасности",
			expected: "Настройки безопасности",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TitleMsg(tt.input)
			
			if !strings.Contains(result, tt.expected) {
				t.Errorf("TitleMsg() = %q, должно содержать %q", result, tt.expected)
			}

			if tt.input != "" && !containsANSICodes(result) {
				t.Errorf("TitleMsg() должно содержать ANSI коды цвета, получено: %q", result)
			}
		})
	}
}

func TestSubtleMsg(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Тонкое сообщение",
			input:    "Дополнительная информация",
			expected: "Дополнительная информация",
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "Подпись",
			input:    "Версия 1.0.0",
			expected: "Версия 1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SubtleMsg(tt.input)
			
			if !strings.Contains(result, tt.expected) {
				t.Errorf("SubtleMsg() = %q, должно содержать %q", result, tt.expected)
			}

			if tt.input != "" && !containsANSICodes(result) {
				t.Errorf("SubtleMsg() должно содержать ANSI коды цвета, получено: %q", result)
			}
		})
	}
}

func TestColorFunctionsWithSpecialCharacters(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	// Тест всех функций с особыми символами
	specialChars := []string{
		"!@#$%^&*()",
		"Русский текст с символами: №;%:?*",
		"Unicode: 🔐🛡️⚡🚀",
		"Tabs\tand\nnewlines",
		"\"Quotes\" and 'apostrophes'",
		"<HTML> & XML tags",
		"JSON: {\"key\": \"value\"}",
	}

	functions := map[string]func(string) string{
		"SuccessMsg":   SuccessMsg,
		"ErrorMsg":     ErrorMsg,
		"InfoMsg":      InfoMsg,
		"PromptMsg":    PromptMsg,
		"GeneratedMsg": GeneratedMsg,
		"TitleMsg":     TitleMsg,
		"SubtleMsg":    SubtleMsg,
	}

	for funcName, fn := range functions {
		for i, input := range specialChars {
			t.Run(fmt.Sprintf("%s_SpecialChar_%d", funcName, i+1), func(t *testing.T) {
				result := fn(input)
				
				if !strings.Contains(result, input) {
					t.Errorf("%s() = %q, должно содержать %q", funcName, result, input)
				}

				if !containsANSICodes(result) {
					t.Errorf("%s() должно содержать ANSI коды цвета", funcName)
				}
			})
		}
	}
}

func TestColorDisabled(t *testing.T) {
	// Тест поведения при отключенных цветах
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()

	// Отключаем цвета
	color.NoColor = true

	testInput := "Тестовое сообщение"
	
	functions := map[string]func(string) string{
		"SuccessMsg":   SuccessMsg,
		"ErrorMsg":     ErrorMsg,
		"InfoMsg":      InfoMsg,
		"PromptMsg":    PromptMsg,
		"GeneratedMsg": GeneratedMsg,
		"TitleMsg":     TitleMsg,
		"SubtleMsg":    SubtleMsg,
	}

	for funcName, fn := range functions {
		t.Run(fmt.Sprintf("%s_NoColor", funcName), func(t *testing.T) {
			result := fn(testInput)
			
			// При отключенных цветах результат должен быть равен входной строке
			if result != testInput {
				t.Errorf("%s() с отключенными цветами = %q, ожидается %q", funcName, result, testInput)
			}

			// Не должно быть ANSI кодов
			if containsANSICodes(result) {
				t.Errorf("%s() с отключенными цветами не должно содержать ANSI коды: %q", funcName, result)
			}
		})
	}
}

func TestAllFunctionsConsistency(t *testing.T) {
	// Принудительно включаем цвета для тестирования
	originalValue := color.NoColor
	defer func() {
		color.NoColor = originalValue
	}()
	color.NoColor = false

	// Тест на консистентность всех функций
	testInput := "Консистентность"
	
	functions := []func(string) string{
		SuccessMsg,
		ErrorMsg,
		InfoMsg,
		PromptMsg,
		GeneratedMsg,
		TitleMsg,
		SubtleMsg,
	}

	for i, fn := range functions {
		t.Run(fmt.Sprintf("Function_%d_Consistency", i+1), func(t *testing.T) {
			result := fn(testInput)
			
			// Все функции должны содержать исходный текст
			if !strings.Contains(result, testInput) {
				t.Errorf("Функция %d не содержит исходный текст", i+1)
			}

			// Все функции должны возвращать непустую строку для непустого ввода
			if len(result) == 0 {
				t.Errorf("Функция %d возвращает пустую строку для непустого ввода", i+1)
			}

			// Результат должен быть длиннее исходной строки (из-за ANSI кодов)
			if len(result) <= len(testInput) {
				t.Errorf("Функция %d: результат должен быть длиннее исходной строки", i+1)
			}
		})
	}
}

// Вспомогательная функция для проверки наличия ANSI кодов
func containsANSICodes(s string) bool {
	// ANSI коды начинаются с ESC[ (или \x1b[)
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.MatchString(s)
}

// Бенчмарки для измерения производительности
func BenchmarkSuccessMsg(b *testing.B) {
	msg := "Операция выполнена успешно"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SuccessMsg(msg)
	}
}

func BenchmarkErrorMsg(b *testing.B) {
	msg := "Произошла ошибка"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrorMsg(msg)
	}
}

func BenchmarkAllColorFunctions(b *testing.B) {
	msg := "Тестовое сообщение для бенчмарка"
	
	functions := []func(string) string{
		SuccessMsg,
		ErrorMsg,
		InfoMsg,
		PromptMsg,
		GeneratedMsg,
		TitleMsg,
		SubtleMsg,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, fn := range functions {
			_ = fn(msg)
		}
	}
}

func BenchmarkLongMessage(b *testing.B) {
	longMsg := strings.Repeat("Очень длинное сообщение с множеством символов ", 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SuccessMsg(longMsg)
	}
}
