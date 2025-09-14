package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"pgen/internal/i18n"
)

func TestDefaultConfig(t *testing.T) {
	// Тест создания конфигурации по умолчанию
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() не должен возвращать nil")
	}

	// Проверяем значения по умолчанию
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"ArgonTime", config.ArgonTime, uint32(3)},
		{"ArgonMemory", config.ArgonMemory, uint32(256 * 1024)},
		{"ArgonThreads", config.ArgonThreads, uint8(4)},
		{"ArgonKeyLen", config.ArgonKeyLen, uint32(32)},
		{"DefaultLength", config.DefaultLength, 16},
		{"DefaultLanguage", config.DefaultLanguage, "auto"},
		{"CharacterSet", config.CharacterSet, "alphanumeric_symbols"},
		{"DefaultCopy", config.DefaultCopy, false},
		{"DefaultClearTimeout", config.DefaultClearTimeout, 45},
		{"ShowPasswordInfo", config.ShowPasswordInfo, false},
		{"ColorOutput", config.ColorOutput, true},
		{"Version", config.Version, "1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("DefaultConfig().%s = %v, ожидается %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestGetConfigPath(t *testing.T) {
	// Тест получения пути к конфигурационному файлу
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() ошибка = %v", err)
	}

	if configPath == "" {
		t.Error("GetConfigPath() не должен возвращать пустой путь")
	}

	// Проверяем, что путь заканчивается на config.json
	if !strings.HasSuffix(configPath, "config.json") {
		t.Errorf("GetConfigPath() = %q, должен заканчиваться на 'config.json'", configPath)
	}

	// Проверяем, что путь содержит pgen
	if !strings.Contains(configPath, "pgen") {
		t.Errorf("GetConfigPath() = %q, должен содержать 'pgen'", configPath)
	}

	// Проверяем, что директория создана
	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("GetConfigPath() должен создать директорию %q", dir)
	}
}

func TestGetConfigPathPlatformSpecific(t *testing.T) {
	// Тест платформо-специфичных путей
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() ошибка = %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(configPath, "AppData") && !strings.Contains(configPath, "Roaming") {
			t.Errorf("Windows: путь должен содержать AppData или Roaming, получен: %q", configPath)
		}
	case "darwin":
		if !strings.Contains(configPath, "Library/Application Support") {
			t.Errorf("macOS: путь должен содержать 'Library/Application Support', получен: %q", configPath)
		}
	default:
		if !strings.Contains(configPath, ".config") && !strings.Contains(configPath, "XDG_CONFIG_HOME") {
			t.Errorf("Linux: путь должен содержать '.config', получен: %q", configPath)
		}
	}
}

func TestConfigValidate(t *testing.T) {
	// Тест валидации конфигурации
	tests := []struct {
		name     string
		config   *Config
		expected *Config
	}{
		{
			name: "Нулевые значения Argon2",
			config: &Config{
				ArgonTime:    0,
				ArgonMemory:  0,
				ArgonThreads: 0,
				ArgonKeyLen:  0,
			},
			expected: &Config{
				ArgonTime:           3,
				ArgonMemory:         256 * 1024,
				ArgonThreads:        4,
				ArgonKeyLen:         32,
				DefaultLength:       16,
				DefaultLanguage:     "auto",
				CharacterSet:        "alphanumeric_symbols",
				DefaultClearTimeout: 0, // validate не изменяет 0, только отрицательные значения
				Version:             "1.0",
			},
		},
		{
			name: "Слишком короткий пароль",
			config: &Config{
				DefaultLength: 2,
			},
			expected: &Config{
				ArgonTime:           3,
				ArgonMemory:         256 * 1024,
				ArgonThreads:        4,
				ArgonKeyLen:         32,
				DefaultLength:       16,
				DefaultLanguage:     "auto",
				CharacterSet:        "alphanumeric_symbols",
				DefaultClearTimeout: 0, // validate не изменяет 0, только отрицательные значения
				Version:             "1.0",
			},
		},
		{
			name: "Слишком длинный пароль",
			config: &Config{
				DefaultLength: 200,
			},
			expected: &Config{
				ArgonTime:           3,
				ArgonMemory:         256 * 1024,
				ArgonThreads:        4,
				ArgonKeyLen:         32,
				DefaultLength:       128,
				DefaultLanguage:     "auto",
				CharacterSet:        "alphanumeric_symbols",
				DefaultClearTimeout: 0, // validate не изменяет 0, только отрицательные значения
				Version:             "1.0",
			},
		},
		{
			name: "Пустые строковые поля",
			config: &Config{
				DefaultLanguage: "",
				CharacterSet:    "",
				Version:         "",
			},
			expected: &Config{
				ArgonTime:           3,
				ArgonMemory:         256 * 1024,
				ArgonThreads:        4,
				ArgonKeyLen:         32,
				DefaultLength:       16,
				DefaultLanguage:     "auto",
				CharacterSet:        "alphanumeric_symbols",
				DefaultClearTimeout: 0, // validate не изменяет 0, только отрицательные значения
				Version:             "1.0",
			},
		},
		{
			name: "Отрицательный таймаут",
			config: &Config{
				DefaultClearTimeout: -10,
			},
			expected: &Config{
				ArgonTime:           3,
				ArgonMemory:         256 * 1024,
				ArgonThreads:        4,
				ArgonKeyLen:         32,
				DefaultLength:       16,
				DefaultLanguage:     "auto",
				CharacterSet:        "alphanumeric_symbols",
				DefaultClearTimeout: 45,
				Version:             "1.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.validate()

			if tt.config.ArgonTime != tt.expected.ArgonTime {
				t.Errorf("validate() ArgonTime = %v, ожидается %v", tt.config.ArgonTime, tt.expected.ArgonTime)
			}
			if tt.config.ArgonMemory != tt.expected.ArgonMemory {
				t.Errorf("validate() ArgonMemory = %v, ожидается %v", tt.config.ArgonMemory, tt.expected.ArgonMemory)
			}
			if tt.config.ArgonThreads != tt.expected.ArgonThreads {
				t.Errorf("validate() ArgonThreads = %v, ожидается %v", tt.config.ArgonThreads, tt.expected.ArgonThreads)
			}
			if tt.config.ArgonKeyLen != tt.expected.ArgonKeyLen {
				t.Errorf("validate() ArgonKeyLen = %v, ожидается %v", tt.config.ArgonKeyLen, tt.expected.ArgonKeyLen)
			}
			if tt.config.DefaultLength != tt.expected.DefaultLength {
				t.Errorf("validate() DefaultLength = %v, ожидается %v", tt.config.DefaultLength, tt.expected.DefaultLength)
			}
			if tt.config.DefaultLanguage != tt.expected.DefaultLanguage {
				t.Errorf("validate() DefaultLanguage = %v, ожидается %v", tt.config.DefaultLanguage, tt.expected.DefaultLanguage)
			}
			if tt.config.CharacterSet != tt.expected.CharacterSet {
				t.Errorf("validate() CharacterSet = %v, ожидается %v", tt.config.CharacterSet, tt.expected.CharacterSet)
			}
			if tt.config.DefaultClearTimeout != tt.expected.DefaultClearTimeout {
				t.Errorf("validate() DefaultClearTimeout = %v, ожидается %v", tt.config.DefaultClearTimeout, tt.expected.DefaultClearTimeout)
			}
			if tt.config.Version != tt.expected.Version {
				t.Errorf("validate() Version = %v, ожидается %v", tt.config.Version, tt.expected.Version)
			}
		})
	}
}

func TestLoadConfigNotExists(t *testing.T) {
	// Тест загрузки конфигурации когда файл не существует
	messages := &i18n.Messages{
		ConfigErrorReading: "Ошибка чтения конфигурации: %v",
		ConfigErrorParsing: "Ошибка парсинга конфигурации: %v",
	}

	// Создаем временную директорию и удаляем её, чтобы путь не существовал
	tempDir, err := os.MkdirTemp("", "pgen_test_*")
	if err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}
	nonExistentPath := filepath.Join(tempDir, "nonexistent", "config.json")
	os.RemoveAll(tempDir)

	// Сохраняем оригинальные переменные окружения
	originalAppData := os.Getenv("APPDATA")
	originalHome := os.Getenv("HOME")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")

	defer func() {
		os.Setenv("APPDATA", originalAppData)
		os.Setenv("HOME", originalHome)
		os.Setenv("XDG_CONFIG_HOME", originalXDG)
	}()

	// Устанавливаем переменные окружения так, чтобы GetConfigPath вернул несуществующий путь
	os.Setenv("APPDATA", filepath.Dir(nonExistentPath))
	os.Setenv("HOME", filepath.Dir(nonExistentPath))
	os.Setenv("XDG_CONFIG_HOME", filepath.Dir(nonExistentPath))

	config, err := Load(messages)
	if err != nil {
		t.Errorf("Load() с несуществующим файлом не должен возвращать ошибку, получена: %v", err)
	}

	if config == nil {
		t.Fatal("Load() не должен возвращать nil конфигурацию")
	}

	// Должна быть возвращена конфигурация по умолчанию
	defaultConfig := DefaultConfig()
	if config.DefaultLength != defaultConfig.DefaultLength {
		t.Errorf("Load() должен возвращать конфигурацию по умолчанию")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	// Тест загрузки конфигурации с невалидным JSON
	// Этот тест проверяет логику парсинга JSON напрямую
	messages := &i18n.Messages{
		ConfigErrorReading: "Ошибка чтения конфигурации: %v",
		ConfigErrorParsing: "Ошибка парсинга конфигурации: %v",
	}

	// Создаем временный файл с невалидным JSON в реальной конфигурационной директории
	configPath, err := GetConfigPath()
	if err != nil {
		t.Skipf("Не удалось получить путь к конфигурации: %v", err)
	}

	// Сохраняем оригинальный файл конфигурации если он существует
	originalData, _ := os.ReadFile(configPath)
	defer func() {
		if originalData != nil {
			os.WriteFile(configPath, originalData, 0644)
		} else {
			os.Remove(configPath)
		}
	}()

	// Записываем невалидный JSON
	invalidJSON := `{"invalid": json content`
	if err := os.WriteFile(configPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Не удалось записать невалидный JSON: %v", err)
	}

	config, err := Load(messages)
	if err == nil {
		t.Error("Load() с невалидным JSON должен возвращать ошибку")
	}

	if config != nil {
		t.Error("Load() с невалидным JSON должен возвращать nil конфигурацию")
	}

	if !strings.Contains(err.Error(), "Ошибка парсинга конфигурации") {
		t.Errorf("Load() должен возвращать ошибку парсинга, получена: %v", err)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Тест сохранения и загрузки конфигурации
	messages := &i18n.Messages{
		ConfigErrorReading:  "Ошибка чтения конфигурации: %v",
		ConfigErrorParsing:  "Ошибка парсинга конфигурации: %v",
		ConfigErrorEncoding: "Ошибка кодирования конфигурации: %v",
		ConfigErrorWriting:  "Ошибка записи конфигурации: %v",
	}

	// Получаем реальный путь к конфигурации
	configPath, err := GetConfigPath()
	if err != nil {
		t.Skipf("Не удалось получить путь к конфигурации: %v", err)
	}

	// Сохраняем оригинальный файл конфигурации если он существует
	originalData, _ := os.ReadFile(configPath)
	defer func() {
		if originalData != nil {
			os.WriteFile(configPath, originalData, 0644)
		} else {
			os.Remove(configPath)
		}
	}()

	// Создаем тестовую конфигурацию
	testConfig := &Config{
		ArgonTime:           5,
		ArgonMemory:         512 * 1024,
		ArgonThreads:        8,
		ArgonKeyLen:         64,
		DefaultLength:       24,
		DefaultLanguage:     "ru",
		CharacterSet:        "custom",
		DefaultCopy:         true,
		DefaultClearTimeout: 60,
		ShowPasswordInfo:    true,
		ColorOutput:         false,
		Version:             "2.0",
	}

	// Сохраняем конфигурацию
	if err := testConfig.Save(messages); err != nil {
		t.Fatalf("Save() ошибка = %v", err)
	}

	// Проверяем, что файл создан
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Save() должен создать файл %q", configPath)
	}

	// Загружаем конфигурацию
	loadedConfig, err := Load(messages)
	if err != nil {
		t.Fatalf("Load() ошибка = %v", err)
	}

	// Сравниваем конфигурации
	if loadedConfig.ArgonTime != testConfig.ArgonTime {
		t.Errorf("Load() ArgonTime = %v, ожидается %v", loadedConfig.ArgonTime, testConfig.ArgonTime)
	}
	if loadedConfig.ArgonMemory != testConfig.ArgonMemory {
		t.Errorf("Load() ArgonMemory = %v, ожидается %v", loadedConfig.ArgonMemory, testConfig.ArgonMemory)
	}
	if loadedConfig.DefaultLength != testConfig.DefaultLength {
		t.Errorf("Load() DefaultLength = %v, ожидается %v", loadedConfig.DefaultLength, testConfig.DefaultLength)
	}
	if loadedConfig.DefaultLanguage != testConfig.DefaultLanguage {
		t.Errorf("Load() DefaultLanguage = %v, ожидается %v", loadedConfig.DefaultLanguage, testConfig.DefaultLanguage)
	}
}

func TestExportConfig(t *testing.T) {
	// Тест экспорта конфигурации
	messages := &i18n.Messages{
		ConfigErrorExportEncode: "Ошибка кодирования при экспорте: %v",
		ConfigErrorExportWrite:  "Ошибка записи при экспорте: %v",
	}

	// Создаем временный файл для экспорта
	tempFile, err := os.CreateTemp("", "export_config_*.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	config := DefaultConfig()
	config.DefaultLength = 20
	config.DefaultLanguage = "en"

	// Экспортируем конфигурацию
	if err := config.Export(tempFile.Name(), messages); err != nil {
		t.Fatalf("Export() ошибка = %v", err)
	}

	// Читаем экспортированный файл
	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Не удалось прочитать экспортированный файл: %v", err)
	}

	// Проверяем, что это валидный JSON
	var exportedConfig Config
	if err := json.Unmarshal(data, &exportedConfig); err != nil {
		t.Fatalf("Экспортированный файл содержит невалидный JSON: %v", err)
	}

	// Проверяем значения
	if exportedConfig.DefaultLength != 20 {
		t.Errorf("Export() DefaultLength = %v, ожидается 20", exportedConfig.DefaultLength)
	}
	if exportedConfig.DefaultLanguage != "en" {
		t.Errorf("Export() DefaultLanguage = %v, ожидается 'en'", exportedConfig.DefaultLanguage)
	}
}

func TestImportConfig(t *testing.T) {
	// Тест импорта конфигурации
	messages := &i18n.Messages{
		ConfigErrorImportRead:  "Ошибка чтения при импорте: %v",
		ConfigErrorImportParse: "Ошибка парсинга при импорте: %v",
	}

	// Создаем временный файл с конфигурацией
	tempFile, err := os.CreateTemp("", "import_config_*.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())

	testConfig := &Config{
		ArgonTime:           7,
		ArgonMemory:         1024 * 1024,
		ArgonThreads:        16,
		ArgonKeyLen:         128,
		DefaultLength:       32,
		DefaultLanguage:     "fr",
		CharacterSet:        "letters_only",
		DefaultCopy:         true,
		DefaultClearTimeout: 120,
		ShowPasswordInfo:    true,
		ColorOutput:         false,
		Version:             "3.0",
	}

	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Не удалось сериализовать тестовую конфигурацию: %v", err)
	}

	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Не удалось записать в временный файл: %v", err)
	}
	tempFile.Close()

	// Импортируем конфигурацию
	importedConfig, err := Import(tempFile.Name(), messages)
	if err != nil {
		t.Fatalf("Import() ошибка = %v", err)
	}

	// Проверяем значения
	if importedConfig.ArgonTime != testConfig.ArgonTime {
		t.Errorf("Import() ArgonTime = %v, ожидается %v", importedConfig.ArgonTime, testConfig.ArgonTime)
	}
	if importedConfig.DefaultLength != testConfig.DefaultLength {
		t.Errorf("Import() DefaultLength = %v, ожидается %v", importedConfig.DefaultLength, testConfig.DefaultLength)
	}
	if importedConfig.DefaultLanguage != testConfig.DefaultLanguage {
		t.Errorf("Import() DefaultLanguage = %v, ожидается %v", importedConfig.DefaultLanguage, testConfig.DefaultLanguage)
	}
}

func TestImportConfigInvalidFile(t *testing.T) {
	// Тест импорта из несуществующего файла
	messages := &i18n.Messages{
		ConfigErrorImportRead:  "Ошибка чтения при импорте: %v",
		ConfigErrorImportParse: "Ошибка парсинга при импорте: %v",
	}

	config, err := Import("/nonexistent/file.json", messages)
	if err == nil {
		t.Error("Import() с несуществующим файлом должен возвращать ошибку")
	}

	if config != nil {
		t.Error("Import() с несуществующим файлом должен возвращать nil конфигурацию")
	}

	if !strings.Contains(err.Error(), "Ошибка чтения при импорте") {
		t.Errorf("Import() должен возвращать ошибку чтения, получена: %v", err)
	}
}

func TestImportConfigInvalidJSON(t *testing.T) {
	// Тест импорта файла с невалидным JSON
	messages := &i18n.Messages{
		ConfigErrorImportRead:  "Ошибка чтения при импорте: %v",
		ConfigErrorImportParse: "Ошибка парсинга при импорте: %v",
	}

	// Создаем временный файл с невалидным JSON
	tempFile, err := os.CreateTemp("", "invalid_import_*.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name())

	invalidJSON := `{"broken": json, content}`
	if _, err := tempFile.WriteString(invalidJSON); err != nil {
		t.Fatalf("Не удалось записать в временный файл: %v", err)
	}
	tempFile.Close()

	config, err := Import(tempFile.Name(), messages)
	if err == nil {
		t.Error("Import() с невалидным JSON должен возвращать ошибку")
	}

	if config != nil {
		t.Error("Import() с невалидным JSON должен возвращать nil конфигурацию")
	}

	if !strings.Contains(err.Error(), "Ошибка парсинга при импорте") {
		t.Errorf("Import() должен возвращать ошибку парсинга, получена: %v", err)
	}
}

func TestConfigJSONSerialization(t *testing.T) {
	// Тест сериализации и десериализации JSON
	config := DefaultConfig()
	config.DefaultLength = 25
	config.DefaultLanguage = "de"
	config.ShowPasswordInfo = true

	// Сериализуем в JSON
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal() ошибка = %v", err)
	}

	// Десериализуем из JSON
	var deserializedConfig Config
	if err := json.Unmarshal(data, &deserializedConfig); err != nil {
		t.Fatalf("json.Unmarshal() ошибка = %v", err)
	}

	// Проверяем значения
	if deserializedConfig.DefaultLength != config.DefaultLength {
		t.Errorf("JSON сериализация DefaultLength = %v, ожидается %v", deserializedConfig.DefaultLength, config.DefaultLength)
	}
	if deserializedConfig.DefaultLanguage != config.DefaultLanguage {
		t.Errorf("JSON сериализация DefaultLanguage = %v, ожидается %v", deserializedConfig.DefaultLanguage, config.DefaultLanguage)
	}
	if deserializedConfig.ShowPasswordInfo != config.ShowPasswordInfo {
		t.Errorf("JSON сериализация ShowPasswordInfo = %v, ожидается %v", deserializedConfig.ShowPasswordInfo, config.ShowPasswordInfo)
	}
}

func TestConfigValidateEdgeCases(t *testing.T) {
	// Тест граничных случаев валидации
	tests := []struct {
		name           string
		defaultLength  int
		expectedLength int
	}{
		{"Минимальная длина", 4, 4},
		{"Максимальная длина", 128, 128},
		{"Длина на границе минимума", 3, 16},
		{"Длина на границе максимума", 129, 128},
		{"Нулевая длина", 0, 16},
		{"Отрицательная длина", -5, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{DefaultLength: tt.defaultLength}
			config.validate()

			if config.DefaultLength != tt.expectedLength {
				t.Errorf("validate() DefaultLength = %v, ожидается %v", config.DefaultLength, tt.expectedLength)
			}
		})
	}
}

// Бенчмарки для измерения производительности
func BenchmarkDefaultConfig(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DefaultConfig()
	}
}

func BenchmarkConfigValidate(b *testing.B) {
	config := &Config{
		ArgonTime:           0,
		ArgonMemory:         0,
		ArgonThreads:        0,
		ArgonKeyLen:         0,
		DefaultLength:       0,
		DefaultLanguage:     "",
		CharacterSet:        "",
		DefaultClearTimeout: -1,
		Version:             "",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.validate()
	}
}

func BenchmarkConfigJSONMarshal(b *testing.B) {
	config := DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(config)
	}
}

func BenchmarkConfigJSONUnmarshal(b *testing.B) {
	config := DefaultConfig()
	data, _ := json.Marshal(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var c Config
		_ = json.Unmarshal(data, &c)
	}
}
