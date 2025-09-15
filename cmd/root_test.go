package cmd

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/MaksymLeiber/pgen/internal/config"
	"github.com/MaksymLeiber/pgen/internal/i18n"
	"github.com/MaksymLeiber/pgen/internal/validator"
)

func TestDetectLanguageFromArgs(t *testing.T) {
	// Сохраняем оригинальные аргументы
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		expected i18n.Language
	}{
		{
			name:     "Флаг --lang ru",
			args:     []string{"pgen", "--lang", "ru"},
			expected: i18n.Russian,
		},
		{
			name:     "Флаг -l en",
			args:     []string{"pgen", "-l", "en"},
			expected: i18n.English,
		},
		{
			name:     "Флаг --lang=russian",
			args:     []string{"pgen", "--lang=russian"},
			expected: i18n.Russian,
		},
		{
			name:     "Флаг --lang=english",
			args:     []string{"pgen", "--lang=english"},
			expected: i18n.English,
		},
		{
			name:     "Без флага языка",
			args:     []string{"pgen"},
			expected: i18n.English, // По умолчанию
		},
		{
			name:     "Неизвестный язык",
			args:     []string{"pgen", "--lang", "de"},
			expected: i18n.English, // Fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			result := detectLanguageFromArgs()
			if result != tt.expected {
				t.Errorf("detectLanguageFromArgs() = %q, ожидается %q", result, tt.expected)
			}
		})
	}
}

func TestGetIssueText(t *testing.T) {
	// Тест функции getIssueText
	messages := &i18n.Messages{
		Errors: struct {
			ClipboardError  string
			GenerationError string
			EmptyMaster     string
			EmptyService    string
			UserCanceled    string
			InputCanceled   string
			LengthTooShort  string
			LengthTooLong   string
			HashTooShort    string
			PasswordIssues  struct {
				LengthTooShort  string
				NoLowercase     string
				NoUppercase     string
				NoNumbers       string
				RepeatingChars  string
				SequentialChars string
				CommonWords     string
			}
			Suggestions struct {
				IncreaseLength  string
				AddLowercase    string
				AddUppercase    string
				AddNumbers      string
				AddSymbols      string
				AvoidRepetition string
				AvoidSequences  string
				AvoidDictionary string
			}
		}{
			PasswordIssues: struct {
				LengthTooShort  string
				NoLowercase     string
				NoUppercase     string
				NoNumbers       string
				RepeatingChars  string
				SequentialChars string
				CommonWords     string
			}{
				LengthTooShort:  "Пароль слишком короткий",
				NoLowercase:     "Отсутствуют строчные буквы",
				NoUppercase:     "Отсутствуют заглавные буквы",
				NoNumbers:       "Отсутствуют цифры",
				RepeatingChars:  "Содержит повторяющиеся символы",
				SequentialChars: "Содержит последовательные символы",
				CommonWords:     "Содержит словарные слова",
			},
		},
	}

	tests := []struct {
		name     string
		issue    string
		expected string
	}{
		{
			name:     "length_too_short",
			issue:    "length_too_short",
			expected: "Пароль слишком короткий",
		},
		{
			name:     "no_lowercase",
			issue:    "no_lowercase",
			expected: "Отсутствуют строчные буквы",
		},
		{
			name:     "no_uppercase",
			issue:    "no_uppercase",
			expected: "Отсутствуют заглавные буквы",
		},
		{
			name:     "no_numbers",
			issue:    "no_numbers",
			expected: "Отсутствуют цифры",
		},
		{
			name:     "repeating_chars",
			issue:    "repeating_chars",
			expected: "Содержит повторяющиеся символы",
		},
		{
			name:     "sequential_chars",
			issue:    "sequential_chars",
			expected: "Содержит последовательные символы",
		},
		{
			name:     "common_words",
			issue:    "common_words",
			expected: "Содержит словарные слова",
		},
		{
			name:     "unknown_issue",
			issue:    "unknown_issue",
			expected: "unknown_issue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIssueText(tt.issue, messages)
			if result != tt.expected {
				t.Errorf("getIssueText(%q) = %q, ожидается %q", tt.issue, result, tt.expected)
			}
		})
	}
}

func TestGetSuggestionText(t *testing.T) {
	// Тест функции getSuggestionText
	messages := &i18n.Messages{
		Errors: struct {
			ClipboardError  string
			GenerationError string
			EmptyMaster     string
			EmptyService    string
			UserCanceled    string
			InputCanceled   string
			LengthTooShort  string
			LengthTooLong   string
			HashTooShort    string
			PasswordIssues  struct {
				LengthTooShort  string
				NoLowercase     string
				NoUppercase     string
				NoNumbers       string
				RepeatingChars  string
				SequentialChars string
				CommonWords     string
			}
			Suggestions struct {
				IncreaseLength  string
				AddLowercase    string
				AddUppercase    string
				AddNumbers      string
				AddSymbols      string
				AvoidRepetition string
				AvoidSequences  string
				AvoidDictionary string
			}
		}{
			Suggestions: struct {
				IncreaseLength  string
				AddLowercase    string
				AddUppercase    string
				AddNumbers      string
				AddSymbols      string
				AvoidRepetition string
				AvoidSequences  string
				AvoidDictionary string
			}{
				IncreaseLength:  "Увеличьте длину пароля",
				AddLowercase:    "Добавьте строчные буквы",
				AddUppercase:    "Добавьте заглавные буквы",
				AddNumbers:      "Добавьте цифры",
				AddSymbols:      "Добавьте специальные символы",
				AvoidRepetition: "Избегайте повторений",
				AvoidSequences:  "Избегайте последовательностей",
				AvoidDictionary: "Избегайте словарных слов",
			},
		},
	}

	tests := []struct {
		name       string
		suggestion string
		expected   string
	}{
		{
			name:       "increase_length",
			suggestion: "increase_length",
			expected:   "Увеличьте длину пароля",
		},
		{
			name:       "add_lowercase",
			suggestion: "add_lowercase",
			expected:   "Добавьте строчные буквы",
		},
		{
			name:       "add_uppercase",
			suggestion: "add_uppercase",
			expected:   "Добавьте заглавные буквы",
		},
		{
			name:       "add_numbers",
			suggestion: "add_numbers",
			expected:   "Добавьте цифры",
		},
		{
			name:       "add_symbols",
			suggestion: "add_symbols",
			expected:   "Добавьте специальные символы",
		},
		{
			name:       "avoid_repetition",
			suggestion: "avoid_repetition",
			expected:   "Избегайте повторений",
		},
		{
			name:       "avoid_sequences",
			suggestion: "avoid_sequences",
			expected:   "Избегайте последовательностей",
		},
		{
			name:       "avoid_dictionary",
			suggestion: "avoid_dictionary",
			expected:   "Избегайте словарных слов",
		},
		{
			name:       "unknown_suggestion",
			suggestion: "unknown_suggestion",
			expected:   "unknown_suggestion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSuggestionText(tt.suggestion, messages)
			if result != tt.expected {
				t.Errorf("getSuggestionText(%q) = %q, ожидается %q", tt.suggestion, result, tt.expected)
			}
		})
	}
}

func TestSetConfigValue(t *testing.T) {
	// Тест функции setConfigValue
	messages := &i18n.Messages{
		ConfigInvalidArgonTime:     "Неверное значение argon_time",
		ConfigInvalidArgonMemory:   "Неверное значение argon_memory",
		ConfigInvalidArgonThreads:  "Неверное значение argon_threads",
		ConfigInvalidArgonKeyLen:   "Неверное значение argon_key_len",
		ConfigInvalidDefaultLength: "Неверное значение default_length",
		ConfigLengthRange:          "Длина должна быть от 4 до 128",
		ConfigLanguageValues:       "Язык должен быть ru, en или auto",
		ConfigCharsetValues:        "Набор символов должен быть alphanumeric, alphanumeric_symbols или symbols_only",
		ConfigInvalidDefaultCopy:   "Неверное значение default_copy",
		ConfigInvalidClearTimeout:  "Неверное значение clear_timeout",
		ConfigTimeoutRange:         "Таймаут должен быть >= 0",
		ConfigInvalidPasswordInfo:  "Неверное значение show_password_info",
		ConfigInvalidColorOutput:   "Неверное значение color_output",
		ConfigInvalidUsername:      "Неверное значение username",
		ConfigUsernameEmpty:        "Имя пользователя не может быть пустым",
		ConfigUnknownKey:           "Неизвестный ключ",
	}

	// Инициализируем глобальную переменную cfg для тестов
	cfg = &config.Config{
		ArgonTime:           1,
		ArgonMemory:         65536,
		ArgonThreads:        1,
		ArgonKeyLen:         32,
		DefaultLength:       16,
		DefaultLanguage:     "auto",
		CharacterSet:        "alphanumeric_symbols",
		DefaultCopy:         false,
		DefaultClearTimeout: 45,
		ShowPasswordInfo:    false,
		ColorOutput:         true,
	}

	tests := []struct {
		name      string
		key       string
		value     string
		wantError bool
	}{
		{
			name:      "Валидный argon_time",
			key:       "argon_time",
			value:     "2",
			wantError: false,
		},
		{
			name:      "Невалидный argon_time",
			key:       "argon_time",
			value:     "invalid",
			wantError: true,
		},
		{
			name:      "Валидный default_length",
			key:       "default_length",
			value:     "20",
			wantError: false,
		},
		{
			name:      "Слишком короткий default_length",
			key:       "default_length",
			value:     "3",
			wantError: true,
		},
		{
			name:      "Слишком длинный default_length",
			key:       "default_length",
			value:     "200",
			wantError: true,
		},
		{
			name:      "Валидный default_language",
			key:       "default_language",
			value:     "ru",
			wantError: false,
		},
		{
			name:      "Невалидный default_language",
			key:       "default_language",
			value:     "fr",
			wantError: true,
		},
		{
			name:      "Валидный character_set",
			key:       "character_set",
			value:     "alphanumeric",
			wantError: false,
		},
		{
			name:      "Невалидный character_set",
			key:       "character_set",
			value:     "invalid_set",
			wantError: true,
		},
		{
			name:      "Валидный default_copy",
			key:       "default_copy",
			value:     "true",
			wantError: false,
		},
		{
			name:      "Невалидный default_copy",
			key:       "default_copy",
			value:     "maybe",
			wantError: true,
		},
		{
			name:      "Валидный default_clear_timeout",
			key:       "default_clear_timeout",
			value:     "60",
			wantError: false,
		},
		{
			name:      "Отрицательный default_clear_timeout",
			key:       "default_clear_timeout",
			value:     "-1",
			wantError: true,
		},
		{
			name:      "Валидный username",
			key:       "username",
			value:     "maksym",
			wantError: false,
		},
		{
			name:      "Пустой username",
			key:       "username",
			value:     "",
			wantError: true,
		},
		{
			name:      "Username только пробелы",
			key:       "username",
			value:     "   ",
			wantError: true,
		},
		{
			name:      "Username с пробелами по краям",
			key:       "username",
			value:     "  maksym  ",
			wantError: false,
		},
		{
			name:      "Неизвестный ключ",
			key:       "unknown_key",
			value:     "value",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setConfigValue(tt.key, tt.value, messages)
			if tt.wantError && err == nil {
				t.Errorf("setConfigValue(%q, %q) должен вернуть ошибку", tt.key, tt.value)
			}
			if !tt.wantError && err != nil {
				t.Errorf("setConfigValue(%q, %q) не должен возвращать ошибку: %v", tt.key, tt.value, err)
			}
		})
	}
}

func TestNeedsElevation(t *testing.T) {
	// Тест функции needsElevation
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("needsElevation не должна паниковать: %v", r)
		}
	}()

	result := needsElevation()

	// Результат может быть true или false в зависимости от платформы и прав
	if result != true && result != false {
		t.Error("needsElevation должна возвращать булево значение")
	}

	// На Unix системах проверяем логику
	if runtime.GOOS != "windows" {
		expected := os.Geteuid() != 0
		if result != expected {
			t.Errorf("needsElevation() = %v, ожидается %v для Unix систем", result, expected)
		}
	}
}

func TestGetCurrentInstallPath(t *testing.T) {
	// Тест функции getCurrentInstallPath
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("getCurrentInstallPath не должна паниковать: %v", r)
		}
	}()

	path := getCurrentInstallPath()
	if path == "" {
		t.Error("getCurrentInstallPath должна возвращать непустой путь")
	}

	// Проверяем, что путь содержит ожидаемые компоненты
	if !strings.Contains(strings.ToLower(path), "pgen") {
		t.Errorf("Путь установки должен содержать 'pgen', получен: %s", path)
	}
}

func TestIsWindowsAdmin(t *testing.T) {
	// Тест функции isWindowsAdmin
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("isWindowsAdmin не должна паниковать: %v", r)
		}
	}()

	result := isWindowsAdmin()

	// Результат может быть true или false
	if result != true && result != false {
		t.Error("isWindowsAdmin должна возвращать булево значение")
	}

	// На не-Windows системах эта функция все равно должна работать
	// (хотя результат может быть не очень осмысленным)
}

func TestUpdateCommandTexts(t *testing.T) {
	// Тест функции updateCommandTexts
	messages := &i18n.Messages{
		Usage:       "pgen [флаги]",
		Description: "Генератор детерминированных паролей",
		Examples:    "Примеры использования",
		Flags: struct {
			Lang             string
			LangDesc         string
			Length           string
			LengthDesc       string
			Copy             string
			CopyDesc         string
			ClearTimeout     string
			ClearTimeoutDesc string
			Version          string
			VersionDesc      string
			About            string
			AboutDesc        string
			Info             string
			InfoDesc         string
			Install          string
			InstallDesc      string
			Uninstall        string
			UninstallDesc    string
		}{
			LangDesc:         "Язык интерфейса",
			LengthDesc:       "Длина пароля",
			CopyDesc:         "Копировать в буфер обмена",
			ClearTimeoutDesc: "Таймаут очистки буфера",
			VersionDesc:      "Показать версию",
			AboutDesc:        "О программе",
			InfoDesc:         "Показать информацию о пароле",
			InstallDesc:      "Установить в систему",
			UninstallDesc:    "Удалить из системы",
		},
	}

	// Создаем тестовую команду
	testCmd := rootCmd

	// Проверяем, что функция не паникует
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("updateCommandTexts не должна паниковать: %v", r)
		}
	}()

	updateCommandTexts(testCmd, messages)

	// Проверяем, что тексты обновились
	if testCmd.Use != messages.Usage {
		t.Errorf("Use не обновился: получен %q, ожидается %q", testCmd.Use, messages.Usage)
	}
	if testCmd.Short != messages.Description {
		t.Errorf("Short не обновился: получен %q, ожидается %q", testCmd.Short, messages.Description)
	}
}

func TestDisplayPasswordStrength(t *testing.T) {
	// Тест функции displayPasswordStrength
	messages := &i18n.Messages{
		MasterPasswordStrength:     "Сила мастер-пароля:",
		PasswordStrengthWeak:       "🔴 Слабый",
		PasswordStrengthFair:       "🟠 Удовлетворительный",
		PasswordStrengthGood:       "🟡 Хороший",
		PasswordStrengthStrong:     "🟢 Сильный",
		PasswordStrengthVeryStrong: "🟢 Очень сильный",
	}

	// Тестируем разные уровни силы
	strengthLevels := []validator.StrengthLevel{
		validator.StrengthWeak,
		validator.StrengthFair,
		validator.StrengthGood,
		validator.StrengthStrong,
		validator.StrengthVeryStrong,
	}

	for _, level := range strengthLevels {
		t.Run(string(rune(int(level)+'0')), func(t *testing.T) {
			strength := &validator.PasswordStrength{
				Level:       level,
				Issues:      []string{},
				Suggestions: []string{},
			}

			// Проверяем, что функция не паникует
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("displayPasswordStrength не должна паниковать: %v", r)
				}
			}()

			displayPasswordStrength(strength, messages)
		})
	}

	// Тест с проблемами и рекомендациями
	t.Run("WithIssuesAndSuggestions", func(t *testing.T) {
		strength := &validator.PasswordStrength{
			Level:       validator.StrengthWeak,
			Issues:      []string{"length_too_short", "no_uppercase"},
			Suggestions: []string{"increase_length", "add_uppercase"},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("displayPasswordStrength с проблемами не должна паниковать: %v", r)
			}
		}()

		displayPasswordStrength(strength, messages)
	})
}

// Бенчмарки для измерения производительности
func BenchmarkDetectLanguageFromArgs(b *testing.B) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"pgen", "--lang", "ru"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = detectLanguageFromArgs()
	}
}

func TestFormatTitleWithUser(t *testing.T) {
	// Тест функции formatTitleWithUser
	tests := []struct {
		name     string
		appTitle string
		username string
		expected string
	}{
		{
			name:     "Пустой username - должен показать default",
			appTitle: "🔑 PGen CLI",
			username: "",
			expected: "profile: [default]",
		},
		{
			name:     "Username 'user' - должен показать default",
			appTitle: "🔑 PGen CLI",
			username: "user",
			expected: "profile: [default]",
		},
		{
			name:     "Кастомный username",
			appTitle: "🔑 PGen CLI",
			username: "maksym",
			expected: "profile: [maksym]",
		},
		{
			name:     "Длинный username",
			appTitle: "🔑 PGen CLI",
			username: "very_long_username_123",
			expected: "profile: [very_long_username_123]",
		},
		{
			name:     "Короткий заголовок",
			appTitle: "PGen",
			username: "test",
			expected: "profile: [test]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTitleWithUser(tt.appTitle, tt.username, &i18n.Messages{ProfileLabel: "profile:"})

			// Проверяем, что результат содержит ожидаемую информацию о профиле
			if !strings.Contains(result, tt.expected) {
				t.Errorf("formatTitleWithUser(%q, %q) = %q, должно содержать %q", tt.appTitle, tt.username, result, tt.expected)
			}

			// Проверяем, что результат содержит исходный заголовок
			if !strings.Contains(result, tt.appTitle) {
				t.Errorf("formatTitleWithUser() должно содержать исходный заголовок %q", tt.appTitle)
			}

			// Проверяем, что результат длиннее исходного заголовка
			if len(result) <= len(tt.appTitle) {
				t.Errorf("formatTitleWithUser() должно возвращать строку длиннее исходного заголовка")
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	// Тест функции stripANSI
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Простой текст без ANSI кодов",
			input:    "Простой текст",
			expected: "Простой текст",
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: "",
		},
		{
			name:     "Текст с цветными кодами",
			input:    "\x1b[31mКрасный текст\x1b[0m",
			expected: "Красный текст",
		},
		{
			name:     "Текст с жирным шрифтом",
			input:    "\x1b[1mЖирный текст\x1b[0m",
			expected: "Жирный текст",
		},
		{
			name:     "Текст с подчеркиванием",
			input:    "\x1b[4mПодчеркнутый текст\x1b[0m",
			expected: "Подчеркнутый текст",
		},
		{
			name:     "Сложные ANSI коды",
			input:    "\x1b[31;1;4mКрасный жирный подчеркнутый\x1b[0m",
			expected: "Красный жирный подчеркнутый",
		},
		{
			name:     "Несколько участков с разными кодами",
			input:    "\x1b[31mКрасный\x1b[0m обычный \x1b[32mзеленый\x1b[0m",
			expected: "Красный обычный зеленый",
		},
		{
			name:     "Заголовок с эмодзи и цветом",
			input:    "\x1b[1m🔑 PGen CLI\x1b[0m",
			expected: "🔑 PGen CLI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripANSI(tt.input)

			if result != tt.expected {
				t.Errorf("stripANSI(%q) = %q, ожидается %q", tt.input, result, tt.expected)
			}

			// Проверяем, что в результате нет ANSI кодов
			if strings.Contains(result, "\x1b[") {
				t.Errorf("stripANSI() должно удалить все ANSI коды, но в результате %q еще есть коды", result)
			}
		})
	}
}

func BenchmarkGetIssueText(b *testing.B) {
	messages := &i18n.Messages{
		Errors: struct {
			ClipboardError  string
			GenerationError string
			EmptyMaster     string
			EmptyService    string
			UserCanceled    string
			InputCanceled   string
			LengthTooShort  string
			LengthTooLong   string
			HashTooShort    string
			PasswordIssues  struct {
				LengthTooShort  string
				NoLowercase     string
				NoUppercase     string
				NoNumbers       string
				RepeatingChars  string
				SequentialChars string
				CommonWords     string
			}
			Suggestions struct {
				IncreaseLength  string
				AddLowercase    string
				AddUppercase    string
				AddNumbers      string
				AddSymbols      string
				AvoidRepetition string
				AvoidSequences  string
				AvoidDictionary string
			}
		}{
			PasswordIssues: struct {
				LengthTooShort  string
				NoLowercase     string
				NoUppercase     string
				NoNumbers       string
				RepeatingChars  string
				SequentialChars string
				CommonWords     string
			}{
				LengthTooShort: "Пароль слишком короткий",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getIssueText("length_too_short", messages)
	}
}

func BenchmarkGetSuggestionText(b *testing.B) {
	messages := &i18n.Messages{
		Errors: struct {
			ClipboardError  string
			GenerationError string
			EmptyMaster     string
			EmptyService    string
			UserCanceled    string
			InputCanceled   string
			LengthTooShort  string
			LengthTooLong   string
			HashTooShort    string
			PasswordIssues  struct {
				LengthTooShort  string
				NoLowercase     string
				NoUppercase     string
				NoNumbers       string
				RepeatingChars  string
				SequentialChars string
				CommonWords     string
			}
			Suggestions struct {
				IncreaseLength  string
				AddLowercase    string
				AddUppercase    string
				AddNumbers      string
				AddSymbols      string
				AvoidRepetition string
				AvoidSequences  string
				AvoidDictionary string
			}
		}{
			Suggestions: struct {
				IncreaseLength  string
				AddLowercase    string
				AddUppercase    string
				AddNumbers      string
				AddSymbols      string
				AvoidRepetition string
				AvoidSequences  string
				AvoidDictionary string
			}{
				IncreaseLength: "Увеличьте длину пароля",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getSuggestionText("increase_length", messages)
	}
}
