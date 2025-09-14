package generator

import (
	"testing"

	"pgen/internal/i18n"
)

func TestNewPasswordGenerator(t *testing.T) {
	tests := []struct {
		name           string
		length         int
		expectedLength int
	}{
		{"Корректная длина", 16, 16},
		{"Нулевая длина", 0, 16},
		{"Отрицательная длина", -5, 16},
		{"Большая длина", 64, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := NewPasswordGenerator(tt.length)
			if pg.length != tt.expectedLength {
				t.Errorf("NewPasswordGenerator() длина = %v, ожидается %v", pg.length, tt.expectedLength)
			}
			if pg.argon != nil {
				t.Error("NewPasswordGenerator() должен иметь пустую конфигурацию argon")
			}
		})
	}
}

func TestNewPasswordGeneratorWithConfig(t *testing.T) {
	config := ArgonConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 1,
		KeyLen:  32,
	}

	pg := NewPasswordGeneratorWithConfig(20, config)

	if pg.length != 20 {
		t.Errorf("NewPasswordGeneratorWithConfig() длина = %v, ожидается %v", pg.length, 20)
	}

	if pg.argon == nil {
		t.Fatal("NewPasswordGeneratorWithConfig() должен иметь конфигурацию argon")
	}

	if pg.argon.Time != config.Time {
		t.Errorf("ArgonConfig Time = %v, ожидается %v", pg.argon.Time, config.Time)
	}
}

func TestGeneratePassword(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")
	pg := NewPasswordGenerator(16)

	tests := []struct {
		name           string
		masterPassword string
		serviceName    string
		wantError      bool
	}{
		{"Корректные данные", "masterpass123", "github.com", false},
		{"Пустой мастер-пароль", "", "github.com", false},
		{"Пустое имя сервиса", "masterpass123", "", false},
		{"Оба поля пустые", "", "", false},
		{"Символы Unicode", "пароль123", "сайт.рф", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := pg.GeneratePassword(tt.masterPassword, tt.serviceName, messages)

			if tt.wantError && err == nil {
				t.Error("GeneratePassword() ожидалась ошибка, получен nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("GeneratePassword() неожиданная ошибка: %v", err)
			}

			if !tt.wantError {
				if len(password) != pg.length {
					t.Errorf("GeneratePassword() длина пароля = %v, ожидается %v", len(password), pg.length)
				}

				if password == "" {
					t.Error("GeneratePassword() вернул пустой пароль")
				}
			}
		})
	}
}

func TestGeneratePasswordDeterministic(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")
	pg := NewPasswordGenerator(16)

	masterPassword := "test123"
	serviceName := "example.com"

	// Генерируем пароль несколько раз
	password1, err1 := pg.GeneratePassword(masterPassword, serviceName, messages)
	password2, err2 := pg.GeneratePassword(masterPassword, serviceName, messages)
	password3, err3 := pg.GeneratePassword(masterPassword, serviceName, messages)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("GeneratePassword() ошибки: %v, %v, %v", err1, err2, err3)
	}

	if password1 != password2 || password2 != password3 {
		t.Errorf("GeneratePassword() не детерминирован: %v, %v, %v", password1, password2, password3)
	}
}

func TestGeneratePasswordDifferentInputs(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")
	pg := NewPasswordGenerator(16)

	password1, _ := pg.GeneratePassword("master1", "service1", messages)
	password2, _ := pg.GeneratePassword("master2", "service1", messages)
	password3, _ := pg.GeneratePassword("master1", "service2", messages)

	if password1 == password2 {
		t.Error("Разные мастер-пароли должны генерировать разные пароли")
	}

	if password1 == password3 {
		t.Error("Разные имена сервисов должны генерировать разные пароли")
	}
}

func TestGeneratePasswordWithCustomConfig(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")
	
	// Быстрые параметры для тестов
	config := ArgonConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 1,
		KeyLen:  32,
	}

	pg := NewPasswordGeneratorWithConfig(20, config)
	password, err := pg.GeneratePassword("testmaster", "testservice", messages)

	if err != nil {
		t.Fatalf("GeneratePassword() ошибка: %v", err)
	}

	if len(password) != 20 {
		t.Errorf("Длина пароля = %v, ожидается %v", len(password), 20)
	}
}

func TestCreateSalt(t *testing.T) {
	tests := []struct {
		serviceName string
	}{
		{"github.com"},
		{"google.com"},
		{""},
		{"тест.рф"},
	}

	for _, tt := range tests {
		t.Run(tt.serviceName, func(t *testing.T) {
			salt1 := createSalt(tt.serviceName)
			salt2 := createSalt(tt.serviceName)

			if len(salt1) != saltLength {
				t.Errorf("createSalt() длина = %v, ожидается %v", len(salt1), saltLength)
			}

			// Проверяем детерминированность
			if string(salt1) != string(salt2) {
				t.Error("createSalt() не детерминирован")
			}
		})
	}

	// Проверяем, что разные сервисы дают разные соли
	salt1 := createSalt("service1")
	salt2 := createSalt("service2")

	if string(salt1) == string(salt2) {
		t.Error("Разные сервисы должны генерировать разные соли")
	}
}

func TestValidateLength(t *testing.T) {
	tests := []struct {
		name      string
		length    int
		wantError bool
	}{
		{"Минимальная корректная", 4, false},
		{"Средняя корректная", 16, false},
		{"Максимальная корректная", 128, false},
		{"Слишком короткая", 3, true},
		{"Слишком длинная", 129, true},
		{"Нулевая", 0, true},
		{"Отрицательная", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLength(tt.length)

			if tt.wantError && err == nil {
				t.Error("ValidateLength() ожидалась ошибка, получен nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("ValidateLength() неожиданная ошибка: %v", err)
			}
		})
	}
}

func TestGenerateFromHash(t *testing.T) {
	pg := NewPasswordGenerator(10)
	
	// Тестовый хеш
	hash := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	
	password1 := pg.generateFromHash(hash)
	password2 := pg.generateFromHash(hash)

	// Проверяем детерминированность
	if password1 != password2 {
		t.Error("generateFromHash() не детерминирован")
	}

	// Проверяем, что пароль содержит только допустимые символы
	for _, char := range password1 {
		found := false
		for _, validChar := range charsetFull {
			if char == validChar {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("generateFromHash() содержит недопустимый символ: %c", char)
		}
	}
}

// Бенчмарки для измерения производительности
func BenchmarkGeneratePassword(b *testing.B) {
	messages := i18n.GetMessages(i18n.English, "test")
	pg := NewPasswordGenerator(16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = pg.GeneratePassword("benchmark", "test", messages)
	}
}

func BenchmarkGeneratePasswordFast(b *testing.B) {
	messages := i18n.GetMessages(i18n.English, "test")
	
	// Быстрые параметры для тестирования
	config := ArgonConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 1,
		KeyLen:  32,
	}
	
	pg := NewPasswordGeneratorWithConfig(16, config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = pg.GeneratePassword("benchmark", "test", messages)
	}
}
