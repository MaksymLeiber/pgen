package clipboard

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/atotto/clipboard"
)

func TestCopyToClipboard(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "Простой текст",
			text: "Hello World",
		},
		{
			name: "Пустая строка",
			text: "",
		},
		{
			name: "Пароль с символами",
			text: "MySecure123!@#",
		},
		{
			name: "Unicode текст",
			text: "Привет мир! 🔐",
		},
		{
			name: "Длинный текст",
			text: "This is a very long password that should still work correctly in the clipboard functionality",
		},
		{
			name: "Специальные символы",
			text: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name: "Многострочный текст",
			text: "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CopyToClipboard(tt.text)
			if err != nil {
				t.Errorf("CopyToClipboard() ошибка = %v", err)
				return
			}

			// Проверяем, что текст действительно скопирован
			clipboardContent, err := clipboard.ReadAll()
			if err != nil {
				t.Errorf("Не удалось прочитать буфер обмена: %v", err)
				return
			}

			if clipboardContent != tt.text {
				t.Errorf("CopyToClipboard() = %q, ожидается %q", clipboardContent, tt.text)
			}
		})
	}
}

func TestCopyToClipboardWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "Корректное копирование с таймаутом",
			text:    "TestPassword123",
			timeout: 100 * time.Millisecond,
			wantErr: false,
		},
		{
			name:    "Нулевой таймаут",
			text:    "NoTimeout",
			timeout: 0,
			wantErr: false,
		},
		{
			name:    "Отрицательный таймаут",
			text:    "NegativeTimeout",
			timeout: -1 * time.Second,
			wantErr: false,
		},
		{
			name:    "Очень короткий таймаут",
			text:    "ShortTimeout",
			timeout: 1 * time.Millisecond,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done, err := CopyToClipboardWithTimeout(tt.text, tt.timeout)

			if (err != nil) != tt.wantErr {
				t.Errorf("CopyToClipboardWithTimeout() ошибка = %v, ожидалась ошибка %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Если ошибка ожидалась, дальше не проверяем
			}

			// Проверяем, что текст скопирован
			clipboardContent, err := clipboard.ReadAll()
			if err != nil {
				t.Errorf("Не удалось прочитать буфер обмена: %v", err)
				return
			}

			if clipboardContent != tt.text {
				t.Errorf("CopyToClipboardWithTimeout() скопированный текст = %q, ожидается %q", clipboardContent, tt.text)
			}

			// Проверяем поведение канала в зависимости от таймаута
			if tt.timeout <= 0 {
				// При нулевом или отрицательном таймауте канал должен быть nil
				if done != nil {
					t.Error("CopyToClipboardWithTimeout() с нулевым таймаутом должен возвращать nil канал")
				}
			} else {
				// При положительном таймауте канал должен существовать
				if done == nil {
					t.Error("CopyToClipboardWithTimeout() с положительным таймаутом должен возвращать канал")
					return
				}

				// Ждем сигнала об очистке
				select {
				case cleared := <-done:
					if !cleared {
						t.Error("CopyToClipboardWithTimeout() канал должен отправить true при очистке")
					}
				case <-time.After(tt.timeout + 50*time.Millisecond):
					t.Error("CopyToClipboardWithTimeout() таймаут не сработал вовремя")
				}

				// Проверяем, что буфер очищен
				clearedContent, err := clipboard.ReadAll()
				if err != nil {
					t.Errorf("Не удалось прочитать буфер после очистки: %v", err)
					return
				}

				if clearedContent != "" {
					t.Errorf("CopyToClipboardWithTimeout() буфер не очищен: %q", clearedContent)
				}
			}
		})
	}
}

func TestCopyToClipboardWithTimeoutMultiple(t *testing.T) {
	// Тест последовательного использования с таймаутами
	tests := []struct {
		text    string
		timeout time.Duration
	}{
		{"First", 30 * time.Millisecond},
		{"Second", 40 * time.Millisecond},
		{"Third", 50 * time.Millisecond},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Множественный_%d", i+1), func(t *testing.T) {
			done, err := CopyToClipboardWithTimeout(tt.text, tt.timeout)
			if err != nil {
				t.Errorf("CopyToClipboardWithTimeout() ошибка = %v", err)
				return
			}

			if done != nil {
				select {
				case cleared := <-done:
					if !cleared {
						t.Error("Канал должен отправить true при очистке")
					}
				case <-time.After(tt.timeout + 100*time.Millisecond):
					t.Error("Таймаут не сработал вовремя")
				}
			}

			// Небольшая пауза между тестами
			time.Sleep(10 * time.Millisecond)
		})
	}
}

func TestCopyToClipboardSequential(t *testing.T) {
	// Тест последовательного копирования разных значений
	values := []string{
		"First",
		"Second", 
		"Third",
		"Четвертый",
		"🔐Fifth",
	}

	for i, value := range values {
		t.Run(fmt.Sprintf("Последовательность_%d", i+1), func(t *testing.T) {
			err := CopyToClipboard(value)
			if err != nil {
				t.Errorf("CopyToClipboard() ошибка = %v", err)
				return
			}

			// Небольшая задержка для стабильности
			time.Sleep(10 * time.Millisecond)

			clipboardContent, err := clipboard.ReadAll()
			if err != nil {
				t.Errorf("Не удалось прочитать буфер: %v", err)
				return
			}

			if clipboardContent != value {
				t.Errorf("Последовательность %d: получен %q, ожидается %q", i+1, clipboardContent, value)
			}
		})
	}
}

func TestCopyToClipboardWithTimeoutChannelClosure(t *testing.T) {
	// Тест на корректное закрытие канала
	text := "ChannelTest"
	timeout := 50 * time.Millisecond

	done, err := CopyToClipboardWithTimeout(text, timeout)
	if err != nil {
		t.Fatalf("CopyToClipboardWithTimeout() ошибка = %v", err)
	}

	if done == nil {
		t.Fatal("CopyToClipboardWithTimeout() должен возвращать канал")
	}

	// Ждем первого сигнала
	select {
	case cleared, ok := <-done:
		if !ok {
			t.Error("Канал закрыт преждевременно")
			return
		}
		if !cleared {
			t.Error("Первый сигнал должен быть true")
		}
	case <-time.After(timeout + 100*time.Millisecond):
		t.Fatal("Таймаут при ожидании первого сигнала")
	}

	// Проверяем, что канал закрыт после первого сигнала
	select {
	case _, ok := <-done:
		if ok {
			t.Error("Канал должен быть закрыт после отправки сигнала")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Канал должен быть закрыт сразу после отправки сигнала")
	}
}

// Бенчмарки для измерения производительности
func BenchmarkCopyToClipboard(b *testing.B) {
	text := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CopyToClipboard(text)
	}
}

func BenchmarkCopyToClipboardWithTimeout(b *testing.B) {
	text := "BenchmarkPassword123!"
	timeout := 100 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = CopyToClipboardWithTimeout(text, timeout)
	}
}

func BenchmarkCopyToClipboardLongText(b *testing.B) {
	// Тест с длинным текстом
	longText := strings.Repeat("A very long password with many characters ", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CopyToClipboard(longText)
	}
}
