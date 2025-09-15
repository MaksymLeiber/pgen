package validator

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/MaksymLeiber/pgen/internal/i18n"
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

// commonWordsMap содержит часто используемые слова и паттерны для быстрого поиска
var commonWordsMap = map[string]bool{
	// Основные пароли
	"password": true, "пароль": true, "pass": true, "pwd": true, "passw0rd": true,
	"123456": true, "12345678": true, "123456789": true, "1234567890": true,
	"qwerty": true, "qwertyui": true, "asdfgh": true, "zxcvbn": true,
	"111111": true, "000000": true, "654321": true,

	// Системные учетные записи
	"admin": true, "administrator": true, "администратор": true,
	"user": true, "пользователь": true, "guest": true, "гость": true,
	"root": true, "рут": true, "superuser": true, "супер": true,
	"test": true, "тест": true, "demo": true, "демо": true,

	// Вход и авторизация
	"login": true, "логин": true, "signin": true, "вход": true,
	"access": true, "доступ": true, "enter": true, "войти": true,
	"welcome": true, "добро пожаловать": true, "start": true, "старт": true,

	// Секретность
	"secret": true, "секрет": true, "private": true, "приватный": true,
	"confidential": true, "конфиденциально": true, "hidden": true, "скрытый": true,
	"secure": true, "безопасный": true, "protect": true, "защита": true,

	// Личная информация
	"love": true, "любовь": true, "family": true, "семья": true,
	"money": true, "деньги": true, "house": true, "дом": true,
	"hello": true, "привет": true, "world": true, "мир": true,
	"life": true, "жизнь": true, "work": true, "работа": true,

	// Популярные слова
	"computer": true, "компьютер": true, "internet": true, "интернет": true,
	"email": true, "почта": true, "phone": true, "телефон": true,
	"birthday": true, "день рождения": true, "name": true, "имя": true,
	"address": true, "адрес": true, "city": true, "город": true,

	// Клавиатурные паттерны
	"asdf": true, "hjkl": true, "wasd": true, "йцук": true,
	"фыва": true, "олдж": true, "ячсм": true, "qaz": true,
	"wsx": true, "edc": true, "rfv": true, "tgb": true,

	// Популярные бренды и сервисы
	"google": true, "гугл": true, "apple": true, "эппл": true,
	"microsoft": true, "майкрософт": true, "windows": true, "виндовс": true,
	"facebook": true, "фейсбук": true, "twitter": true, "твиттер": true,
	"instagram": true, "инстаграм": true, "youtube": true, "ютуб": true,

	// Даты и годы
	"2023": true, "2024": true, "2025": true, "2022": true, "2021": true, "2020": true,
	"january": true, "январь": true, "february": true, "февраль": true,
	"march": true, "март": true, "april": true, "апрель": true,
	"monday": true, "понедельник": true, "sunday": true, "воскресенье": true,

	// Простые замены
	"p@ssword": true, "p@ssw0rd": true,
	"adm1n": true, "@dmin": true, "r00t": true, "t3st": true,
	"s3cret": true, "l0ve": true, "h3llo": true, "w0rld": true,
}

// hasCommonWords проверяет на наличие словарных слов (оптимизированная версия)
func hasCommonWords(password string) bool {
	lower := strings.ToLower(password)

	// Быстрая проверка прямых совпадений в map
	for word := range commonWordsMap {
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
