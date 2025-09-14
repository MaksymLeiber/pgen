package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/MaksymLeiber/pgen/internal/analyzer"
	"github.com/MaksymLeiber/pgen/internal/clipboard"
	"github.com/MaksymLeiber/pgen/internal/colors"
	"github.com/MaksymLeiber/pgen/internal/config"
	"github.com/MaksymLeiber/pgen/internal/generator"
	"github.com/MaksymLeiber/pgen/internal/i18n"
	"github.com/MaksymLeiber/pgen/internal/input"
	"github.com/MaksymLeiber/pgen/internal/installer"
	"github.com/MaksymLeiber/pgen/internal/validator"
)

var (
	langFlag      string
	lengthFlag    int
	copyFlag      bool
	clearTimeout  int
	versionFlag   bool
	aboutFlag     bool
	showInfoFlag  bool
	installFlag   bool
	uninstallFlag bool
	Version       string
	cfg           *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "pgen",
	Short: "",
	Long:  "",
	Run:   runRootCommand,
}

func Execute() {
	// Загружаем конфигурацию
	var err error
	// Используем английский язык по умолчанию для загрузки конфигурации
	defaultMessages := i18n.GetMessages(i18n.English, Version)
	cfg, err = config.Load(defaultMessages)
	if err != nil {
		cfg = config.DefaultConfig()
	}

	// Добавляем команды управления конфигурацией
	rootCmd.AddCommand(configCmd)

	lang := detectLanguageFromArgs()
	messages := i18n.GetMessages(lang, Version)
	updateCommandTexts(rootCmd, messages)
	updateConfigCommandTexts(messages)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Персистентные флаги (доступны во всех подкомандах)
	rootCmd.PersistentFlags().StringVarP(&langFlag, "lang", "l", "", "")

	// Локальные флаги (только для корневой команды)
	rootCmd.Flags().IntVarP(&lengthFlag, "length", "n", 16, "")
	rootCmd.Flags().BoolVarP(&copyFlag, "copy", "c", false, "")
	rootCmd.Flags().IntVarP(&clearTimeout, "clear-timeout", "t", 45, "")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "")
	rootCmd.Flags().BoolVarP(&aboutFlag, "about", "a", false, "")
	rootCmd.Flags().BoolVarP(&showInfoFlag, "info", "i", false, "")
	rootCmd.Flags().BoolVarP(&installFlag, "install", "", false, "")
	rootCmd.Flags().BoolVarP(&uninstallFlag, "uninstall", "", false, "")

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Определяем эффективную длину
		length := lengthFlag
		if !cmd.Flags().Changed("length") {
			length = cfg.DefaultLength
		}

		if err := generator.ValidateLength(length); err != nil {
			messages := i18n.GetMessages(detectLanguageFromArgs(), Version)
			switch err.Error() {
			case "length_too_short":
				return errors.New(messages.Errors.LengthTooShort)
			case "length_too_long":
				return errors.New(messages.Errors.LengthTooLong)
			default:
				return err
			}
		}
		return nil
	}
}

func runRootCommand(cmd *cobra.Command, args []string) {
	// Используем уже определенный язык
	language := detectLanguageFromArgs()
	messages := i18n.GetMessages(language, Version)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintf(os.Stderr, "\n%s\n", colors.ErrorMsg(messages.Errors.UserCanceled))
		os.Exit(1)
	}()

	fmt.Println(colors.TitleMsg(messages.AppTitle))
	fmt.Println(colors.SubtleMsg(messages.AppSubtitle + "\n"))

	if versionFlag {
		fmt.Println(colors.InfoMsg(messages.Version))
		return
	}

	if aboutFlag {
		fmt.Println(colors.InfoMsg(messages.About))
		return
	}

	if installFlag {
		runInstallation(messages)
		return
	}

	if uninstallFlag {
		runUninstallation(messages)
		return
	}

	fmt.Print(colors.PromptMsg(messages.EnterMasterPassword + " "))
	masterPassword, err := input.ReadPasswordWithStarsAndMessages(&input.InputMessages{
		UserCanceled:  messages.Errors.UserCanceled,
		InputCanceled: messages.Errors.InputCanceled,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", colors.ErrorMsg(messages.Errors.GenerationError+":"), err)
		os.Exit(1)
	}

	if len(masterPassword) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", colors.ErrorMsg(messages.Errors.EmptyMaster))
		os.Exit(1)
	}

	// Проверка силы мастер-пароля
	strength := validator.ValidatePasswordStrength(masterPassword, messages)
	displayPasswordStrength(strength, messages)

	fmt.Print(colors.PromptMsg(messages.EnterServiceName + " "))
	serviceName, err := input.ReadLine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", colors.ErrorMsg(messages.Errors.GenerationError+":"), err)
		os.Exit(1)
	}

	if serviceName == "" {
		fmt.Fprintf(os.Stderr, "%s\n", colors.ErrorMsg(messages.Errors.EmptyService))
		os.Exit(1)
	}

	fmt.Print(colors.SubtleMsg(messages.GeneratingPassword + "\n"))

	// Определяем эффективную длину пароля
	length := lengthFlag
	if !cmd.Flags().Changed("length") {
		length = cfg.DefaultLength // Используем из конфигурации
	}

	// Создаем генератор с конфигурацией
	gen := generator.NewPasswordGeneratorWithConfig(length, generator.ArgonConfig{
		Time:    cfg.ArgonTime,
		Memory:  cfg.ArgonMemory,
		Threads: cfg.ArgonThreads,
		KeyLen:  cfg.ArgonKeyLen,
	})
	password, err := gen.GeneratePassword(masterPassword, serviceName, messages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", colors.ErrorMsg(messages.Errors.GenerationError+":"), err)
		os.Exit(1)
	}

	fmt.Printf("\n%s %s\n", colors.InfoMsg(messages.PasswordGenerated), colors.GeneratedMsg(password))
	fmt.Printf("%s %s\n", colors.SubtleMsg(messages.LengthLabel), colors.SubtleMsg(fmt.Sprintf("%d %s", len(password), messages.CharactersLabel)))

	// Показ информации о пароле
	if showInfoFlag || cfg.ShowPasswordInfo {
		displayPasswordInfo(password, messages)
	}

	if copyFlag || cfg.DefaultCopy {
		// Определяем эффективный таймаут
		effectiveTimeout := clearTimeout
		if !cmd.Flags().Changed("clear-timeout") {
			effectiveTimeout = cfg.DefaultClearTimeout
		}

		// Используем настраиваемый таймаут для очистки
		timeoutDuration := time.Duration(effectiveTimeout) * time.Second
		done, err := clipboard.CopyToClipboardWithTimeout(password, timeoutDuration)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %v\n", colors.ErrorMsg(messages.Errors.ClipboardError+":"), err)
		} else {
			fmt.Printf("%s %s\n", colors.SuccessMsg("✓"), colors.SuccessMsg(messages.CopiedToClipboard))
			if effectiveTimeout > 0 && done != nil {
				// Объясняем пользователю, что происходит
				fmt.Printf("%s %s %ds\n", colors.SubtleMsg("⏱"), colors.SubtleMsg(messages.ClipboardWillClear), effectiveTimeout)
				fmt.Printf("%s\n", colors.SubtleMsg(messages.ClipboardSecurityInfo))
				fmt.Printf("%s\n", colors.InfoMsg(messages.ClipboardWaitingInfo))
				// Ожидаем завершения очистки буфера
				<-done
				fmt.Printf("%s %s\n", colors.SuccessMsg("✓"), colors.SuccessMsg(messages.ClipboardCleared))
			}
		}
	}

	fmt.Println(colors.SubtleMsg("\n" + messages.GetRandomTip()))
}

// runInstallation выполняет установку приложения в системные пути
func runInstallation(messages *i18n.Messages) {
	fmt.Println(colors.InfoMsg(messages.InstallCheckingPath))

	// Создаем инсталлятор для текущей платформы
	installer := installer.NewSystemInstaller(messages)

	// Проверяем, установлено ли уже приложение
	if installer.IsInstalled() {
		fmt.Println(colors.InfoMsg(messages.InstallAlreadyExists))
		fmt.Printf("%s %s\n", colors.SubtleMsg(messages.InstallLocation), colors.SubtleMsg(installer.GetInstallPath()))
		return
	}

	// Проверяем права доступа
	if needsElevation() {
		fmt.Println(colors.ErrorMsg(messages.InstallPermissionDenied))
		os.Exit(1)
	}

	fmt.Println(colors.InfoMsg(messages.InstallCopyingFile))

	// Выполняем установку
	if err := installer.Install(messages); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", colors.ErrorMsg(messages.InstallError), err)
		os.Exit(1)
	}

	fmt.Println(colors.SuccessMsg(messages.InstallSuccess))
	fmt.Printf("%s %s\n", colors.SubtleMsg(messages.InstallInstalledTo), colors.SubtleMsg(installer.GetInstallPath()))
	fmt.Println(colors.InfoMsg(messages.InstallPathAdded))
	fmt.Println(colors.InfoMsg(messages.InstallComplete))

	// Показываем инструкции по перезапуску
	if runtime.GOOS != "windows" {
		fmt.Println(colors.SubtleMsg(messages.InstallRestart))
	}
}

// runUninstallation выполняет удаление приложения из системы
func runUninstallation(messages *i18n.Messages) {
	fmt.Println(colors.InfoMsg(messages.UninstallCheckingPath))

	// Создаем инсталлятор для текущей платформы
	installer := installer.NewSystemInstaller(messages)

	// Проверяем, установлено ли приложение
	if !installer.IsInstalled() {
		fmt.Println(colors.InfoMsg(messages.UninstallNotInstalled))
		return
	}

	// Подтверждение удаления
	fmt.Printf("%s ", colors.PromptMsg(messages.UninstallConfirm))
	var confirmation string
	fmt.Scanln(&confirmation)

	if strings.ToLower(confirmation) != "y" && strings.ToLower(confirmation) != "д" {
		fmt.Println(colors.InfoMsg(messages.UninstallCancelled))
		return
	}

	// Проверяем права доступа
	if needsElevation() {
		fmt.Println(colors.ErrorMsg(messages.UninstallPermissionDenied))
		os.Exit(1)
	}

	fmt.Println(colors.InfoMsg(messages.UninstallRemoving))

	// Выполняем удаление
	if err := installer.Uninstall(messages); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", colors.ErrorMsg(messages.UninstallError), err)
		os.Exit(1)
	}

	fmt.Println(colors.SuccessMsg(messages.UninstallSuccess))
	fmt.Printf("%s %s\n", colors.SubtleMsg(messages.UninstallRemovedFrom), colors.SubtleMsg(installer.GetInstallPath()))
	fmt.Println(colors.InfoMsg(messages.UninstallComplete))
}

// needsElevation определяет, нужны ли повышенные права для установки
func needsElevation() bool {
	switch runtime.GOOS {
	case "windows":
		// На Windows проверяем, пытаемся ли мы установить в системную директорию
		programFiles := os.Getenv("PROGRAMFILES")
		if programFiles == "" {
			programFiles = "C:\\Program Files"
		}

		// Получаем текущий путь установки
		currentPath := getCurrentInstallPath()

		// Если путь установки находится в Program Files, нужны права администратора
		if strings.HasPrefix(currentPath, programFiles) {
			return !isWindowsAdmin()
		}

		return false
	default:
		// На Unix системах для /usr/local/bin нужны права root
		return os.Geteuid() != 0
	}
}

// getCurrentInstallPath получает текущий путь установки
func getCurrentInstallPath() string {
	// Создаем временный инсталлятор для получения пути
	messages := i18n.GetMessages(i18n.English, "test")
	installer := installer.NewSystemInstaller(messages)
	return installer.GetInstallPath()
}

// isWindowsAdmin проверяет, запущен ли процесс с правами администратора на Windows
func isWindowsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func updateCommandTexts(cmd *cobra.Command, messages *i18n.Messages) {
	cmd.Use = messages.Usage
	cmd.Short = messages.Description
	cmd.Long = messages.Description + "\n\n" + messages.Examples

	// Настраиваем шаблон для русской локализации
	if strings.Contains(messages.Usage, "флаги") {
		cmd.SetUsageTemplate(`Использование:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [команда]{{end}}{{if gt (len .Aliases) 0}}

Псевдонимы:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Примеры:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Доступные команды:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Флаги:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Глобальные флаги:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Дополнительные команды:
{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Используйте "{{.CommandPath}} [command] --help" для получения информации о команде.{{end}}
`)
	}

	if flag := cmd.Flag("lang"); flag != nil {
		flag.Usage = messages.Flags.LangDesc
	}
	if flag := cmd.Flag("length"); flag != nil {
		flag.Usage = messages.Flags.LengthDesc
	}
	if flag := cmd.Flag("copy"); flag != nil {
		flag.Usage = messages.Flags.CopyDesc
	}
	if flag := cmd.Flag("clear-timeout"); flag != nil {
		flag.Usage = messages.Flags.ClearTimeoutDesc
	}
	if flag := cmd.Flag("version"); flag != nil {
		flag.Usage = messages.Flags.VersionDesc
	}
	if flag := cmd.Flag("about"); flag != nil {
		flag.Usage = messages.Flags.AboutDesc
	}
	if flag := cmd.Flag("info"); flag != nil {
		flag.Usage = messages.Flags.InfoDesc
	}
	if flag := cmd.Flag("install"); flag != nil {
		flag.Usage = messages.Flags.InstallDesc
	}
	if flag := cmd.Flag("uninstall"); flag != nil {
		flag.Usage = messages.Flags.UninstallDesc
	}

}

func detectLanguageFromArgs() i18n.Language {
	for i, arg := range os.Args {
		if (arg == "--lang" || arg == "-l") && i+1 < len(os.Args) {
			return i18n.DetectLanguage(os.Args[i+1])
		}
		if len(arg) > 7 && arg[:7] == "--lang=" {
			return i18n.DetectLanguage(arg[7:])
		}
	}
	// Если флаг не указан, используем значение из конфигурации
	if cfg != nil && cfg.DefaultLanguage != "auto" {
		return i18n.DetectLanguage(cfg.DefaultLanguage)
	}
	return i18n.DetectLanguage("")
}

// displayPasswordStrength отображает информацию о силе мастер-пароля
func displayPasswordStrength(strength *validator.PasswordStrength, messages *i18n.Messages) {
	fmt.Printf("%s ", colors.SubtleMsg(messages.MasterPasswordStrength))

	// Отображаем уровень силы
	switch strength.Level {
	case validator.StrengthWeak:
		fmt.Println(colors.ErrorMsg(messages.PasswordStrengthWeak))
	case validator.StrengthFair:
		fmt.Println(colors.InfoMsg(messages.PasswordStrengthFair))
	case validator.StrengthGood:
		fmt.Println(colors.InfoMsg(messages.PasswordStrengthGood))
	case validator.StrengthStrong:
		fmt.Println(colors.SuccessMsg(messages.PasswordStrengthStrong))
	case validator.StrengthVeryStrong:
		fmt.Println(colors.SuccessMsg(messages.PasswordStrengthVeryStrong))
	}

	// Показываем проблемы и рекомендации
	if len(strength.Issues) > 0 || len(strength.Suggestions) > 0 {
		for _, issue := range strength.Issues {
			issueText := getIssueText(issue, messages)
			fmt.Printf("%s %s\n", colors.ErrorMsg("⚠️"), colors.SubtleMsg(issueText))
		}

		for _, suggestion := range strength.Suggestions {
			suggestionText := getSuggestionText(suggestion, messages)
			fmt.Printf("%s %s\n", colors.InfoMsg("💡"), colors.SubtleMsg(suggestionText))
		}
		fmt.Println()
	}
}

// displayPasswordInfo отображает информацию о сгенерированном пароле
func displayPasswordInfo(password string, messages *i18n.Messages) {
	info := analyzer.AnalyzePassword(password, messages)

	fmt.Printf("\n%s\n", colors.InfoMsg(messages.PasswordInfo))
	fmt.Printf("%s %s (%d символов)\n", colors.SubtleMsg(messages.Charset), colors.SubtleMsg(info.Charset), info.CharsetSize)
	fmt.Printf("%s %.1f бит\n", colors.SubtleMsg(messages.Entropy), info.Entropy)
	fmt.Printf("%s %s\n", colors.SubtleMsg(messages.TimeToCrack), colors.SubtleMsg(info.TimeToCrack.HumanTime))
	fmt.Printf("%s %d %s, %d %s, %d %s, %d %s\n",
		colors.SubtleMsg(messages.Composition),
		info.Composition.Uppercase, messages.Uppercase,
		info.Composition.Lowercase, messages.Lowercase,
		info.Composition.Numbers, messages.Numbers,
		info.Composition.Symbols, messages.Symbols)
	fmt.Printf("%s %s\n", colors.SubtleMsg(messages.Strength), colors.SubtleMsg(info.Strength))
	fmt.Printf("%s %s\n", colors.SubtleMsg("📝"), colors.SubtleMsg(messages.CrackAssumptions))
}

// getIssueText возвращает текст проблемы на соответствующем языке
func getIssueText(issue string, messages *i18n.Messages) string {
	switch issue {
	case "length_too_short":
		return messages.Errors.PasswordIssues.LengthTooShort
	case "no_lowercase":
		return messages.Errors.PasswordIssues.NoLowercase
	case "no_uppercase":
		return messages.Errors.PasswordIssues.NoUppercase
	case "no_numbers":
		return messages.Errors.PasswordIssues.NoNumbers
	case "repeating_chars":
		return messages.Errors.PasswordIssues.RepeatingChars
	case "sequential_chars":
		return messages.Errors.PasswordIssues.SequentialChars
	case "common_words":
		return messages.Errors.PasswordIssues.CommonWords
	default:
		return issue
	}
}

// getSuggestionText возвращает текст рекомендации на соответствующем языке
func getSuggestionText(suggestion string, messages *i18n.Messages) string {
	switch suggestion {
	case "increase_length":
		return messages.Errors.Suggestions.IncreaseLength
	case "add_lowercase":
		return messages.Errors.Suggestions.AddLowercase
	case "add_uppercase":
		return messages.Errors.Suggestions.AddUppercase
	case "add_numbers":
		return messages.Errors.Suggestions.AddNumbers
	case "add_symbols":
		return messages.Errors.Suggestions.AddSymbols
	case "avoid_repetition":
		return messages.Errors.Suggestions.AvoidRepetition
	case "avoid_sequences":
		return messages.Errors.Suggestions.AvoidSequences
	case "avoid_dictionary":
		return messages.Errors.Suggestions.AvoidDictionary
	default:
		return suggestion
	}
}

// Команды для управления конфигурацией
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "",
	Long:  "",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		lang := detectLanguageFromArgs()
		messages := i18n.GetMessages(lang, Version)
		data, _ := json.MarshalIndent(cfg, "", "  ")
		fmt.Printf("%s %s\n%s\n", colors.InfoMsg("📋"), messages.ConfigCurrentConfig, string(data))
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		lang := detectLanguageFromArgs()
		messages := i18n.GetMessages(lang, Version)
		key, value := args[0], args[1]
		if err := setConfigValue(key, value, messages); err != nil {
			fmt.Fprintf(os.Stderr, "%s %s %v\n", colors.ErrorMsg("❌"), messages.ConfigInvalidKey, err)
			os.Exit(1)
		}
		if err := cfg.Save(messages); err != nil {
			fmt.Fprintf(os.Stderr, "%s %s %v\n", colors.ErrorMsg("❌"), messages.ConfigErrorSaving, err)
			os.Exit(1)
		}
		fmt.Printf("%s %s %s = %s\n", colors.SuccessMsg("✓"), messages.ConfigUpdated, key, value)
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		lang := detectLanguageFromArgs()
		messages := i18n.GetMessages(lang, Version)
		cfg = config.DefaultConfig()
		if err := cfg.Save(messages); err != nil {
			fmt.Fprintf(os.Stderr, "%s %s %v\n", colors.ErrorMsg("❌"), messages.ConfigErrorSaving, err)
			os.Exit(1)
		}
		fmt.Printf("%s %s\n", colors.SuccessMsg("✓"), messages.ConfigReset_)
	},
}

var configExportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lang := detectLanguageFromArgs()
		messages := i18n.GetMessages(lang, Version)
		filename := args[0]
		if err := cfg.Export(filename, messages); err != nil {
			fmt.Fprintf(os.Stderr, "%s %s %v\n", colors.ErrorMsg("❌"), messages.ConfigErrorExporting, err)
			os.Exit(1)
		}
		fmt.Printf("%s %s %s\n", colors.SuccessMsg("✓"), messages.ConfigExported, filename)
	},
}

var configImportCmd = &cobra.Command{
	Use:   "import [filename]",
	Short: "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lang := detectLanguageFromArgs()
		messages := i18n.GetMessages(lang, Version)
		filename := args[0]
		importedCfg, err := config.Import(filename, messages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s %v\n", colors.ErrorMsg("❌"), messages.ConfigErrorImporting, err)
			os.Exit(1)
		}
		cfg = importedCfg
		if err := cfg.Save(messages); err != nil {
			fmt.Fprintf(os.Stderr, "%s %s %v\n", colors.ErrorMsg("❌"), messages.ConfigErrorSavingImported, err)
			os.Exit(1)
		}
		fmt.Printf("%s %s %s\n", colors.SuccessMsg("✓"), messages.ConfigImported, filename)
	},
}

// setConfigValue устанавливает значение конфигурации
func setConfigValue(key, value string, messages *i18n.Messages) error {
	switch key {
	case "argon_time":
		val, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidArgonTime, err)
		}
		cfg.ArgonTime = uint32(val)
	case "argon_memory":
		val, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidArgonMemory, err)
		}
		cfg.ArgonMemory = uint32(val)
	case "argon_threads":
		val, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidArgonThreads, err)
		}
		cfg.ArgonThreads = uint8(val)
	case "argon_key_len":
		val, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidArgonKeyLen, err)
		}
		cfg.ArgonKeyLen = uint32(val)
	case "default_length":
		val, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidDefaultLength, err)
		}
		if val < 4 || val > 128 {
			return fmt.Errorf("%s", messages.ConfigLengthRange)
		}
		cfg.DefaultLength = val
	case "default_language":
		if value != "ru" && value != "en" && value != "auto" {
			return fmt.Errorf("%s", messages.ConfigLanguageValues)
		}
		cfg.DefaultLanguage = value
	case "character_set":
		if value != "alphanumeric" && value != "alphanumeric_symbols" && value != "symbols_only" {
			return fmt.Errorf("%s", messages.ConfigCharsetValues)
		}
		cfg.CharacterSet = value
	case "default_copy":
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidDefaultCopy, err)
		}
		cfg.DefaultCopy = val
	case "default_clear_timeout":
		val, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidClearTimeout, err)
		}
		if val < 0 {
			return fmt.Errorf("%s", messages.ConfigTimeoutRange)
		}
		cfg.DefaultClearTimeout = val
	case "show_password_info":
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidPasswordInfo, err)
		}
		cfg.ShowPasswordInfo = val
	case "color_output":
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("%s %v", messages.ConfigInvalidColorOutput, err)
		}
		cfg.ColorOutput = val
	default:
		return fmt.Errorf("%s %s", messages.ConfigUnknownKey, key)
	}
	return nil
}

func init() {
	// Добавляем подкоманды к config
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configExportCmd)
	configCmd.AddCommand(configImportCmd)
}

// updateConfigCommandTexts обновляет тексты команд конфигурации
func updateConfigCommandTexts(messages *i18n.Messages) {
	configCmd.Short = messages.ConfigManagement
	configCmd.Long = messages.ConfigManageSettings

	configShowCmd.Short = messages.ConfigShow
	configSetCmd.Short = messages.ConfigSet
	configResetCmd.Short = messages.ConfigReset
	configExportCmd.Short = messages.ConfigExport
	configImportCmd.Short = messages.ConfigImport
}
