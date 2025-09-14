package i18n

import (
	"os"
	"strings"
	"testing"
)

func TestLanguageConstants(t *testing.T) {
	// Тест констант языков
	if Russian != "ru" {
		t.Errorf("Russian = %q, ожидается 'ru'", Russian)
	}
	if English != "en" {
		t.Errorf("English = %q, ожидается 'en'", English)
	}
}

func TestGetMessagesRussian(t *testing.T) {
	// Тест получения русских сообщений
	version := "1.0.0"
	messages := GetMessages(Russian, version)

	if messages == nil {
		t.Fatal("GetMessages(Russian) не должен возвращать nil")
	}

	// Проверяем основные сообщения
	tests := []struct {
		name     string
		got      string
		contains string
	}{
		{"EnterMasterPassword", messages.EnterMasterPassword, "мастер-пароль"},
		{"EnterServiceName", messages.EnterServiceName, "сервис"},
		{"PasswordGenerated", messages.PasswordGenerated, "пароль"},
		{"CopiedToClipboard", messages.CopiedToClipboard, "скопирован"},
		{"Version", messages.Version, version},
		{"AppTitle", messages.AppTitle, "PGen"},
		{"Description", messages.Description, "генерации"},
		{"Usage", messages.Usage, "pgen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got == "" {
				t.Errorf("GetMessages(Russian).%s не должно быть пустым", tt.name)
			}
			if !strings.Contains(strings.ToLower(tt.got), strings.ToLower(tt.contains)) {
				t.Errorf("GetMessages(Russian).%s = %q, должно содержать %q", tt.name, tt.got, tt.contains)
			}
		})
	}
}

func TestGetMessagesEnglish(t *testing.T) {
	// Тест получения английских сообщений
	version := "2.0.0"
	messages := GetMessages(English, version)

	if messages == nil {
		t.Fatal("GetMessages(English) не должен возвращать nil")
	}

	// Проверяем основные сообщения
	tests := []struct {
		name     string
		got      string
		contains string
	}{
		{"EnterMasterPassword", messages.EnterMasterPassword, "master password"},
		{"EnterServiceName", messages.EnterServiceName, "service"},
		{"PasswordGenerated", messages.PasswordGenerated, "password"},
		{"CopiedToClipboard", messages.CopiedToClipboard, "copied"},
		{"Version", messages.Version, version},
		{"AppTitle", messages.AppTitle, "PGen"},
		{"Description", messages.Description, "generating"},
		{"Usage", messages.Usage, "pgen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got == "" {
				t.Errorf("GetMessages(English).%s не должно быть пустым", tt.name)
			}
			if !strings.Contains(strings.ToLower(tt.got), strings.ToLower(tt.contains)) {
				t.Errorf("GetMessages(English).%s = %q, должно содержать %q", tt.name, tt.got, tt.contains)
			}
		})
	}
}

func TestGetMessagesDefaultLanguage(t *testing.T) {
	// Тест получения сообщений для неизвестного языка (должен вернуть английский)
	version := "1.5.0"
	messages := GetMessages("unknown", version)

	if messages == nil {
		t.Fatal("GetMessages('unknown') не должен возвращать nil")
	}

	// Должен вернуть английские сообщения
	if !strings.Contains(messages.EnterMasterPassword, "master password") {
		t.Errorf("GetMessages('unknown') должен возвращать английские сообщения, получено: %q", messages.EnterMasterPassword)
	}

	if !strings.Contains(messages.Version, version) {
		t.Errorf("GetMessages('unknown').Version должен содержать версию %q, получено: %q", version, messages.Version)
	}
}

func TestGetMessagesVersionIntegration(t *testing.T) {
	// Тест интеграции версии в сообщения
	testVersions := []string{"1.0.0", "2.5.3", "0.1.0-beta", "3.0.0-rc1"}

	for _, version := range testVersions {
		t.Run("Version_"+version, func(t *testing.T) {
			ruMessages := GetMessages(Russian, version)
			enMessages := GetMessages(English, version)

			if !strings.Contains(ruMessages.Version, version) {
				t.Errorf("Русские сообщения должны содержать версию %q, получено: %q", version, ruMessages.Version)
			}

			if !strings.Contains(enMessages.Version, version) {
				t.Errorf("Английские сообщения должны содержать версию %q, получено: %q", version, enMessages.Version)
			}
		})
	}
}

func TestPasswordStrengthMessages(t *testing.T) {
	// Тест сообщений о силе пароля
	ruMessages := GetMessages(Russian, "1.0.0")
	enMessages := GetMessages(English, "1.0.0")

	// Проверяем русские сообщения о силе
	ruStrengthTests := []struct {
		name string
		msg  string
	}{
		{"PasswordStrengthWeak", ruMessages.PasswordStrengthWeak},
		{"PasswordStrengthFair", ruMessages.PasswordStrengthFair},
		{"PasswordStrengthGood", ruMessages.PasswordStrengthGood},
		{"PasswordStrengthStrong", ruMessages.PasswordStrengthStrong},
		{"PasswordStrengthVeryStrong", ruMessages.PasswordStrengthVeryStrong},
	}

	for _, tt := range ruStrengthTests {
		t.Run("Russian_"+tt.name, func(t *testing.T) {
			if tt.msg == "" {
				t.Errorf("Русское сообщение %s не должно быть пустым", tt.name)
			}
			if !strings.Contains(tt.msg, "🔴") && !strings.Contains(tt.msg, "🟠") && 
			   !strings.Contains(tt.msg, "🟡") && !strings.Contains(tt.msg, "🟢") {
				// Проверяем, что есть хотя бы текст
				if len(tt.msg) < 3 {
					t.Errorf("Русское сообщение %s слишком короткое: %q", tt.name, tt.msg)
				}
			}
		})
	}

	// Проверяем английские сообщения о силе
	enStrengthTests := []struct {
		name string
		msg  string
	}{
		{"PasswordStrengthWeak", enMessages.PasswordStrengthWeak},
		{"PasswordStrengthFair", enMessages.PasswordStrengthFair},
		{"PasswordStrengthGood", enMessages.PasswordStrengthGood},
		{"PasswordStrengthStrong", enMessages.PasswordStrengthStrong},
		{"PasswordStrengthVeryStrong", enMessages.PasswordStrengthVeryStrong},
	}

	for _, tt := range enStrengthTests {
		t.Run("English_"+tt.name, func(t *testing.T) {
			if tt.msg == "" {
				t.Errorf("Английское сообщение %s не должно быть пустым", tt.name)
			}
			if len(tt.msg) < 3 {
				t.Errorf("Английское сообщение %s слишком короткое: %q", tt.name, tt.msg)
			}
		})
	}
}

func TestConfigurationMessages(t *testing.T) {
	// Тест сообщений конфигурации
	ruMessages := GetMessages(Russian, "1.0.0")
	enMessages := GetMessages(English, "1.0.0")

	configTests := []struct {
		name string
		ruMsg string
		enMsg string
	}{
		{"ConfigErrorReading", ruMessages.ConfigErrorReading, enMessages.ConfigErrorReading},
		{"ConfigErrorParsing", ruMessages.ConfigErrorParsing, enMessages.ConfigErrorParsing},
		{"ConfigErrorWriting", ruMessages.ConfigErrorWriting, enMessages.ConfigErrorWriting},
		{"ConfigShow", ruMessages.ConfigShow, enMessages.ConfigShow},
		{"ConfigSet", ruMessages.ConfigSet, enMessages.ConfigSet},
		{"ConfigReset", ruMessages.ConfigReset, enMessages.ConfigReset},
	}

	for _, tt := range configTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ruMsg == "" {
				t.Errorf("Русское сообщение %s не должно быть пустым", tt.name)
			}
			if tt.enMsg == "" {
				t.Errorf("Английское сообщение %s не должно быть пустым", tt.name)
			}
			
			// Проверяем, что сообщения разные (не одинаковые)
			if tt.ruMsg == tt.enMsg {
				t.Errorf("Русское и английское сообщения %s не должны быть одинаковыми: %q", tt.name, tt.ruMsg)
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	// Тест сообщений об ошибках
	ruMessages := GetMessages(Russian, "1.0.0")
	enMessages := GetMessages(English, "1.0.0")

	// Проверяем основные ошибки
	errorTests := []struct {
		name string
		ruMsg string
		enMsg string
	}{
		{"ClipboardError", ruMessages.Errors.ClipboardError, enMessages.Errors.ClipboardError},
		{"GenerationError", ruMessages.Errors.GenerationError, enMessages.Errors.GenerationError},
		{"EmptyMaster", ruMessages.Errors.EmptyMaster, enMessages.Errors.EmptyMaster},
		{"EmptyService", ruMessages.Errors.EmptyService, enMessages.Errors.EmptyService},
		{"UserCanceled", ruMessages.Errors.UserCanceled, enMessages.Errors.UserCanceled},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ruMsg == "" {
				t.Errorf("Русское сообщение об ошибке %s не должно быть пустым", tt.name)
			}
			if tt.enMsg == "" {
				t.Errorf("Английское сообщение об ошибке %s не должно быть пустым", tt.name)
			}
		})
	}

	// Проверяем вложенные структуры ошибок
	passwordIssueTests := []struct {
		name string
		ruMsg string
		enMsg string
	}{
		{"LengthTooShort", ruMessages.Errors.PasswordIssues.LengthTooShort, enMessages.Errors.PasswordIssues.LengthTooShort},
		{"NoLowercase", ruMessages.Errors.PasswordIssues.NoLowercase, enMessages.Errors.PasswordIssues.NoLowercase},
		{"NoUppercase", ruMessages.Errors.PasswordIssues.NoUppercase, enMessages.Errors.PasswordIssues.NoUppercase},
		{"NoNumbers", ruMessages.Errors.PasswordIssues.NoNumbers, enMessages.Errors.PasswordIssues.NoNumbers},
	}

	for _, tt := range passwordIssueTests {
		t.Run("PasswordIssues_"+tt.name, func(t *testing.T) {
			if tt.ruMsg == "" {
				t.Errorf("Русское сообщение о проблеме пароля %s не должно быть пустым", tt.name)
			}
			if tt.enMsg == "" {
				t.Errorf("Английское сообщение о проблеме пароля %s не должно быть пустым", tt.name)
			}
		})
	}
}

func TestFlagsMessages(t *testing.T) {
	// Тест сообщений флагов
	ruMessages := GetMessages(Russian, "1.0.0")
	enMessages := GetMessages(English, "1.0.0")

	flagTests := []struct {
		name string
		ruMsg string
		enMsg string
	}{
		{"Lang", ruMessages.Flags.Lang, enMessages.Flags.Lang},
		{"LangDesc", ruMessages.Flags.LangDesc, enMessages.Flags.LangDesc},
		{"Length", ruMessages.Flags.Length, enMessages.Flags.Length},
		{"LengthDesc", ruMessages.Flags.LengthDesc, enMessages.Flags.LengthDesc},
		{"Copy", ruMessages.Flags.Copy, enMessages.Flags.Copy},
		{"CopyDesc", ruMessages.Flags.CopyDesc, enMessages.Flags.CopyDesc},
	}

	for _, tt := range flagTests {
		t.Run("Flags_"+tt.name, func(t *testing.T) {
			if tt.ruMsg == "" {
				t.Errorf("Русское сообщение флага %s не должно быть пустым", tt.name)
			}
			if tt.enMsg == "" {
				t.Errorf("Английское сообщение флага %s не должно быть пустым", tt.name)
			}
		})
	}
}

func TestDetectLanguageWithFlag(t *testing.T) {
	// Тест определения языка по флагу
	tests := []struct {
		name     string
		langFlag string
		expected Language
	}{
		{"Russian_ru", "ru", Russian},
		{"Russian_russian", "russian", Russian},
		{"Russian_RU", "RU", Russian},
		{"Russian_RUSSIAN", "RUSSIAN", Russian},
		{"English_en", "en", English},
		{"English_english", "english", English},
		{"English_EN", "EN", English},
		{"English_ENGLISH", "ENGLISH", English},
		{"Unknown_de", "de", English}, // Неизвестный язык должен вернуть English
		{"Unknown_fr", "fr", English},
		{"Empty", "", English}, // Пустой флаг должен определяться по окружению
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем оригинальные переменные окружения
			originalLang := os.Getenv("LANG")
			originalLCAll := os.Getenv("LC_ALL")
			originalLCMessages := os.Getenv("LC_MESSAGES")

			defer func() {
				os.Setenv("LANG", originalLang)
				os.Setenv("LC_ALL", originalLCAll)
				os.Setenv("LC_MESSAGES", originalLCMessages)
			}()

			// Очищаем переменные окружения для чистого теста
			os.Setenv("LANG", "")
			os.Setenv("LC_ALL", "")
			os.Setenv("LC_MESSAGES", "")

			result := DetectLanguage(tt.langFlag)
			if result != tt.expected {
				t.Errorf("DetectLanguage(%q) = %q, ожидается %q", tt.langFlag, result, tt.expected)
			}
		})
	}
}

func TestDetectLanguageFromEnvironment(t *testing.T) {
	// Тест определения языка из переменных окружения
	tests := []struct {
		name     string
		langEnv  string
		lcAllEnv string
		lcMsgEnv string
		expected Language
	}{
		{"LANG_ru", "ru_RU.UTF-8", "", "", Russian},
		{"LANG_russian", "russian", "", "", Russian},
		{"LANG_en", "en_US.UTF-8", "", "", English},
		{"LC_ALL_ru", "", "ru_RU.UTF-8", "", Russian},
		{"LC_MESSAGES_ru", "", "", "ru_RU.UTF-8", Russian},
		{"Multiple_ru", "", "ru_RU.UTF-8", "", Russian}, // LC_ALL проверяется только если LANG пустой
		{"No_ru_anywhere", "en_US.UTF-8", "de_DE.UTF-8", "fr_FR.UTF-8", English},
		{"Empty_all", "", "", "", English},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем оригинальные переменные окружения
			originalLang := os.Getenv("LANG")
			originalLCAll := os.Getenv("LC_ALL")
			originalLCMessages := os.Getenv("LC_MESSAGES")

			defer func() {
				os.Setenv("LANG", originalLang)
				os.Setenv("LC_ALL", originalLCAll)
				os.Setenv("LC_MESSAGES", originalLCMessages)
			}()

			// Устанавливаем тестовые переменные
			os.Setenv("LANG", tt.langEnv)
			os.Setenv("LC_ALL", tt.lcAllEnv)
			os.Setenv("LC_MESSAGES", tt.lcMsgEnv)

			result := DetectLanguage("") // Пустой флаг для проверки окружения
			if result != tt.expected {
				t.Errorf("DetectLanguage('') с окружением LANG=%q, LC_ALL=%q, LC_MESSAGES=%q = %q, ожидается %q", 
					tt.langEnv, tt.lcAllEnv, tt.lcMsgEnv, result, tt.expected)
			}
		})
	}
}

func TestDetectLanguageFlagPriority(t *testing.T) {
	// Тест приоритета флага над переменными окружения
	originalLang := os.Getenv("LANG")
	defer os.Setenv("LANG", originalLang)

	// Устанавливаем русское окружение
	os.Setenv("LANG", "ru_RU.UTF-8")

	// Флаг должен иметь приоритет над окружением
	result := DetectLanguage("en")
	if result != English {
		t.Errorf("DetectLanguage('en') с русским окружением должен возвращать English, получено %q", result)
	}

	// И наоборот
	os.Setenv("LANG", "en_US.UTF-8")
	result = DetectLanguage("ru")
	if result != Russian {
		t.Errorf("DetectLanguage('ru') с английским окружением должен возвращать Russian, получено %q", result)
	}
}

func TestInstallationMessages(t *testing.T) {
	// Тест сообщений установки
	ruMessages := GetMessages(Russian, "1.0.0")
	enMessages := GetMessages(English, "1.0.0")

	installTests := []struct {
		name string
		ruMsg string
		enMsg string
	}{
		{"InstallSuccess", ruMessages.InstallSuccess, enMessages.InstallSuccess},
		{"InstallError", ruMessages.InstallError, enMessages.InstallError},
		{"InstallAlreadyExists", ruMessages.InstallAlreadyExists, enMessages.InstallAlreadyExists},
		{"UninstallSuccess", ruMessages.UninstallSuccess, enMessages.UninstallSuccess},
		{"UninstallError", ruMessages.UninstallError, enMessages.UninstallError},
		{"UninstallNotInstalled", ruMessages.UninstallNotInstalled, enMessages.UninstallNotInstalled},
	}

	for _, tt := range installTests {
		t.Run("Install_"+tt.name, func(t *testing.T) {
			if tt.ruMsg == "" {
				t.Errorf("Русское сообщение установки %s не должно быть пустым", tt.name)
			}
			if tt.enMsg == "" {
				t.Errorf("Английское сообщение установки %s не должно быть пустым", tt.name)
			}
		})
	}
}

func TestTimeMessages(t *testing.T) {
	// Тест сообщений времени взлома
	ruMessages := GetMessages(Russian, "1.0.0")
	enMessages := GetMessages(English, "1.0.0")

	timeTests := []struct {
		name string
		ruMsg string
		enMsg string
	}{
		{"TimeInstantly", ruMessages.TimeInstantly, enMessages.TimeInstantly},
		{"TimeSeconds", ruMessages.TimeSeconds, enMessages.TimeSeconds},
		{"TimeMinutes", ruMessages.TimeMinutes, enMessages.TimeMinutes},
		{"TimeHours", ruMessages.TimeHours, enMessages.TimeHours},
		{"TimeDays", ruMessages.TimeDays, enMessages.TimeDays},
		{"TimeYears", ruMessages.TimeYears, enMessages.TimeYears},
	}

	for _, tt := range timeTests {
		t.Run("Time_"+tt.name, func(t *testing.T) {
			if tt.ruMsg == "" {
				t.Errorf("Русское сообщение времени %s не должно быть пустым", tt.name)
			}
			if tt.enMsg == "" {
				t.Errorf("Английское сообщение времени %s не должно быть пустым", tt.name)
			}

			// Проверяем, что сообщения времени содержат форматирование для чисел
			if strings.Contains(tt.name, "Seconds") || strings.Contains(tt.name, "Minutes") || 
			   strings.Contains(tt.name, "Hours") || strings.Contains(tt.name, "Days") || 
			   strings.Contains(tt.name, "Years") {
				if !strings.Contains(tt.ruMsg, "%.0f") {
					t.Errorf("Русское сообщение времени %s должно содержать форматирование %%.0f, получено: %q", tt.name, tt.ruMsg)
				}
				if !strings.Contains(tt.enMsg, "%.0f") {
					t.Errorf("Английское сообщение времени %s должно содержать форматирование %%.0f, получено: %q", tt.name, tt.enMsg)
				}
			}
		})
	}
}

func TestMessageConsistency(t *testing.T) {
	// Тест консистентности сообщений между языками
	ruMessages := GetMessages(Russian, "1.0.0")
	enMessages := GetMessages(English, "1.0.0")

	// Проверяем, что все основные поля заполнены в обоих языках
	if ruMessages.EnterMasterPassword == "" || enMessages.EnterMasterPassword == "" {
		t.Error("EnterMasterPassword должно быть заполнено в обоих языках")
	}

	if ruMessages.AppTitle == "" || enMessages.AppTitle == "" {
		t.Error("AppTitle должно быть заполнено в обоих языках")
	}

	if ruMessages.Examples == "" || enMessages.Examples == "" {
		t.Error("Examples должно быть заполнено в обоих языках")
	}

	// Проверяем, что структуры флагов заполнены
	if ruMessages.Flags.Lang == "" || enMessages.Flags.Lang == "" {
		t.Error("Flags.Lang должно быть заполнено в обоих языках")
	}

	// Проверяем, что структуры ошибок заполнены
	if ruMessages.Errors.EmptyMaster == "" || enMessages.Errors.EmptyMaster == "" {
		t.Error("Errors.EmptyMaster должно быть заполнено в обоих языках")
	}
}

// Бенчмарки для измерения производительности
func BenchmarkGetMessagesRussian(b *testing.B) {
	version := "1.0.0"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetMessages(Russian, version)
	}
}

func BenchmarkGetMessagesEnglish(b *testing.B) {
	version := "1.0.0"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetMessages(English, version)
	}
}

func BenchmarkDetectLanguage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DetectLanguage("ru")
	}
}

func BenchmarkDetectLanguageFromEnv(b *testing.B) {
	os.Setenv("LANG", "ru_RU.UTF-8")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DetectLanguage("")
	}
}
