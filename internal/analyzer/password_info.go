package analyzer

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/MaksymLeiber/pgen/internal/i18n"
)

// PasswordInfo информация о сгенерированном пароле
type PasswordInfo struct {
	Length      int
	Charset     string
	CharsetSize int
	Entropy     float64
	TimeToCrack CrackTime
	Composition CharComposition
	Strength    string
}

// CharComposition состав символов в пароле
type CharComposition struct {
	Uppercase int
	Lowercase int
	Numbers   int
	Symbols   int
	Total     int
}

// CrackTime время взлома пароля
type CrackTime struct {
	Seconds     float64
	HumanTime   string
	Assumptions string
}

// AnalyzePassword анализирует сгенерированный пароль
func AnalyzePassword(password string, messages *i18n.Messages) *PasswordInfo {
	info := &PasswordInfo{
		Length:      len(password),
		Composition: analyzeComposition(password),
	}

	// Определяем используемый набор символов
	info.Charset, info.CharsetSize = detectCharset(password)

	// Вычисляем энтропию
	info.Entropy = calculateEntropy(info.Length, info.CharsetSize)

	// Оцениваем время взлома
	info.TimeToCrack = estimateCrackTime(info.Entropy, messages)

	// Определяем общую силу
	info.Strength = determineStrength(info.Entropy, info.Composition, messages)

	return info
}

// analyzeComposition анализирует состав символов
func analyzeComposition(password string) CharComposition {
	comp := CharComposition{}

	for _, char := range password {
		comp.Total++
		switch {
		case unicode.IsUpper(char):
			comp.Uppercase++
		case unicode.IsLower(char):
			comp.Lowercase++
		case unicode.IsNumber(char):
			comp.Numbers++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			comp.Symbols++
		}
	}

	return comp
}

// detectCharset определяет используемый набор символов
func detectCharset(password string) (string, int) {
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSymbol := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	var charsetParts []string
	charsetSize := 0

	if hasUpper {
		charsetParts = append(charsetParts, "A-Z")
		charsetSize += 26
	}
	if hasLower {
		charsetParts = append(charsetParts, "a-z")
		charsetSize += 26
	}
	if hasNumber {
		charsetParts = append(charsetParts, "0-9")
		charsetSize += 10
	}
	if hasSymbol {
		charsetParts = append(charsetParts, "!@#$...")
		charsetSize += 32 // Приблизительно
	}

	charset := strings.Join(charsetParts, ", ")
	return charset, charsetSize
}

// calculateEntropy вычисляет энтропию пароля
func calculateEntropy(length, charsetSize int) float64 {
	if charsetSize <= 1 || length <= 0 {
		return 0
	}
	return float64(length) * math.Log2(float64(charsetSize))
}

// estimateCrackTime оценивает время взлома
func estimateCrackTime(entropy float64, messages *i18n.Messages) CrackTime {
	// Предполагаем современную GPU (RTX 4090) - ~100 млрд попыток в секунду для простых хешей
	// Для Argon2 будет значительно медленнее - примерно 1000 попыток в секунду
	attemptsPerSecond := 1000.0

	// Среднее время взлома = половина от всех возможных комбинаций
	totalCombinations := math.Pow(2, entropy)
	averageAttempts := totalCombinations / 2
	seconds := averageAttempts / attemptsPerSecond

	humanTime := formatTime(seconds, messages)
	assumptions := messages.CrackAssumptions

	return CrackTime{
		Seconds:     seconds,
		HumanTime:   humanTime,
		Assumptions: assumptions,
	}
}

// formatTime форматирует время в читаемый вид
func formatTime(seconds float64, messages *i18n.Messages) string {
	if seconds < 1 {
		return messages.TimeInstantly
	}
	if seconds < 60 {
		return fmt.Sprintf(messages.TimeSeconds, seconds)
	}
	if seconds < 3600 {
		minutes := seconds / 60
		return fmt.Sprintf(messages.TimeMinutes, minutes)
	}
	if seconds < 86400 {
		hours := seconds / 3600
		return fmt.Sprintf(messages.TimeHours, hours)
	}
	if seconds < 31536000 {
		days := seconds / 86400
		return fmt.Sprintf(messages.TimeDays, days)
	}
	if seconds < 31536000000 {
		years := seconds / 31536000
		if years < 1000 {
			return fmt.Sprintf(messages.TimeYears, years)
		}
		if years < 1000000 {
			thousands := years / 1000
			return fmt.Sprintf(messages.TimeThousandYears, thousands)
		}
		if years < 1000000000 {
			millions := years / 1000000
			return fmt.Sprintf(messages.TimeMillionYears, millions)
		}
		billions := years / 1000000000
		return fmt.Sprintf(messages.TimeBillionYears, billions)
	}
	return messages.TimeMoreThanTrillion
}

// determineStrength определяет общую силу пароля
func determineStrength(entropy float64, comp CharComposition, messages *i18n.Messages) string {
	// Базовая оценка по энтропии
	var strengthByEntropy string
	switch {
	case entropy < 30:
		strengthByEntropy = messages.StrengthVeryWeak
	case entropy < 50:
		strengthByEntropy = messages.StrengthWeak
	case entropy < 70:
		strengthByEntropy = messages.StrengthFair
	case entropy < 90:
		strengthByEntropy = messages.StrengthGood
	default:
		strengthByEntropy = messages.StrengthVeryStrong
	}

	// Корректировка на основе разнообразия символов
	diversity := 0
	if comp.Uppercase > 0 {
		diversity++
	}
	if comp.Lowercase > 0 {
		diversity++
	}
	if comp.Numbers > 0 {
		diversity++
	}
	if comp.Symbols > 0 {
		diversity++
	}

	if diversity >= 3 {
		return strengthByEntropy
	} else if diversity == 2 {
		// Понижаем на один уровень
		switch strengthByEntropy {
		case messages.StrengthVeryStrong:
			return messages.StrengthStrong
		case messages.StrengthStrong:
			return messages.StrengthFair
		case messages.StrengthFair:
			return messages.StrengthWeak
		default:
			return strengthByEntropy
		}
	} else {
		// Понижаем на два уровня
		switch strengthByEntropy {
		case messages.StrengthVeryStrong:
			return messages.StrengthFair
		case messages.StrengthStrong:
			return messages.StrengthWeak
		case messages.StrengthFair:
			return messages.StrengthVeryWeak
		default:
			return messages.StrengthVeryWeak
		}
	}
}
