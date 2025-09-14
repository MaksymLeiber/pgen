package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/MaksymLeiber/pgen/internal/i18n"
)

// Config структура конфигурации
type Config struct {
	// Argon2 параметры
	ArgonTime    uint32 `json:"argon_time"`
	ArgonMemory  uint32 `json:"argon_memory"`
	ArgonThreads uint8  `json:"argon_threads"`
	ArgonKeyLen  uint32 `json:"argon_key_len"`

	// Настройки генерации
	DefaultLength   int    `json:"default_length"`
	DefaultLanguage string `json:"default_language"`
	CharacterSet    string `json:"character_set"`

	// Настройки буфера обмена
	DefaultCopy         bool `json:"default_copy"`
	DefaultClearTimeout int  `json:"default_clear_timeout"`

	// Настройки отображения
	ShowPasswordInfo bool `json:"show_password_info"`
	ColorOutput      bool `json:"color_output"`

	// Версия конфигурации для совместимости
	Version string `json:"config_version"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		ArgonTime:           3,
		ArgonMemory:         256 * 1024,
		ArgonThreads:        4,
		ArgonKeyLen:         32,
		DefaultLength:       16,
		DefaultLanguage:     "auto",
		CharacterSet:        "alphanumeric_symbols",
		DefaultCopy:         false,
		DefaultClearTimeout: 45,
		ShowPasswordInfo:    false,
		ColorOutput:         true,
		Version:             "1.0",
	}
}

// GetConfigPath возвращает путь к файлу конфигурации
func GetConfigPath() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(homeDir, "AppData", "Roaming")
		}
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, "Library", "Application Support")
	default: // linux и другие unix
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(homeDir, ".config")
		}
	}

	pgenConfigDir := filepath.Join(configDir, "pgen")
	if err := os.MkdirAll(pgenConfigDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(pgenConfigDir, "config.json"), nil
}

// Load загружает конфигурацию из файла
func Load(messages *i18n.Messages) (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		// Файл не существует, возвращаем конфигурацию по умолчанию
		return DefaultConfig(), nil
	}
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigErrorReading, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf(messages.ConfigErrorParsing, err)
	}

	// Валидация и установка значений по умолчанию
	config.validate()

	return &config, nil
}

// Save сохраняет конфигурацию в файл
func (c *Config) Save(messages *i18n.Messages) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf(messages.ConfigErrorEncoding, err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf(messages.ConfigErrorWriting, err)
	}

	return nil
}

// validate проверяет и исправляет значения конфигурации
func (c *Config) validate() {
	if c.ArgonTime == 0 {
		c.ArgonTime = 3
	}
	if c.ArgonMemory == 0 {
		c.ArgonMemory = 256 * 1024
	}
	if c.ArgonThreads == 0 {
		c.ArgonThreads = 4
	}
	if c.ArgonKeyLen == 0 {
		c.ArgonKeyLen = 32
	}
	if c.DefaultLength < 4 {
		c.DefaultLength = 16
	}
	if c.DefaultLength > 128 {
		c.DefaultLength = 128
	}
	if c.DefaultLanguage == "" {
		c.DefaultLanguage = "auto"
	}
	if c.CharacterSet == "" {
		c.CharacterSet = "alphanumeric_symbols"
	}
	if c.DefaultClearTimeout < 0 {
		c.DefaultClearTimeout = 45
	}
	if c.Version == "" {
		c.Version = "1.0"
	}
}

// Export экспортирует конфигурацию в указанный файл
func (c *Config) Export(filename string, messages *i18n.Messages) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf(messages.ConfigErrorExportEncode, err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf(messages.ConfigErrorExportWrite, err)
	}

	return nil
}

// Import импортирует конфигурацию из указанного файла
func Import(filename string, messages *i18n.Messages) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigErrorImportRead, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf(messages.ConfigErrorImportParse, err)
	}

	config.validate()
	return &config, nil
}
