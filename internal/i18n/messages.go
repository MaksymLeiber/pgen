package i18n

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

type Language string

const (
	Russian Language = "ru"
	English Language = "en"
)

func (m *Messages) GetRandomTip() string {
	if len(m.Tips) == 0 {
		return m.Tip // Фолбэк к обычному совету
	}

	// генератор случайных чисел на основе времени для подсказок
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(m.Tips))
	return m.Tips[index]
}

type Messages struct {
	EnterMasterPassword   string
	EnterServiceName      string
	PasswordGenerated     string
	CopiedToClipboard     string
	ClipboardCleared      string
	ClipboardWillClear    string
	ClipboardSecurityInfo string
	ClipboardWaitingInfo  string
	Version               string
	Usage                 string
	Description           string
	Examples              string
	AppTitle              string
	AppSubtitle           string
	GeneratingPassword    string
	LengthLabel           string
	CharactersLabel       string
	Tip                   string
	Tips                  []string
	About                 string

	// Установка приложения
	InstallSuccess          string
	InstallError            string
	InstallAlreadyExists    string
	InstallPermissionDenied string
	InstallPathAdded        string
	InstallComplete         string
	InstallCheckingPath     string
	InstallCopyingFile      string
	InstallUpdatingPath     string
	InstallRestart          string
	InstallLocation         string
	InstallInstalledTo      string

	// Технические ошибки установки
	InstallErrorCreateDir   string
	InstallErrorGetExePath  string
	InstallErrorCopyFile    string
	InstallErrorSetPerms    string
	InstallErrorAddPath     string
	InstallErrorOpenFile    string
	InstallErrorWriteFile   string
	InstallErrorResolvePath string

	// Служебные сообщения установщика
	InstallProfileComment   string
	InstallPanicWindowsFunc string
	InstallPanicUnixFunc    string

	// Удаление приложения
	UninstallSuccess          string
	UninstallError            string
	UninstallNotInstalled     string
	UninstallPermissionDenied string
	UninstallPathRemoved      string
	UninstallComplete         string
	UninstallCheckingPath     string
	UninstallRemoving         string
	UninstallConfirm          string
	UninstallCancelled        string
	UninstallRemovedFrom      string

	// Валидация мастер-пароля
	MasterPasswordStrength     string
	PasswordStrengthWeak       string
	PasswordStrengthFair       string
	PasswordStrengthGood       string
	PasswordStrengthStrong     string
	PasswordStrengthVeryStrong string

	// Информация о пароле
	PasswordInfo     string
	Charset          string
	Entropy          string
	TimeToCrack      string
	Composition      string
	Uppercase        string
	Lowercase        string
	Numbers          string
	Symbols          string
	Strength         string
	CrackAssumptions string

	// Единицы времени для взлома
	TimeInstantly        string
	TimeSeconds          string
	TimeMinutes          string
	TimeHours            string
	TimeDays             string
	TimeYears            string
	TimeThousandYears    string
	TimeMillionYears     string
	TimeBillionYears     string
	TimeMoreThanTrillion string

	// Уровни силы пароля (без эмодзи для использования в коде)
	StrengthVeryWeak   string
	StrengthWeak       string
	StrengthFair       string
	StrengthGood       string
	StrengthStrong     string
	StrengthVeryStrong string

	// Configuration error messages
	ConfigErrorReading      string
	ConfigErrorParsing      string
	ConfigErrorEncoding     string
	ConfigErrorWriting      string
	ConfigErrorExportEncode string
	ConfigErrorExportWrite  string
	ConfigErrorImportRead   string
	ConfigErrorImportParse  string

	// Команды конфигурации
	ConfigShow                 string
	ConfigSet                  string
	ConfigReset                string
	ConfigExport               string
	ConfigImport               string
	ConfigCurrentConfig        string
	ConfigUpdated              string
	ConfigReset_               string
	ConfigExported             string
	ConfigImported             string
	ConfigInvalidKey           string
	ConfigManagement           string
	ConfigManageSettings       string
	ConfigErrorSaving          string
	ConfigErrorExporting       string
	ConfigErrorImporting       string
	ConfigErrorSavingImported  string
	ConfigInvalidArgonTime     string
	ConfigInvalidArgonMemory   string
	ConfigInvalidArgonThreads  string
	ConfigInvalidArgonKeyLen   string
	ConfigInvalidDefaultLength string
	ConfigInvalidDefaultLang   string
	ConfigInvalidCharset       string
	ConfigInvalidDefaultCopy   string
	ConfigInvalidClearTimeout  string
	ConfigInvalidPasswordInfo  string
	ConfigInvalidColorOutput   string
	ConfigInvalidUsername      string
	ConfigUsernameEmpty        string
	ConfigUnknownKey           string
	ProfileLabel               string
	ConfigLengthRange          string
	ConfigLanguageValues       string
	ConfigCharsetValues        string
	ConfigTimeoutRange         string

	Flags struct {
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
	}

	Errors struct {
		ClipboardError  string
		GenerationError string
		EmptyMaster     string
		EmptyService    string
		UserCanceled    string
		InputCanceled   string
		LengthTooShort  string
		LengthTooLong   string
		HashTooShort    string

		// Проблемы с мастер-паролем
		PasswordIssues struct {
			LengthTooShort  string
			NoLowercase     string
			NoUppercase     string
			NoNumbers       string
			RepeatingChars  string
			SequentialChars string
			CommonWords     string
		}

		// Рекомендации по улучшению
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
	}
}

func GetMessages(lang Language, version string) *Messages {
	switch lang {
	case Russian:
		return &Messages{
			EnterMasterPassword:   "Введите мастер-пароль:",
			EnterServiceName:      "Введите название сервиса:",
			PasswordGenerated:     "Сгенерированный пароль:",
			CopiedToClipboard:     "Пароль скопирован в буфер обмена",
			ClipboardCleared:      "Буфер обмена очищен",
			ClipboardWillClear:    "Автоматическая очистка через",
			ClipboardSecurityInfo: "🛡️ Для безопасности пароль будет удален из буфера обмена",
			ClipboardWaitingInfo:  "⏳ Ожидание очистки... (нажмите Ctrl+C для немедленного выхода)",
			Version:               "PGen CLI v" + version,
			Usage:                 "pgen [флаги]",
			Description:           "Утилита для генерации паролей из мастер-пароля",
			AppTitle:              "🔑 PGen CLI",
			AppSubtitle:           "Генератор безопасных паролей",
			GeneratingPassword:    "🔄 Генерирую пароль...",
			LengthLabel:           "Длина:",
			CharactersLabel:       "символов",
			Tip:                   "💡 Совет: Используйте одинаковые мастер-пароль и название сервиса для получения того же пароля.",
			Tips: []string{
				"💡 Совет: Используйте одинаковые мастер-пароль и название сервиса для получения того же пароля.",
				"🔒 Совет: Используйте сильный мастер-пароль - он защищает все ваши сервисы.",
				"📝 Совет: Запомните свой мастер-пароль надежно - никуда его не записывайте!",
				"⚡ Совет: Используйте описательные имена сервисов: 'gmail.com', 'work-email', 'banking'.",
				"📱 Совет: Можно создавать разные пароли для разных целей: работа, личное, банки.",
				"🔄 Совет: Меняйте мастер-пароль раз в год для максимальной безопасности.",
				"📊 Совет: Увеличьте длину пароля для важных сервисов: --length 24.",
				"📎 Совет: Используйте флаг --copy для автоматического копирования в буфер.",
				"🚫 Совет: Никому не говорите свой мастер-пароль - это ключ ко всем вашим аккаунтам!",
				"🌍 Совет: PGen работает оффлайн - ваши пароли никуда не отправляются!",
				"⚙️ Совет: Настройте конфигурацию: pgen config set default_length 20.",
				"📲 Совет: Можно сохранить список сервисов в заметках для удобства.",
			},

			MasterPasswordStrength:     "Сила мастер-пароля:",
			PasswordStrengthWeak:       "🔴 Очень слабый",
			PasswordStrengthFair:       "🟠 Слабый",
			PasswordStrengthGood:       "🟡 Средний",
			PasswordStrengthStrong:     "🟢 Сильный",
			PasswordStrengthVeryStrong: "🟢 Очень сильный",

			PasswordInfo:     "📊 Информация о пароле:",
			Charset:          "Набор символов:",
			Entropy:          "Энтропия:",
			TimeToCrack:      "Время взлома:",
			Composition:      "Состав:",
			Uppercase:        "заглавные",
			Lowercase:        "строчные",
			Numbers:          "цифры",
			Symbols:          "символы",
			Strength:         "Сила:",
			CrackAssumptions: "При использовании Argon2 на современном оборудовании",

			// Единицы времени для взлома
			TimeInstantly:        "Мгновенно",
			TimeSeconds:          "%.0f секунд",
			TimeMinutes:          "%.0f минут",
			TimeHours:            "%.0f часов",
			TimeDays:             "%.0f дней",
			TimeYears:            "%.0f лет",
			TimeThousandYears:    "%.0f тысяч лет",
			TimeMillionYears:     "%.0f миллионов лет",
			TimeBillionYears:     "%.0f миллиардов лет",
			TimeMoreThanTrillion: "Больше триллиона лет",

			// Уровни силы пароля (без эмодзи для использования в коде)
			StrengthVeryWeak:   "Очень слабый",
			StrengthWeak:       "Слабый",
			StrengthFair:       "Средний",
			StrengthGood:       "Сильный",
			StrengthStrong:     "Очень сильный",
			StrengthVeryStrong: "Очень сильный",

			// Ошибки конфигурации
			ConfigErrorReading:      "ошибка чтения файла конфигурации: %v",
			ConfigErrorParsing:      "ошибка разбора файла конфигурации: %v",
			ConfigErrorEncoding:     "ошибка кодирования конфигурации: %v",
			ConfigErrorWriting:      "ошибка записи файла конфигурации: %v",
			ConfigErrorExportEncode: "ошибка кодирования конфигурации для экспорта: %v",
			ConfigErrorExportWrite:  "ошибка записи файла экспорта: %v",
			ConfigErrorImportRead:   "ошибка чтения файла импорта: %v",
			ConfigErrorImportParse:  "ошибка разбора файла импорта: %v",

			// Команды конфигурации
			ConfigShow:                 "Показать текущую конфигурацию",
			ConfigSet:                  "Установить значение конфигурации",
			ConfigReset:                "Сбросить конфигурацию к стандартным значениям",
			ConfigExport:               "Экспортировать конфигурацию в файл",
			ConfigImport:               "Импортировать конфигурацию из файла",
			ConfigCurrentConfig:        "Текущая конфигурация:",
			ConfigUpdated:              "Конфигурация обновлена:",
			ConfigReset_:               "Конфигурация сброшена к стандартным значениям",
			ConfigExported:             "Конфигурация экспортирована в",
			ConfigImported:             "Конфигурация импортирована из",
			ConfigInvalidKey:           "Неизвестный ключ конфигурации:",
			ConfigManagement:           "Управление конфигурацией",
			ConfigManageSettings:       "Управление настройками PGen",
			ConfigErrorSaving:          "Ошибка сохранения конфигурации:",
			ConfigErrorExporting:       "Ошибка экспорта конфигурации:",
			ConfigErrorImporting:       "Ошибка импорта конфигурации:",
			ConfigErrorSavingImported:  "Ошибка сохранения импортированной конфигурации:",
			ConfigInvalidArgonTime:     "Неверное значение argon_time:",
			ConfigInvalidArgonMemory:   "Неверное значение argon_memory:",
			ConfigInvalidArgonThreads:  "Неверное значение argon_threads:",
			ConfigInvalidArgonKeyLen:   "Неверное значение argon_key_len:",
			ConfigInvalidDefaultLength: "Неверное значение default_length:",
			ConfigInvalidDefaultLang:   "Неверное значение default_language:",
			ConfigInvalidCharset:       "Неверное значение character_set:",
			ConfigInvalidDefaultCopy:   "Неверное значение default_copy:",
			ConfigInvalidClearTimeout:  "Неверное значение default_clear_timeout:",
			ConfigInvalidPasswordInfo:  "Неверное значение show_password_info:",
			ConfigInvalidColorOutput:   "Неверное значение color_output:",
			ConfigInvalidUsername:      "Неверное значение username:",
			ConfigUsernameEmpty:        "Имя пользователя не может быть пустым",
			ConfigUnknownKey:           "Неизвестный ключ конфигурации:",
			ProfileLabel:               "профиль:",
			ConfigLengthRange:          "default_length должен быть от 4 до 128",
			ConfigLanguageValues:       "default_language должен быть 'ru', 'en' или 'auto'",
			ConfigCharsetValues:        "character_set должен быть 'alphanumeric', 'alphanumeric_symbols' или 'symbols_only'",
			ConfigTimeoutRange:         "default_clear_timeout должен быть >= 0",
			About:                      "PGen CLI - безопасный генератор паролей\n\nОписание:\n  Утилита для генерации детерминированных паролей из мастер-пароля\n  с использованием криптографически стойкого алгоритма Argon2.\n\nОсобенности:\n  • Одинаковые входные данные всегда дают одинаковый результат\n  • Высокая криптографическая стойкость (Argon2)\n  • Кроссплатформенность (Windows, Linux, macOS)\n  • Поддержка русского и английского языков\n  • Интеграция с буфером обмена\n\nБезопасность:\n  • Пароли не сохраняются и не передаются по сети\n  • Все вычисления производятся локально\n  • Исходный код открыт для аудита\n\nАвтор: Макс Лейбер ©2025\nEmail: max@leiber.pro\nTelegram: @leiberpro\nЛицензия: MIT",

			// Установка приложения
			InstallSuccess:          "✅ PGen успешно установлен в систему",
			InstallError:            "❌ Ошибка установки PGen",
			InstallAlreadyExists:    "ℹ️ PGen уже установлен в системе",
			InstallPermissionDenied: "🔒 Недостаточно прав для установки. Запустите от имени администратора",
			InstallPathAdded:        "📝 Путь добавлен в переменную PATH",
			InstallComplete:         "🎉 Установка завершена! Перезапустите терминал для применения изменений",
			InstallCheckingPath:     "🔍 Проверяю существующую установку...",
			InstallLocation:         "📍 Расположение:",
			InstallInstalledTo:      "📍 Установлено в:",
			InstallRestart:          "🔄 Для применения изменений перезапустите терминал или выполните: source ~/.bashrc",

			// Удаление приложения
			UninstallSuccess:          "✅ PGen успешно удален из системы",
			UninstallError:            "❌ Ошибка удаления PGen",
			UninstallNotInstalled:     "ℹ️ PGen не установлен в системе",
			UninstallPermissionDenied: "🔒 Недостаточно прав для удаления. Запустите от имени администратора",
			UninstallPathRemoved:      "📝 Путь удален из переменной PATH",
			UninstallComplete:         "🎉 Удаление завершено!",
			UninstallCheckingPath:     "🔍 Проверяю существующую установку...",
			UninstallRemoving:         "🗑️ Удаляю файлы...",
			UninstallConfirm:          "Вы уверены, что хотите удалить PGen? (y/N):",
			UninstallCancelled:        "❌ Удаление отменено",
			UninstallRemovedFrom:      "📍 Удалено из:",

			// Технические ошибки установки
			InstallErrorCreateDir:   "Не удалось создать папку установки",
			InstallErrorGetExePath:  "Не удалось получить путь к исполняемому файлу",
			InstallErrorCopyFile:    "Не удалось скопировать исполняемый файл",
			InstallErrorSetPerms:    "Не удалось установить права на выполнение",
			InstallErrorAddPath:     "Не удалось добавить в PATH",
			InstallErrorOpenFile:    "Не удалось открыть файл",
			InstallErrorWriteFile:   "Не удалось записать в файл",
			InstallErrorResolvePath: "Не удалось разрешить путь к исполняемому файлу",

			// Служебные сообщения установщика
			InstallProfileComment:   "Добавлено установщиком PGen",
			InstallPanicWindowsFunc: "newWindowsInstaller не должно вызываться на Unix",
			InstallPanicUnixFunc:    "newUnixInstaller не должно вызываться на Windows",
			Examples: `Примеры:
  pgen                         # Интерактивный режим
  pgen --copy                  # Скопировать пароль в буфер
  pgen -c -t 30                # Скопировать с очисткой через 30 сек
  pgen --length 20             # Установить длину пароля
  pgen --lang en               # Использовать английский язык
  pgen --install               # Установить PGen в системные пути`,
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
				Lang:             "язык",
				LangDesc:         "Язык интерфейса (ru, en)",
				Length:           "длина",
				LengthDesc:       "Длина генерируемого пароля",
				Copy:             "копировать",
				CopyDesc:         "Скопировать пароль в буфер обмена",
				ClearTimeout:     "время-очистки",
				ClearTimeoutDesc: "Время до автоочистки буфера (секунды, 0=отключить)",
				Version:          "версия",
				VersionDesc:      "Показать версию программы",
				About:            "о-программе",
				AboutDesc:        "Показать информацию о программе",
				Info:             "информация",
				InfoDesc:         "Показать подробную информацию о пароле",
				Install:          "установить",
				InstallDesc:      "Установить PGen в системные пути (PATH)",
				Uninstall:        "удалить",
				UninstallDesc:    "Удалить PGen из системы",
			},
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

				// Проблемы с мастер-паролем
				PasswordIssues struct {
					LengthTooShort  string
					NoLowercase     string
					NoUppercase     string
					NoNumbers       string
					RepeatingChars  string
					SequentialChars string
					CommonWords     string
				}

				// Рекомендации по улучшению
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
				ClipboardError:  "Ошибка копирования в буфер обмена",
				GenerationError: "Ошибка генерации пароля",
				EmptyMaster:     "Мастер-пароль не может быть пустым",
				EmptyService:    "Название сервиса не может быть пустым",
				UserCanceled:    "Операция прервана пользователем",
				InputCanceled:   "Ввод отменен",
				LengthTooShort:  "Минимальная длина пароля: 4 символа",
				LengthTooLong:   "Максимальная длина пароля: 128 символов",
				HashTooShort:    "Не удается сгенерировать пароль требуемой длины",

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
					SequentialChars: "Содержит последовательности символов",
					CommonWords:     "Содержит словарные слова",
				},

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
					IncreaseLength:  "Увеличьте длину пароля (минимум 12 символов)",
					AddLowercase:    "Добавьте строчные буквы (a-z)",
					AddUppercase:    "Добавьте заглавные буквы (A-Z)",
					AddNumbers:      "Добавьте цифры (0-9)",
					AddSymbols:      "Добавьте специальные символы (!@#$%)",
					AvoidRepetition: "Избегайте повторения символов подряд",
					AvoidSequences:  "Избегайте последовательностей (abc, 123)",
					AvoidDictionary: "Избегайте словарных слов",
				},
			},
		}
	default: // English
		return &Messages{
			EnterMasterPassword:   "Enter master password:",
			EnterServiceName:      "Enter service name:",
			PasswordGenerated:     "Generated password:",
			CopiedToClipboard:     "Password copied to clipboard",
			ClipboardCleared:      "Clipboard cleared",
			ClipboardWillClear:    "Auto-clear in",
			ClipboardSecurityInfo: "🛡️ For security, password will be removed from clipboard",
			ClipboardWaitingInfo:  "⏳ Waiting for cleanup... (press Ctrl+C to exit immediately)",
			Version:               "PGen CLI v" + version,
			Usage:                 "pgen [flags]",
			Description:           "Utility for generating passwords from master password",
			AppTitle:              "🔑 PGen CLI",
			AppSubtitle:           "Secure Password Generator",
			GeneratingPassword:    "🔄 Generating password...",
			LengthLabel:           "Length:",
			CharactersLabel:       "characters",
			Tip:                   "💡 Tip: Use the same master password and service name to get the same password.",
			Tips: []string{
				"💡 Tip: Use the same master password and service name to get the same password.",
				"🔒 Tip: Use a strong master password - it protects all your services.",
				"📝 Tip: Remember your master password securely - never write it down!",
				"⚡ Tip: Use descriptive service names: 'gmail.com', 'work-email', 'banking'.",
				"📱 Tip: Create different passwords for different purposes: work, personal, banking.",
				"🔄 Tip: Change your master password once a year for maximum security.",
				"📊 Tip: Increase password length for important services: --length 24.",
				"📎 Tip: Use --copy flag for automatic clipboard copying.",
				"🚫 Tip: Never share your master password - it's the key to all your accounts!",
				"🌍 Tip: PGen works offline - your passwords never leave your device!",
				"⚙️ Tip: Configure settings: pgen config set default_length 20.",
				"📲 Tip: Keep a list of your services in notes for convenience.",
			},

			MasterPasswordStrength:     "Master password strength:",
			PasswordStrengthWeak:       "🔴 Very weak",
			PasswordStrengthFair:       "🟠 Weak",
			PasswordStrengthGood:       "🟡 Fair",
			PasswordStrengthStrong:     "🟢 Strong",
			PasswordStrengthVeryStrong: "🟢 Very strong",

			PasswordInfo:     "📊 Password information:",
			Charset:          "Character set:",
			Entropy:          "Entropy:",
			TimeToCrack:      "Time to crack:",
			Composition:      "Composition:",
			Uppercase:        "uppercase",
			Lowercase:        "lowercase",
			Numbers:          "numbers",
			Symbols:          "symbols",
			Strength:         "Strength:",
			CrackAssumptions: "Using Argon2 on modern hardware",

			// Единицы времени для взлома
			TimeInstantly:        "Instantly",
			TimeSeconds:          "%.0f seconds",
			TimeMinutes:          "%.0f minutes",
			TimeHours:            "%.0f hours",
			TimeDays:             "%.0f days",
			TimeYears:            "%.0f years",
			TimeThousandYears:    "%.0f thousand years",
			TimeMillionYears:     "%.0f million years",
			TimeBillionYears:     "%.0f billion years",
			TimeMoreThanTrillion: "More than trillion years",

			// Уровни силы пароля (без эмодзи для использования в коде)
			StrengthVeryWeak:   "Very weak",
			StrengthWeak:       "Weak",
			StrengthFair:       "Fair",
			StrengthGood:       "Good",
			StrengthStrong:     "Strong",
			StrengthVeryStrong: "Very strong",

			// Configuration error messages
			ConfigErrorReading:      "error reading config file: %v",
			ConfigErrorParsing:      "error parsing config file: %v",
			ConfigErrorEncoding:     "error encoding config: %v",
			ConfigErrorWriting:      "error writing config file: %v",
			ConfigErrorExportEncode: "error encoding config for export: %v",
			ConfigErrorExportWrite:  "error writing export file: %v",
			ConfigErrorImportRead:   "error reading import file: %v",
			ConfigErrorImportParse:  "error parsing import file: %v",

			// Configuration commands
			ConfigShow:                 "Show current configuration",
			ConfigSet:                  "Set configuration value",
			ConfigReset:                "Reset configuration to defaults",
			ConfigExport:               "Export configuration to file",
			ConfigImport:               "Import configuration from file",
			ConfigCurrentConfig:        "Current Configuration:",
			ConfigUpdated:              "Configuration updated:",
			ConfigReset_:               "Configuration reset to defaults",
			ConfigExported:             "Configuration exported to",
			ConfigImported:             "Configuration imported from",
			ConfigInvalidKey:           "Unknown configuration key:",
			ConfigManagement:           "Configuration management",
			ConfigManageSettings:       "Manage PGen configuration settings",
			ConfigErrorSaving:          "Error saving config:",
			ConfigErrorExporting:       "Error exporting config:",
			ConfigErrorImporting:       "Error importing config:",
			ConfigErrorSavingImported:  "Error saving imported config:",
			ConfigInvalidArgonTime:     "Invalid argon_time value:",
			ConfigInvalidArgonMemory:   "Invalid argon_memory value:",
			ConfigInvalidArgonThreads:  "Invalid argon_threads value:",
			ConfigInvalidArgonKeyLen:   "Invalid argon_key_len value:",
			ConfigInvalidDefaultLength: "Invalid default_length value:",
			ConfigInvalidDefaultLang:   "Invalid default_language value:",
			ConfigInvalidCharset:       "Invalid character_set value:",
			ConfigInvalidDefaultCopy:   "Invalid default_copy value:",
			ConfigInvalidClearTimeout:  "Invalid default_clear_timeout value:",
			ConfigInvalidPasswordInfo:  "Invalid show_password_info value:",
			ConfigInvalidColorOutput:   "Invalid color_output value:",
			ConfigInvalidUsername:      "Invalid username value:",
			ConfigUsernameEmpty:        "Username cannot be empty",
			ConfigUnknownKey:           "Unknown configuration key:",
			ProfileLabel:               "profile:",
			ConfigLengthRange:          "default_length must be between 4 and 128",
			ConfigLanguageValues:       "default_language must be 'ru', 'en' or 'auto'",
			ConfigCharsetValues:        "character_set must be 'alphanumeric', 'alphanumeric_symbols', or 'symbols_only'",
			ConfigTimeoutRange:         "default_clear_timeout must be >= 0",
			About:                      "PGen CLI - Secure Password Generator\n\nDescription:\n  A utility for generating deterministic passwords from a master password\n  using the cryptographically strong Argon2 algorithm.\n\nFeatures:\n  • Same input always produces the same result\n  • High cryptographic strength (Argon2)\n  • Cross-platform support (Windows, Linux, macOS)\n  • Russian and English language support\n  • Clipboard integration\n\nSecurity:\n  • Passwords are not stored or transmitted over the network\n  • All computations are performed locally\n  • Source code is open for audit\n\nAuthor: Max Leiber ©2025\nEmail: max@leiber.pro\nTelegram: @leiberpro\nLicense: MIT",

			// Installation messages
			InstallSuccess:          "✅ PGen successfully installed to system",
			InstallError:            "❌ Error installing PGen",
			InstallAlreadyExists:    "ℹ️ PGen is already installed in the system",
			InstallPermissionDenied: "🔒 Insufficient permissions for installation. Run as administrator",
			InstallPathAdded:        "📝 Path added to PATH environment variable",
			InstallComplete:         "🎉 Installation completed! Restart your terminal to apply changes",
			InstallCheckingPath:     "🔍 Checking existing installation...",
			InstallLocation:         "📍 Location:",
			InstallInstalledTo:      "📍 Installed to:",
			InstallRestart:          "🔄 To apply changes, restart terminal or run: source ~/.bashrc",

			// Uninstallation messages
			UninstallSuccess:          "✅ PGen successfully uninstalled from system",
			UninstallError:            "❌ Error uninstalling PGen",
			UninstallNotInstalled:     "ℹ️ PGen is not installed in the system",
			UninstallPermissionDenied: "🔒 Insufficient permissions for uninstallation. Run as administrator",
			UninstallPathRemoved:      "📝 Path removed from PATH environment variable",
			UninstallComplete:         "🎉 Uninstallation completed!",
			UninstallCheckingPath:     "🔍 Checking existing installation...",
			UninstallRemoving:         "🗑️ Removing files...",
			UninstallConfirm:          "Are you sure you want to uninstall PGen? (y/N):",
			UninstallCancelled:        "❌ Uninstallation cancelled",
			UninstallRemovedFrom:      "📍 Removed from:",

			// Technical installation errors
			InstallErrorCreateDir:   "Failed to create install directory",
			InstallErrorGetExePath:  "Failed to get executable path",
			InstallErrorCopyFile:    "Failed to copy executable",
			InstallErrorSetPerms:    "Failed to set executable permissions",
			InstallErrorAddPath:     "Failed to add to PATH",
			InstallErrorOpenFile:    "Failed to open file",
			InstallErrorWriteFile:   "Failed to write to file",
			InstallErrorResolvePath: "Failed to resolve executable path",

			// Installer service messages
			InstallProfileComment:   "Added by PGen installer",
			InstallPanicWindowsFunc: "newWindowsInstaller should not be called on Unix",
			InstallPanicUnixFunc:    "newUnixInstaller should not be called on Windows",
			Examples: `Examples:
  pgen                         # Interactive mode
  pgen --copy                  # Copy password to clipboard
  pgen -c -t 30                # Copy with 30sec auto-clear
  pgen --length 20             # Set password length
  pgen --lang ru               # Use Russian language
  pgen --install               # Install PGen to system PATH`,
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
				Lang:             "lang",
				LangDesc:         "Interface language (ru, en)",
				Length:           "length",
				LengthDesc:       "Length of generated password",
				Copy:             "copy",
				CopyDesc:         "Copy password to clipboard",
				ClearTimeout:     "clear-timeout",
				ClearTimeoutDesc: "Time before auto-clearing clipboard (seconds, 0=disable)",
				Version:          "version",
				VersionDesc:      "Show program version",
				About:            "about",
				AboutDesc:        "Show information about the program",
				Info:             "info",
				InfoDesc:         "Show detailed password information",
				Install:          "install",
				InstallDesc:      "Install PGen to system PATH",
				Uninstall:        "uninstall",
				UninstallDesc:    "Uninstall PGen from system",
			},
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

				// Проблемы с мастер-паролем
				PasswordIssues struct {
					LengthTooShort  string
					NoLowercase     string
					NoUppercase     string
					NoNumbers       string
					RepeatingChars  string
					SequentialChars string
					CommonWords     string
				}

				// Рекомендации по улучшению
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
				ClipboardError:  "Error copying to clipboard",
				GenerationError: "Error generating password",
				EmptyMaster:     "Master password cannot be empty",
				EmptyService:    "Service name cannot be empty",
				UserCanceled:    "Operation canceled by user",
				InputCanceled:   "Input canceled",
				LengthTooShort:  "Minimum password length: 4 characters",
				LengthTooLong:   "Maximum password length: 128 characters",
				HashTooShort:    "Cannot generate password of requested length",

				PasswordIssues: struct {
					LengthTooShort  string
					NoLowercase     string
					NoUppercase     string
					NoNumbers       string
					RepeatingChars  string
					SequentialChars string
					CommonWords     string
				}{
					LengthTooShort:  "Password is too short",
					NoLowercase:     "Missing lowercase letters",
					NoUppercase:     "Missing uppercase letters",
					NoNumbers:       "Missing numbers",
					RepeatingChars:  "Contains repeating characters",
					SequentialChars: "Contains sequential characters",
					CommonWords:     "Contains dictionary words",
				},

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
					IncreaseLength:  "Increase password length (minimum 12 characters)",
					AddLowercase:    "Add lowercase letters (a-z)",
					AddUppercase:    "Add uppercase letters (A-Z)",
					AddNumbers:      "Add numbers (0-9)",
					AddSymbols:      "Add special symbols (!@#$%)",
					AvoidRepetition: "Avoid repeating characters",
					AvoidSequences:  "Avoid sequences (abc, 123)",
					AvoidDictionary: "Avoid dictionary words",
				},
			},
		}
	}
}

func DetectLanguage(langFlag string) Language {
	if langFlag != "" {
		switch strings.ToLower(langFlag) {
		case "ru", "russian":
			return Russian
		case "en", "english":
			return English
		}
	}

	locale := os.Getenv("LANG")
	if locale == "" {
		locale = os.Getenv("LC_ALL")
	}
	if locale == "" {
		locale = os.Getenv("LC_MESSAGES")
	}

	if strings.Contains(strings.ToLower(locale), "ru") {
		return Russian
	}

	return English
}
