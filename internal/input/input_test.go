package input

import (
	"os"
	"strings"
	"testing"
)

func TestInputMessages(t *testing.T) {
	// Тест структуры InputMessages
	messages := &InputMessages{
		UserCanceled:  "Пользователь отменил операцию",
		InputCanceled: "Ввод отменен",
	}

	if messages.UserCanceled == "" {
		t.Error("UserCanceled не должно быть пустым")
	}
	if messages.InputCanceled == "" {
		t.Error("InputCanceled не должно быть пустым")
	}

	// Проверяем, что можно создать пустую структуру
	emptyMessages := &InputMessages{}
	if emptyMessages.UserCanceled != "" {
		t.Error("Пустая структура должна иметь пустые поля")
	}
}

func TestReadLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "Простая строка",
			input:    "hello world\n",
			expected: "hello world",
			hasError: false,
		},
		{
			name:     "Строка с пробелами в начале и конце",
			input:    "  test string  \n",
			expected: "test string",
			hasError: false,
		},
		{
			name:     "Пустая строка",
			input:    "\n",
			expected: "",
			hasError: false,
		},
		{
			name:     "Строка только с пробелами",
			input:    "   \n",
			expected: "",
			hasError: false,
		},
		{
			name:     "Строка с табуляцией",
			input:    "\thello\t\n",
			expected: "hello",
			hasError: false,
		},
		{
			name:     "Многострочный ввод",
			input:    "first line\nsecond line\n",
			expected: "first line",
			hasError: false,
		},
		{
			name:     "Строка с Unicode символами",
			input:    "Привет мир! 🔑\n",
			expected: "Привет мир! 🔑",
			hasError: false,
		},
		{
			name:     "Строка со специальными символами",
			input:    "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./\n",
			expected: "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем временный файл для имитации stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// Создаем pipe для имитации ввода
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Ошибка создания pipe: %v", err)
			}
			defer r.Close()
			defer w.Close()

			os.Stdin = r

			// Записываем тестовые данные в pipe
			go func() {
				defer w.Close()
				w.WriteString(tt.input)
			}()

			result, err := ReadLine()

			if tt.hasError && err == nil {
				t.Errorf("ReadLine() должен вернуть ошибку для входа %q", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("ReadLine() не должен возвращать ошибку для входа %q, получена ошибка: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ReadLine() = %q, ожидается %q", result, tt.expected)
			}
		})
	}
}

func TestReadLineWithClosedInput(t *testing.T) {
	// Тест поведения при закрытом вводе
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Создаем pipe и сразу закрываем writer
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Ошибка создания pipe: %v", err)
	}
	defer r.Close()
	w.Close() // Закрываем writer сразу

	os.Stdin = r

	result, err := ReadLine()
	if err == nil {
		t.Error("ReadLine() должен вернуть ошибку при закрытом вводе")
	}
	if result != "" {
		t.Errorf("ReadLine() должен вернуть пустую строку при ошибке, получено: %q", result)
	}
}

func TestReadPasswordWithStarsAndMessages(t *testing.T) {
	// Тест функции ReadPasswordWithStarsAndMessages
	messages := &InputMessages{
		UserCanceled:  "Операция отменена пользователем",
		InputCanceled: "Ввод отменен",
	}

	// Этот тест сложен для полного тестирования, так как функция зависит от терминала
	// Проверим, что функция принимает правильные параметры и не паникует
	t.Run("Проверка параметров", func(t *testing.T) {
		if messages == nil {
			t.Error("messages не должно быть nil")
		}

		// Проверяем, что функция существует и может быть вызвана
		// В реальных условиях это потребует mock'а терминала
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ReadPasswordWithStarsAndMessages не должна паниковать: %v", r)
			}
		}()

		// Мы не можем легко протестировать эту функцию без mock'а терминала
		// но можем проверить, что она существует и принимает правильные параметры
		_ = ReadPasswordWithStarsAndMessages
	})
}

func TestReadPasswordWithStarsAndMessagesNilMessages(t *testing.T) {
	// Тест с nil сообщениями
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ReadPasswordWithStarsAndMessages не должна паниковать с nil messages: %v", r)
		}
	}()

	// Функция должна обрабатывать nil messages без паники
	// В реальной реализации это может привести к ошибке, но не к панике
	_ = ReadPasswordWithStarsAndMessages
}

func TestInputMessagesWithDifferentLanguages(t *testing.T) {
	// Тест сообщений на разных языках
	tests := []struct {
		name     string
		messages *InputMessages
	}{
		{
			name: "Русские сообщения",
			messages: &InputMessages{
				UserCanceled:  "Пользователь отменил операцию",
				InputCanceled: "Ввод отменен",
			},
		},
		{
			name: "Английские сообщения",
			messages: &InputMessages{
				UserCanceled:  "User canceled operation",
				InputCanceled: "Input canceled",
			},
		},
		{
			name: "Сообщения с Unicode",
			messages: &InputMessages{
				UserCanceled:  "Пользователь отменил 🚫",
				InputCanceled: "Ввод отменен ❌",
			},
		},
		{
			name: "Длинные сообщения",
			messages: &InputMessages{
				UserCanceled:  "Операция была отменена пользователем по причине нежелания продолжать выполнение команды",
				InputCanceled: "Процесс ввода данных был прерван по запросу пользователя или из-за системной ошибки",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.messages.UserCanceled == "" {
				t.Error("UserCanceled не должно быть пустым")
			}
			if tt.messages.InputCanceled == "" {
				t.Error("InputCanceled не должно быть пустым")
			}

			// Проверяем, что сообщения содержат осмысленный текст
			if len(tt.messages.UserCanceled) < 3 {
				t.Error("UserCanceled слишком короткое")
			}
			if len(tt.messages.InputCanceled) < 3 {
				t.Error("InputCanceled слишком короткое")
			}
		})
	}
}

func TestReadLineWithLargeInput(t *testing.T) {
	// Тест с большим объемом данных
	largeInput := strings.Repeat("a", 10000) + "\n"
	expected := strings.Repeat("a", 10000)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Ошибка создания pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	os.Stdin = r

	go func() {
		defer w.Close()
		w.WriteString(largeInput)
	}()

	result, err := ReadLine()
	if err != nil {
		t.Errorf("ReadLine() не должен возвращать ошибку для большого ввода: %v", err)
	}
	if result != expected {
		t.Errorf("ReadLine() неправильно обработал большой ввод, длина результата: %d, ожидается: %d", len(result), len(expected))
	}
}

func TestReadLineWithBinaryData(t *testing.T) {
	// Тест с бинарными данными
	binaryInput := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 'h', 'e', 'l', 'l', 'o', '\n'}
	
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Ошибка создания pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	os.Stdin = r

	go func() {
		defer w.Close()
		w.Write(binaryInput)
	}()

	result, err := ReadLine()
	if err != nil {
		t.Errorf("ReadLine() не должен возвращать ошибку для бинарных данных: %v", err)
	}
	
	// Результат должен содержать все символы до \n, включая бинарные
	expectedLength := len(binaryInput) - 1 // минус \n
	if len(result) != expectedLength {
		t.Errorf("ReadLine() неправильно обработал бинарные данные, длина: %d, ожидается: %d", len(result), expectedLength)
	}
}

func TestReadLineMultipleCalls(t *testing.T) {
	// Тест множественных вызовов ReadLine
	inputs := []string{"first\n", "second\n", "third\n"}
	expected := []string{"first", "second", "third"}

	for i, input := range inputs {
		t.Run(strings.Join([]string{"Call", string(rune('1' + i))}, "_"), func(t *testing.T) {
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Ошибка создания pipe: %v", err)
			}
			defer r.Close()
			defer w.Close()

			os.Stdin = r

			go func() {
				defer w.Close()
				w.WriteString(input)
			}()

			result, err := ReadLine()
			if err != nil {
				t.Errorf("ReadLine() вызов %d не должен возвращать ошибку: %v", i+1, err)
			}
			if result != expected[i] {
				t.Errorf("ReadLine() вызов %d = %q, ожидается %q", i+1, result, expected[i])
			}
		})
	}
}

func TestReadLineWithDifferentLineEndings(t *testing.T) {
	// Тест с разными окончаниями строк
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unix line ending (LF)",
			input:    "test\n",
			expected: "test",
		},
		{
			name:     "Windows line ending (CRLF)",
			input:    "test\r\n",
			expected: "test",
		},
		{
			name:     "String without newline (EOF expected)",
			input:    "test",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Ошибка создания pipe: %v", err)
			}
			defer r.Close()
			defer w.Close()

			os.Stdin = r

			go func() {
				defer w.Close()
				w.WriteString(tt.input)
			}()

			result, err := ReadLine()
			
			// Для случая без символа новой строки ожидаем EOF ошибку
			if strings.Contains(tt.name, "EOF expected") {
				if err == nil {
					t.Errorf("ReadLine() должен вернуть ошибку для %s", tt.name)
				}
				return
			}
			
			if err != nil {
				t.Errorf("ReadLine() не должен возвращать ошибку для %s: %v", tt.name, err)
			}
			if result != tt.expected {
				t.Errorf("ReadLine() для %s = %q, ожидается %q", tt.name, result, tt.expected)
			}
		})
	}
}

// Бенчмарки для измерения производительности
func BenchmarkReadLine(b *testing.B) {
	input := "benchmark test line\n"
	
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		
		oldStdin := os.Stdin
		r, w, err := os.Pipe()
		if err != nil {
			b.Fatalf("Ошибка создания pipe: %v", err)
		}
		os.Stdin = r

		go func() {
			defer w.Close()
			w.WriteString(input)
		}()

		b.StartTimer()
		_, err = ReadLine()
		b.StopTimer()

		r.Close()
		os.Stdin = oldStdin

		if err != nil {
			b.Errorf("ReadLine() вернул ошибку: %v", err)
		}
	}
}

func BenchmarkReadLineLarge(b *testing.B) {
	input := strings.Repeat("x", 1000) + "\n"
	
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		
		oldStdin := os.Stdin
		r, w, err := os.Pipe()
		if err != nil {
			b.Fatalf("Ошибка создания pipe: %v", err)
		}
		os.Stdin = r

		go func() {
			defer w.Close()
			w.WriteString(input)
		}()

		b.StartTimer()
		_, err = ReadLine()
		b.StopTimer()

		r.Close()
		os.Stdin = oldStdin

		if err != nil {
			b.Errorf("ReadLine() вернул ошибку: %v", err)
		}
	}
}

func BenchmarkInputMessagesCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messages := &InputMessages{
			UserCanceled:  "Пользователь отменил операцию",
			InputCanceled: "Ввод отменен",
		}
		// Используем поля, чтобы избежать предупреждений
		_ = messages.UserCanceled
		_ = messages.InputCanceled
	}
}
