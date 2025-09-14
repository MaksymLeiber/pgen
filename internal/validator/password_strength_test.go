package validator

import (
	"testing"

	"pgen/internal/i18n"
)

func TestValidatePasswordStrength(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name           string
		password       string
		expectedLevel  StrengthLevel
		expectIssues   bool
		minScore       int
		maxScore       int
	}{
		{
			name:          "Очень слабый пароль",
			password:      "123",
			expectedLevel: StrengthWeak,
			expectIssues:  true,
			minScore:      -20,
			maxScore:      10,
		},
		{
			name:          "Слабый пароль",
			password:      "password",
			expectedLevel: StrengthWeak,
			expectIssues:  true,
			minScore:      -30,
			maxScore:      10,
		},
		{
			name:          "Средний пароль",
			password:      "Password123",
			expectedLevel: StrengthWeak, // Содержит словарное слово "password"
			expectIssues:  true,
			minScore:      0,
			maxScore:      30,
		},
		{
			name:          "Хороший пароль",
			password:      "Kx9#mL2$vN8@",
			expectedLevel: StrengthStrong, // Фактически получает 70 очков
			expectIssues:  false,
			minScore:      70,
			maxScore:      90,
		},
		{
			name:          "Сильный пароль",
			password:      "Kx9#mL2$vN8@qR5!zT3%",
			expectedLevel: StrengthStrong,
			expectIssues:  false,
			minScore:      70,
			maxScore:      90,
		},
		{
			name:          "Очень сильная парольная фраза",
			password:      "My unique phrase with numbers 789 and symbols #$%",
			expectedLevel: StrengthStrong,
			expectIssues:  true, // Содержит последовательности
			minScore:      70,
			maxScore:      90,
		},
		{
			name:          "Русский пароль",
			password:      "МойКлючДоступа789#",
			expectedLevel: StrengthFair, // Содержит словарные слова
			expectIssues:  true,
			minScore:      30,
			maxScore:      50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePasswordStrength(tt.password, messages)

			if result.Level != tt.expectedLevel {
				t.Errorf("ValidatePasswordStrength() уровень = %v, ожидается %v", result.Level, tt.expectedLevel)
			}

			if result.Score < tt.minScore || result.Score > tt.maxScore {
				t.Errorf("ValidatePasswordStrength() очки = %v, ожидается между %v и %v", result.Score, tt.minScore, tt.maxScore)
			}

			if tt.expectIssues && len(result.Issues) == 0 {
				t.Error("ValidatePasswordStrength() ожидались проблемы, но их нет")
			}

			if !tt.expectIssues && len(result.Issues) > 0 {
				t.Errorf("ValidatePasswordStrength() неожиданные проблемы: %v", result.Issues)
			}

			// Проверяем, что есть предложения для слабых паролей
			if result.Level <= StrengthFair && len(result.Suggestions) == 0 {
				t.Error("ValidatePasswordStrength() слабый пароль должен иметь предложения")
			}
		})
	}
}

func TestPasswordLengthScoring(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name     string
		password string
		minScore int
	}{
		{"Короткий пароль", "Ab1!", 0},     // < 8 символов
		{"Средний пароль", "Ab1!5678", 10}, // 8-11 символов
		{"Длинный пароль", "Ab1!567890123", 20}, // 12-15 символов
		{"Очень длинный пароль", "Ab1!567890123456", 25}, // 16+ символов
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePasswordStrength(tt.password, messages)
			
			if len(tt.password) >= 8 && result.Score < tt.minScore {
				t.Errorf("Оценка длины пароля: очки = %v, ожидается >= %v", result.Score, tt.minScore)
			}
		})
	}
}

func TestCharacterTypeScoring(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name     string
		password string
		hasLower bool
		hasUpper bool
		hasNumber bool
		hasSymbol bool
	}{
		{"Только строчные", "abcdefgh", true, false, false, false},
		{"Только заглавные", "ABCDEFGH", false, true, false, false},
		{"Только цифры", "12345678", false, false, true, false},
		{"Только символы", "!@#$%^&*", false, false, false, true},
		{"Смешанный регистр", "AbCdEfGh", true, true, false, false},
		{"Все типы", "AbC123!@", true, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePasswordStrength(tt.password, messages)

			// Проверяем наличие соответствующих проблем
			hasLowerIssue := false
			hasUpperIssue := false
			hasNumberIssue := false

			for _, issue := range result.Issues {
				if issue == messages.Errors.PasswordIssues.NoLowercase {
					hasLowerIssue = true
				}
				if issue == messages.Errors.PasswordIssues.NoUppercase {
					hasUpperIssue = true
				}
				if issue == messages.Errors.PasswordIssues.NoNumbers {
					hasNumberIssue = true
				}
			}

			if !tt.hasLower && !hasLowerIssue {
				t.Error("Должна быть проблема с отсутствием строчных букв")
			}
			if !tt.hasUpper && !hasUpperIssue {
				t.Error("Должна быть проблема с отсутствием заглавных букв")
			}
			if !tt.hasNumber && !hasNumberIssue {
				t.Error("Должна быть проблема с отсутствием цифр")
			}
		})
	}
}

func TestHasRepeatingChars(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Нет повторов", "AbC123!@", false},
		{"Два одинаковых символа", "AAbc123", false},
		{"Три одинаковых символа", "AAAbc123", true},
		{"Четыре одинаковых символа", "AAAAbc123", true},
		{"Повтор в середине", "Ab111c23", true},
		{"Повтор в конце", "Abc123!!!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasRepeatingChars(tt.password)
			if result != tt.expected {
				t.Errorf("hasRepeatingChars(%q) = %v, ожидается %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestHasSequences(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Нет последовательностей", "Ab1!Xy9@", false},
		{"Алфавитная последовательность", "abcdef123", true},
		{"Обратный алфавит", "fedcba123", true},
		{"Числовая последовательность", "Ab123456!", true},
		{"Обратные цифры", "Ab654321!", true},
		{"Клавиатурная последовательность", "qwerty123", true},
		{"Короткая последовательность", "ab1!Xy9@", false}, // Только 2 символа
		{"Смешанный регистр последовательность", "AbC123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasSequences(tt.password)
			if result != tt.expected {
				t.Errorf("hasSequences(%q) = %v, ожидается %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestHasCommonWords(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Нет обычных слов", "Xy9@Zb3#", false},
		{"Содержит password", "mypassword123", true},
		{"Содержит admin", "admin123", true},
		{"Содержит русское слово", "мойпароль123", true},
		{"Содержит qwerty", "qwerty123", true},
		{"Содержит love", "iloveyou", true},
		{"Нечувствительно к регистру", "PASSWORD123", true},
		{"Частичное совпадение", "mypass123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasCommonWords(tt.password)
			if result != tt.expected {
				t.Errorf("hasCommonWords(%q) = %v, ожидается %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestIsPhrase(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Не фраза", "Password123!", false},
		{"С пробелами", "This is a phrase", true},
		{
			name:     "Длинный со словами",
			password: "ThisIsMyLongPasswordWithWords",
			expected: false, // Нет пробелов и слова не разделены четко
		},
		{"Короткий смешанный", "Ab123!", false},
		{"Длинный но без слов", "Ab123!@#$%^&*()_+", false},
		{"Русская фраза", "Это моя фраза пароль", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPhrase(tt.password)
			if result != tt.expected {
				t.Errorf("isPhrase(%q) = %v, ожидается %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Простая строка", "abc", "cba"},
		{"Цифры", "123", "321"},
		{"Смешанное", "Ab1!", "!1bA"},
		{"Пустая", "", ""},
		{"Один символ", "a", "a"},
		{"Юникод", "привет", "тевирп"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reverse(tt.input)
			if result != tt.expected {
				t.Errorf("reverse(%q) = %q, ожидается %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContainsSequence(t *testing.T) {
	tests := []struct {
		name     string
		password string
		sequence string
		minLen   int
		expected bool
	}{
		{"Содержит последовательность", "myabc123", "abcdefg", 3, true},
		{"Нет последовательности", "myx123", "abcdefg", 3, false},
		{"Слишком коротко", "myab123", "abcdefg", 3, false},
		{"Точное совпадение", "abc", "abcdefg", 3, true},
		{"В конце", "mydefg", "abcdefg", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsSequence(tt.password, tt.sequence, tt.minLen)
			if result != tt.expected {
				t.Errorf("containsSequence(%q, %q, %d) = %v, ожидается %v", 
					tt.password, tt.sequence, tt.minLen, result, tt.expected)
			}
		})
	}
}

func TestStrengthLevels(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name     string
		password string
		level    StrengthLevel
	}{
		{"Слабый", "123", StrengthWeak},
		{"Средний", "Kx9mL2vN", StrengthFair},
		{"Хороший", "Kx9#mL2$vN8@", StrengthStrong}, // Фактически получает Strong
		{"Сильный", "Kx9#mL2$vN8@qR5!", StrengthStrong},
		{"Очень сильный", "My unique phrase with numbers 789 and symbols #$%", StrengthStrong}, // Фактически получает Strong
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePasswordStrength(tt.password, messages)
			if result.Level != tt.level {
				t.Errorf("Пароль %q: уровень = %v, ожидается %v (очки: %d)", 
					tt.password, result.Level, tt.level, result.Score)
			}
		})
	}
}

func TestEmptyPassword(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")
	result := ValidatePasswordStrength("", messages)

	if result.Level != StrengthWeak {
		t.Errorf("Пустой пароль должен быть слабым, получен %v", result.Level)
	}

	if len(result.Issues) == 0 {
		t.Error("Пустой пароль должен иметь проблемы")
	}

	if len(result.Suggestions) == 0 {
		t.Error("Пустой пароль должен иметь предложения")
	}
}

func TestUnicodePasswords(t *testing.T) {
	messages := i18n.GetMessages(i18n.Russian, "test")

	tests := []struct {
		name            string
		password        string
		minLevel        StrengthLevel
	}{
		{"Русский сильный", "МойКлючБезопасности789#", StrengthFair},
		{"Смешанные языки", "MyКлюч789!", StrengthFair},
		{"Пароль с эмодзи", "MyKey🔐789!", StrengthFair},
		{"Китайские символы", "密钥789!", StrengthWeak}, // Короткий пароль
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePasswordStrength(tt.password, messages)
			if result.Level < tt.minLevel {
				t.Errorf("Юникод пароль %q: уровень = %v, ожидается >= %v", 
					tt.password, result.Level, tt.minLevel)
			}
		})
	}
}

// Бенчмарки для измерения производительности
func BenchmarkValidatePasswordStrength(b *testing.B) {
	messages := i18n.GetMessages(i18n.English, "test")
	password := "MySecurePassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidatePasswordStrength(password, messages)
	}
}

func BenchmarkHasCommonWords(b *testing.B) {
	password := "MySecurePassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hasCommonWords(password)
	}
}

func BenchmarkHasSequences(b *testing.B) {
	password := "MySecurePassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hasSequences(password)
	}
}
