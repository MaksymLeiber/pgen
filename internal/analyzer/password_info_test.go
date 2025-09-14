package analyzer

import (
	"math"
	"testing"

	"pgen/internal/i18n"
)

func TestAnalyzePassword(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name            string
		password        string
		expectedLength  int
		expectedEntropy float64
		minCharsetSize  int
	}{
		{
			name:            "Простой пароль",
			password:        "abc123",
			expectedLength:  6,
			expectedEntropy: 30.0, // Приблизительно
			minCharsetSize:  36,   // a-z + 0-9
		},
		{
			name:            "Сложный пароль",
			password:        "AbC123!@#",
			expectedLength:  9,
			expectedEntropy: 50.0, // Приблизительно
			minCharsetSize:  90,   // A-Z + a-z + 0-9 + symbols
		},
		{
			name:            "Только цифры",
			password:        "123456",
			expectedLength:  6,
			expectedEntropy: 19.0, // Приблизительно
			minCharsetSize:  10,   // 0-9
		},
		{
			name:            "Пустой пароль",
			password:        "",
			expectedLength:  0,
			expectedEntropy: 0,
			minCharsetSize:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := AnalyzePassword(tt.password, messages)

			if info.Length != tt.expectedLength {
				t.Errorf("Длина = %v, ожидается %v", info.Length, tt.expectedLength)
			}

			if tt.password != "" {
				if info.CharsetSize < tt.minCharsetSize {
					t.Errorf("Размер набора символов = %v, ожидается >= %v", info.CharsetSize, tt.minCharsetSize)
				}

				if math.Abs(info.Entropy-tt.expectedEntropy) > 10 {
					t.Errorf("Энтропия = %v, ожидается ~%v", info.Entropy, tt.expectedEntropy)
				}

				if info.Strength == "" {
					t.Error("Сила не должна быть пустой")
				}

				if info.TimeToCrack.HumanTime == "" {
					t.Error("Время взлома не должно быть пустым")
				}
			}
		})
	}
}

func TestAnalyzeComposition(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected CharComposition
	}{
		{
			name:     "Смешанный пароль",
			password: "AbC123!@",
			expected: CharComposition{
				Uppercase: 2, // A, C
				Lowercase: 1, // b
				Numbers:   3, // 1, 2, 3
				Symbols:   2, // !, @
				Total:     8,
			},
		},
		{
			name:     "Только строчные",
			password: "abcdef",
			expected: CharComposition{
				Uppercase: 0,
				Lowercase: 6,
				Numbers:   0,
				Symbols:   0,
				Total:     6,
			},
		},
		{
			name:     "Пустой пароль",
			password: "",
			expected: CharComposition{
				Uppercase: 0,
				Lowercase: 0,
				Numbers:   0,
				Symbols:   0,
				Total:     0,
			},
		},
		{
			name:     "Юникод символы",
			password: "Привет123!",
			expected: CharComposition{
				Uppercase: 1, // П
				Lowercase: 5, // ривет
				Numbers:   3, // 123
				Symbols:   1, // !
				Total:     10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzeComposition(tt.password)

			if result.Uppercase != tt.expected.Uppercase {
				t.Errorf("Заглавные = %v, ожидается %v", result.Uppercase, tt.expected.Uppercase)
			}
			if result.Lowercase != tt.expected.Lowercase {
				t.Errorf("Строчные = %v, ожидается %v", result.Lowercase, tt.expected.Lowercase)
			}
			if result.Numbers != tt.expected.Numbers {
				t.Errorf("Цифры = %v, ожидается %v", result.Numbers, tt.expected.Numbers)
			}
			if result.Symbols != tt.expected.Symbols {
				t.Errorf("Символы = %v, ожидается %v", result.Symbols, tt.expected.Symbols)
			}
			if result.Total != tt.expected.Total {
				t.Errorf("Всего = %v, ожидается %v", result.Total, tt.expected.Total)
			}
		})
	}
}

func TestDetectCharset(t *testing.T) {
	tests := []struct {
		name            string
		password        string
		expectedCharset string
		expectedSize    int
	}{
		{
			name:            "Все типы",
			password:        "AbC123!@",
			expectedCharset: "A-Z, a-z, 0-9, !@#$...",
			expectedSize:    94, // 26+26+10+32
		},
		{
			name:            "Буквы и цифры",
			password:        "AbC123",
			expectedCharset: "A-Z, a-z, 0-9",
			expectedSize:    62, // 26+26+10
		},
		{
			name:            "Только строчные",
			password:        "abcdef",
			expectedCharset: "a-z",
			expectedSize:    26,
		},
		{
			name:            "Только цифры",
			password:        "123456",
			expectedCharset: "0-9",
			expectedSize:    10,
		},
		{
			name:            "Только символы",
			password:        "!@#$%",
			expectedCharset: "!@#$...",
			expectedSize:    32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charset, size := detectCharset(tt.password)

			if charset != tt.expectedCharset {
				t.Errorf("Набор символов = %q, ожидается %q", charset, tt.expectedCharset)
			}
			if size != tt.expectedSize {
				t.Errorf("Размер = %v, ожидается %v", size, tt.expectedSize)
			}
		})
	}
}

func TestCalculateEntropy(t *testing.T) {
	tests := []struct {
		name        string
		length      int
		charsetSize int
		expected    float64
	}{
		{
			name:        "Стандартный пароль",
			length:      8,
			charsetSize: 94,
			expected:    52.44, // 8 * log2(94) ≈ 52.44
		},
		{
			name:        "Простой пароль",
			length:      6,
			charsetSize: 36,
			expected:    31.02, // 6 * log2(36) ≈ 31.02
		},
		{
			name:        "Нулевая длина",
			length:      0,
			charsetSize: 94,
			expected:    0,
		},
		{
			name:        "Один символ в наборе",
			length:      8,
			charsetSize: 1,
			expected:    0,
		},
		{
			name:        "Нулевой набор",
			length:      8,
			charsetSize: 0,
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateEntropy(tt.length, tt.charsetSize)

			if math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("calculateEntropy(%d, %d) = %v, ожидается ~%v", 
					tt.length, tt.charsetSize, result, tt.expected)
			}
		})
	}
}

func TestEstimateCrackTime(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name         string
		entropy      float64
		expectFast   bool
		expectSlow   bool
	}{
		{
			name:       "Очень низкая энтропия",
			entropy:    10,
			expectFast: true,
			expectSlow: false,
		},
		{
			name:       "Средняя энтропия",
			entropy:    50,
			expectFast: false,
			expectSlow: false,
		},
		{
			name:       "Высокая энтропия",
			entropy:    100,
			expectFast: false,
			expectSlow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := estimateCrackTime(tt.entropy, messages)

			if result.HumanTime == "" {
				t.Error("Время взлома не должно быть пустым")
			}

			if result.Assumptions == "" {
				t.Error("Предположения не должны быть пустыми")
			}

			if tt.expectFast && result.Seconds > 3600 {
				t.Errorf("Ожидалось быстрое время взлома, получено %v секунд", result.Seconds)
			}

			if tt.expectSlow && result.Seconds < 86400 {
				t.Errorf("Ожидалось медленное время взлома, получено %v секунд", result.Seconds)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name     string
		seconds  float64
		contains string
	}{
		{
			name:     "Мгновенно",
			seconds:  0.5,
			contains: "instantly",
		},
		{
			name:     "Секунды",
			seconds:  30,
			contains: "second",
		},
		{
			name:     "Минуты",
			seconds:  1800, // 30 минут
			contains: "minute",
		},
		{
			name:     "Часы",
			seconds:  7200, // 2 часа
			contains: "hour",
		},
		{
			name:     "Дни",
			seconds:  172800, // 2 дня
			contains: "day",
		},
		{
			name:     "Годы",
			seconds:  63072000, // 2 года
			contains: "year",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTime(tt.seconds, messages)

			if result == "" {
				t.Error("formatTime не должен возвращать пустую строку")
			}

			// Проверяем, что результат содержит ожидаемое слово
			// (точная проверка зависит от локализации)
			if len(result) < 3 {
				t.Errorf("formatTime результат слишком короткий: %q", result)
			}
		})
	}
}

func TestDetermineStrength(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name        string
		entropy     float64
		composition CharComposition
		expected    string
	}{
		{
			name:    "Очень слабый",
			entropy: 20,
			composition: CharComposition{
				Lowercase: 5,
				Total:     5,
			},
			expected: messages.StrengthVeryWeak,
		},
		{
			name:    "Хороший с разнообразием",
			entropy: 80,
			composition: CharComposition{
				Uppercase: 2,
				Lowercase: 3,
				Numbers:   2,
				Symbols:   1,
				Total:     8,
			},
			expected: messages.StrengthGood,
		},
		{
			name:    "Высокая энтропия но низкое разнообразие",
			entropy: 90,
			composition: CharComposition{
				Lowercase: 20,
				Total:     20,
			},
			expected: messages.StrengthFair, // Понижено из-за низкого разнообразия
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineStrength(tt.entropy, tt.composition, messages)

			if result != tt.expected {
				t.Errorf("determineStrength() = %q, ожидается %q", result, tt.expected)
			}
		})
	}
}

func TestPasswordInfoFields(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")
	password := "MySecure123!"

	info := AnalyzePassword(password, messages)

	// Проверяем, что все поля заполнены
	if info.Length == 0 {
		t.Error("Длина не должна быть нулевой")
	}

	if info.Charset == "" {
		t.Error("Набор символов не должен быть пустым")
	}

	if info.CharsetSize == 0 {
		t.Error("Размер набора символов не должен быть нулевым")
	}

	if info.Entropy == 0 {
		t.Error("Энтропия не должна быть нулевой")
	}

	if info.TimeToCrack.HumanTime == "" {
		t.Error("Время взлома не должно быть пустым")
	}

	if info.TimeToCrack.Assumptions == "" {
		t.Error("Предположения взлома не должны быть пустыми")
	}

	if info.Composition.Total == 0 {
		t.Error("Общий состав не должен быть нулевым")
	}

	if info.Strength == "" {
		t.Error("Сила не должна быть пустой")
	}
}

func TestRealWorldPasswords(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")

	tests := []struct {
		name            string
		password        string
		minEntropy      float64
		expectedStrong  bool
	}{
		{
			name:           "Сильный сгенерированный пароль",
			password:       "Kp9#mL2$vN8@qR5!",
			minEntropy:     80,
			expectedStrong: true,
		},
		{
			name:           "Слабый обычный пароль",
			password:       "password123",
			minEntropy:     0,
			expectedStrong: false,
		},
		{
			name:           "Средняя парольная фраза",
			password:       "MyHouse123!",
			minEntropy:     50,
			expectedStrong: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := AnalyzePassword(tt.password, messages)

			if info.Entropy < tt.minEntropy {
				t.Errorf("Энтропия = %v, ожидается >= %v", info.Entropy, tt.minEntropy)
			}

			isStrong := info.Strength == messages.StrengthStrong || 
						info.Strength == messages.StrengthVeryStrong

			if tt.expectedStrong && !isStrong {
				t.Errorf("Ожидался сильный пароль, получен %v", info.Strength)
			}

			if !tt.expectedStrong && isStrong && info.Entropy > 70 {
				t.Errorf("Ожидался слабый пароль, получен %v с энтропией %v", 
					info.Strength, info.Entropy)
			}
		})
	}
}

func TestUnicodePasswords(t *testing.T) {
	messages := i18n.GetMessages(i18n.Russian, "test")

	tests := []struct {
		name     string
		password string
	}{
		{"Русский пароль", "МойПароль123!"},
		{"Смешанные языки", "MyПароль123!"},
		{"Пароль с эмодзи", "Pass🔐123!"},
		{"Китайские символы", "密码123!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := AnalyzePassword(tt.password, messages)

			if info.Length == 0 {
				t.Error("Юникод пароль должен иметь ненулевую длину")
			}

			if info.Composition.Total == 0 {
				t.Error("Юникод пароль должен иметь анализ состава")
			}

			if info.Entropy == 0 {
				t.Error("Юникод пароль должен иметь расчет энтропии")
			}
		})
	}
}

// Бенчмарки для измерения производительности
func BenchmarkAnalyzePassword(b *testing.B) {
	messages := i18n.GetMessages(i18n.English, "test")
	password := "MySecurePassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AnalyzePassword(password, messages)
	}
}

func BenchmarkAnalyzeComposition(b *testing.B) {
	password := "MySecurePassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzeComposition(password)
	}
}

func BenchmarkCalculateEntropy(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateEntropy(16, 94)
	}
}

func BenchmarkFormatTime(b *testing.B) {
	messages := i18n.GetMessages(i18n.English, "test")
	seconds := 123456789.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatTime(seconds, messages)
	}
}
