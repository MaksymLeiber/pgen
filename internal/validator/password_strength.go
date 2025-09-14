package validator

import (
	"regexp"
	"strings"
	"unicode"

	"pgen/internal/i18n"
)

// StrengthLevel уровень силы пароля
type StrengthLevel int

const (
	StrengthWeak StrengthLevel = iota
	StrengthFair
	StrengthGood
	StrengthStrong
	StrengthVeryStrong
)

// PasswordStrength результат анализа силы пароля
type PasswordStrength struct {
	Level       StrengthLevel
	Score       int
	Issues      []string
	Suggestions []string
}

// ValidatePasswordStrength анализирует силу мастер-пароля
func ValidatePasswordStrength(password string, messages *i18n.Messages) *PasswordStrength {
	result := &PasswordStrength{
		Score:       0,
		Issues:      []string{},
		Suggestions: []string{},
	}

	// Базовые проверки
	length := len(password)

	// Длина пароля
	if length < 8 {
		result.Issues = append(result.Issues, messages.Errors.PasswordIssues.LengthTooShort)
		result.Suggestions = append(result.Suggestions, messages.Errors.Suggestions.IncreaseLength)
	} else if length >= 8 && length < 12 {
		result.Score += 10
	} else if length >= 12 && length < 16 {
		result.Score += 20
	} else if length >= 16 {
		result.Score += 25
	}

	// Типы символов
	hasLower := false
	hasUpper := false
	hasNumber := false
	hasSymbol := false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	charTypeCount := 0
	if hasLower {
		charTypeCount++
		result.Score += 5
	} else {
		result.Issues = append(result.Issues, messages.Errors.PasswordIssues.NoLowercase)
		result.Suggestions = append(result.Suggestions, messages.Errors.Suggestions.AddLowercase)
	}

	if hasUpper {
		charTypeCount++
		result.Score += 5
	} else {
		result.Issues = append(result.Issues, messages.Errors.PasswordIssues.NoUppercase)
		result.Suggestions = append(result.Suggestions, messages.Errors.Suggestions.AddUppercase)
	}

	if hasNumber {
		charTypeCount++
		result.Score += 5
	} else {
		result.Issues = append(result.Issues, messages.Errors.PasswordIssues.NoNumbers)
		result.Suggestions = append(result.Suggestions, messages.Errors.Suggestions.AddNumbers)
	}

	if hasSymbol {
		charTypeCount++
		result.Score += 10
	} else {
		result.Suggestions = append(result.Suggestions, messages.Errors.Suggestions.AddSymbols)
	}

	// Бонус за разнообразие символов
	if charTypeCount >= 3 {
		result.Score += 10
	}
	if charTypeCount == 4 {
		result.Score += 15
	}

	// Проверка на повторяющиеся символы
	if hasRepeatingChars(password) {
		result.Issues = append(result.Issues, messages.Errors.PasswordIssues.RepeatingChars)
		result.Suggestions = append(result.Suggestions, messages.Errors.Suggestions.AvoidRepetition)
		result.Score -= 10
	}

	// Проверка на последовательности
	if hasSequences(password) {
		result.Issues = append(result.Issues, messages.Errors.PasswordIssues.SequentialChars)
		result.Suggestions = append(result.Suggestions, messages.Errors.Suggestions.AvoidSequences)
		result.Score -= 15
	}

	// Проверка на словарные слова
	if hasCommonWords(password) {
		result.Issues = append(result.Issues, messages.Errors.PasswordIssues.CommonWords)
		result.Suggestions = append(result.Suggestions, messages.Errors.Suggestions.AvoidDictionary)
		result.Score -= 20
	}

	// Проверка на фразы (предпочтительно)
	if isPhrase(password) {
		result.Score += 20
	}

	// Определение уровня
	if result.Score < 30 {
		result.Level = StrengthWeak
	} else if result.Score < 50 {
		result.Level = StrengthFair
	} else if result.Score < 70 {
		result.Level = StrengthGood
	} else if result.Score < 90 {
		result.Level = StrengthStrong
	} else {
		result.Level = StrengthVeryStrong
	}

	return result
}

// hasRepeatingChars проверяет на повторяющиеся символы подряд
func hasRepeatingChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}
	return false
}

// hasSequences проверяет на последовательности символов
func hasSequences(password string) bool {
	sequences := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"qwertyuiopasdfghjklzxcvbnm",
		"0123456789",
	}

	lower := strings.ToLower(password)
	for _, seq := range sequences {
		if containsSequence(lower, seq, 3) || containsSequence(lower, reverse(seq), 3) {
			return true
		}
	}
	return false
}

// containsSequence проверяет наличие последовательности длиной minLen
func containsSequence(password, sequence string, minLen int) bool {
	for i := 0; i <= len(sequence)-minLen; i++ {
		substr := sequence[i : i+minLen]
		if strings.Contains(password, substr) {
			return true
		}
	}
	return false
}

// reverse переворачивает строку
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// hasCommonWords проверяет на наличие словарных слов
func hasCommonWords(password string) bool {
	// Расширенный список частых слов и паттернов
	commonWords := []string{
		// Основные пароли
		"password", "пароль", "pass", "pwd", "passw0rd",
		"123456", "12345678", "123456789", "1234567890",
		"qwerty", "qwertyui", "asdfgh", "zxcvbn",
		"111111", "000000", "654321",
		
		// Системные учетные записи
		"admin", "administrator", "администратор",
		"user", "пользователь", "guest", "гость",
		"root", "рут", "superuser", "супер",
		"test", "тест", "demo", "демо",
		
		// Вход и авторизация
		"login", "логин", "signin", "вход",
		"access", "доступ", "enter", "войти",
		"welcome", "добро пожаловать", "start", "старт",
		
		// Секретность
		"secret", "секрет", "private", "приватный",
		"confidential", "конфиденциально", "hidden", "скрытый",
		"secure", "безопасный", "protect", "защита",
		
		// Личная информация
		"love", "любовь", "family", "семья",
		"money", "деньги", "house", "дом",
		"hello", "привет", "world", "мир",
		"life", "жизнь", "work", "работа",
		
		// Популярные слова
		"computer", "компьютер", "internet", "интернет",
		"email", "почта", "phone", "телефон",
		"birthday", "день рождения", "name", "имя",
		"address", "адрес", "city", "город",
		
		// Клавиатурные паттерны
		"asdf", "hjkl", "wasd", "йцук",
		"фыва", "олдж", "ячсм", "qaz",
		"wsx", "edc", "rfv", "tgb",
		
		// Популярные бренды и сервисы
		"google", "гугл", "apple", "эппл",
		"microsoft", "майкрософт", "windows", "виндовс",
		"facebook", "фейсбук", "twitter", "твиттер",
		"instagram", "инстаграм", "youtube", "ютуб",
		
		// Даты и годы
		"2023", "2024", "2025", "2022", "2021", "2020",
		"january", "январь", "february", "февраль",
		"march", "март", "april", "апрель",
		"monday", "понедельник", "sunday", "воскресенье",
		
		// Простые замены
		"passw0rd", "p@ssword", "p@ssw0rd",
		"adm1n", "@dmin", "r00t", "t3st",
		"s3cret", "l0ve", "h3llo", "w0rld",
	}

	lower := strings.ToLower(password)
	for _, word := range commonWords {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

// isPhrase определяет, является ли пароль фразой (содержит пробелы или много слов)
func isPhrase(password string) bool {
	// Если содержит пробелы или длинный и содержит разнообразные символы
	if strings.Contains(password, " ") {
		return true
	}

	// Или если длинный и содержит много типов символов
	if len(password) >= 15 {
		wordPattern := regexp.MustCompile(`[a-zA-Zа-яА-Я]{3,}`)
		words := wordPattern.FindAllString(password, -1)
		return len(words) >= 2
	}

	return false
}
