package generator

import (
	"testing"

	"github.com/MaksymLeiber/pgen/internal/i18n"
	"github.com/MaksymLeiber/pgen/internal/security"
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
			masterPassword := security.NewSecureString(tt.masterPassword)
			defer masterPassword.Clear()
			
			password, err := pg.GeneratePassword(masterPassword, tt.serviceName, "testuser", messages)

			if tt.wantError && err == nil {
				t.Error("GeneratePassword() ожидалась ошибка, получен nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("GeneratePassword() неожиданная ошибка: %v", err)
			}

			if !tt.wantError {
				defer password.Clear()
				
				if password.Len() != pg.length {
					t.Errorf("GeneratePassword() длина пароля = %v, ожидается %v", password.Len(), pg.length)
				}

				if password.IsEmpty() {
					t.Error("GeneratePassword() пустой пароль")
				}
			}
		})
	}
}

func TestGeneratePasswordDeterministic(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")
	pg := NewPasswordGenerator(16)

	masterPassword := security.NewSecureString("testmaster123")
	defer masterPassword.Clear()
	serviceName := "github.com"
	username := "testuser"

	// Генерируем пароль несколько раз
	password1, err := pg.GeneratePassword(masterPassword, serviceName, username, messages)
	if err != nil {
		t.Fatalf("GeneratePassword() ошибка: %v", err)
	}
	defer password1.Clear()

	password2, err := pg.GeneratePassword(masterPassword, serviceName, username, messages)
	if err != nil {
		t.Fatalf("GeneratePassword() ошибка: %v", err)
	}
	defer password2.Clear()

	password3, err := pg.GeneratePassword(masterPassword, serviceName, username, messages)
	if err != nil {
		t.Fatalf("GeneratePassword() ошибка: %v", err)
	}
	defer password3.Clear()

	// Все пароли должны быть одинаковыми
	if !password1.SecureCompare(password2) || !password2.SecureCompare(password3) {
		t.Errorf("Пароли не детерминированы: %s, %s, %s", password1.String(), password2.String(), password3.String())
	}
}

func TestGeneratePasswordDifferentInputs(t *testing.T) {
	messages := i18n.GetMessages(i18n.English, "test")
	pg := NewPasswordGenerator(16)

	master1 := security.NewSecureString("master1")
	defer master1.Clear()
	master2 := security.NewSecureString("master2")
	defer master2.Clear()

	password1, _ := pg.GeneratePassword(master1, "service1", "user1", messages)
	defer password1.Clear()
	password2, _ := pg.GeneratePassword(master2, "service1", "user1", messages)
	defer password2.Clear()
	password3, _ := pg.GeneratePassword(master1, "service2", "user1", messages)
	defer password3.Clear()

	if password1.SecureCompare(password2) {
		t.Error("Разные мастер-пароли должны генерировать разные пароли")
	}

	if password1.SecureCompare(password3) {
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
	masterPassword := security.NewSecureString("testmaster")
	defer masterPassword.Clear()
	
	password, err := pg.GeneratePassword(masterPassword, "testservice", "testuser", messages)

	if err != nil {
		t.Fatalf("GeneratePassword() ошибка: %v", err)
	}
	defer password.Clear()

	if password.Len() != 20 {
		t.Errorf("Длина пароля = %v, ожидается %v", password.Len(), 20)
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
			salt1 := createSalt(tt.serviceName, "testuser")
			salt2 := createSalt(tt.serviceName, "testuser")

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
	salt1 := createSalt("service1", "user1")
	salt2 := createSalt("service2", "user1")

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
	masterPassword := security.NewSecureString("benchmark")
	defer masterPassword.Clear()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		password, _ := pg.GeneratePassword(masterPassword, "test", "benchuser", messages)
		if password != nil {
			password.Clear()
		}
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
	masterPassword := security.NewSecureString("benchmark")
	defer masterPassword.Clear()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		password, _ := pg.GeneratePassword(masterPassword, "test", "benchuser", messages)
		if password != nil {
			password.Clear()
		}
	}
}
