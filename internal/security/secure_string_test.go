package security

import (
	"testing"
)

func TestSecureString(t *testing.T) {
	t.Run("Создание SecureString", func(t *testing.T) {
		original := "test_password_123"
		ss := NewSecureString(original)
		
		if ss.String() != original {
			t.Errorf("Ожидалось %s, получено %s", original, ss.String())
		}
		
		if ss.Len() != len(original) {
			t.Errorf("Ожидалась длина %d, получено %d", len(original), ss.Len())
		}
		
		if ss.IsEmpty() {
			t.Error("SecureString не должна быть пустой")
		}
	})
	
	t.Run("Создание SecureString из байтов", func(t *testing.T) {
		original := []byte("test_password_123")
		ss := NewSecureStringFromBytes(original)
		
		if string(ss.Bytes()) != string(original) {
			t.Errorf("Ожидалось %s, получено %s", string(original), string(ss.Bytes()))
		}
	})
	
	t.Run("Очистка SecureString", func(t *testing.T) {
		ss := NewSecureString("sensitive_data")
		
		if ss.IsEmpty() {
			t.Error("SecureString не должна быть пустой до очистки")
		}
		
		ss.Clear()
		
		if !ss.IsEmpty() {
			t.Error("SecureString должна быть пустой после очистки")
		}
		
		if ss.Len() != 0 {
			t.Error("Длина должна быть 0 после очистки")
		}
	})
	
	t.Run("Безопасное сравнение", func(t *testing.T) {
		ss1 := NewSecureString("password123")
		ss2 := NewSecureString("password123")
		ss3 := NewSecureString("different")
		
		if !ss1.SecureCompare(ss2) {
			t.Error("Одинаковые строки должны быть равны")
		}
		
		if ss1.SecureCompare(ss3) {
			t.Error("Разные строки не должны быть равны")
		}
		
		// Тест с nil значениями
		var nilSS *SecureString
		if !nilSS.SecureCompare(nil) {
			t.Error("Два nil SecureString должны быть равны")
		}
		
		if ss1.SecureCompare(nilSS) {
			t.Error("SecureString и nil не должны быть равны")
		}
	})
}

func TestConstantTimeCompare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []byte
		expected bool
	}{
		{
			name:     "Одинаковые срезы",
			a:        []byte("password"),
			b:        []byte("password"),
			expected: true,
		},
		{
			name:     "Разные срезы одинаковой длины",
			a:        []byte("password"),
			b:        []byte("passwort"),
			expected: false,
		},
		{
			name:     "Разные длины",
			a:        []byte("password"),
			b:        []byte("pass"),
			expected: false,
		},
		{
			name:     "Пустые срезы",
			a:        []byte{},
			b:        []byte{},
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := constantTimeCompare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Ожидалось %v, получено %v", tt.expected, result)
			}
		})
	}
}

func TestZeroMemory(t *testing.T) {
	data := []byte("sensitive_data_here")
	
	// Проверяем, что данные не пусты
	hasNonZero := false
	for _, b := range data {
		if b != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Error("Тестовые данные должны содержать ненулевые байты")
	}
	
	// Очищаем память
	ZeroMemory(data)
	
	// Проверяем, что все байты обнулены
	for i, b := range data {
		if b != 0 {
			t.Errorf("Память не обнулена в позиции %d: получено %d", i, b)
		}
	}
}

func TestSecureWipe(t *testing.T) {
	data := []byte("sensitive_data_to_wipe")
	
	// Проверяем, что данные не пусты
	hasNonZero := false
	for _, b := range data {
		if b != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Error("Тестовые данные должны содержать ненулевые байты")
	}
	
	// Безопасно очищаем
	SecureWipe(data)
	
	// Проверяем, что все байты обнулены
	for i, b := range data {
		if b != 0 {
			t.Errorf("Память не очищена в позиции %d: получено %d", i, b)
		}
	}
}

func TestSecureRandom(t *testing.T) {
	size := 32
	random1, err := SecureRandom(size)
	if err != nil {
		t.Fatalf("SecureRandom завершилась с ошибкой: %v", err)
	}
	
	if len(random1) != size {
		t.Errorf("Ожидался размер %d, получено %d", size, len(random1))
	}
	
	// Генерируем второй набор случайных данных
	random2, err := SecureRandom(size)
	if err != nil {
		t.Fatalf("SecureRandom завершилась с ошибкой: %v", err)
	}
	
	// Проверяем, что они разные (с очень высокой вероятностью)
	if constantTimeCompare(random1, random2) {
		t.Error("Два случайных набора должны быть разными")
	}
}

// Бенчмарки
func BenchmarkСозданиеSecureString(b *testing.B) {
	password := "очень_длинный_пароль_для_бенчмарка_123456789"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ss := NewSecureString(password)
		ss.Clear()
	}
}

func BenchmarkБезопасноеСравнение(b *testing.B) {
	ss1 := NewSecureString("пароль_для_бенчмарка_123")
	ss2 := NewSecureString("пароль_для_бенчмарка_123")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ss1.SecureCompare(ss2)
	}
}

func BenchmarkОчисткаПамяти(b *testing.B) {
	data := make([]byte, 1024)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ZeroMemory(data)
	}
}

func BenchmarkБезопаснаяОчистка(b *testing.B) {
	data := make([]byte, 1024)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SecureWipe(data)
	}
}
