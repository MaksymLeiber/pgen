package clipboard

import (
	"time"

	"github.com/atotto/clipboard"
)

// CopyToClipboard копирует текст в буфер обмена
func CopyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}

// CopyToClipboardWithTimeout копирует текст в буфер и возвращает канал для ожидания очистки
func CopyToClipboardWithTimeout(text string, timeout time.Duration) (<-chan bool, error) {
	if err := clipboard.WriteAll(text); err != nil {
		return nil, err
	}

	// Если таймаут 0 - не очищаем автоматически
	if timeout <= 0 {
		return nil, nil
	}

	// Создаем канал для уведомления об очистке
	done := make(chan bool, 1)

	// Запускаем очистку в отдельной горутине
	go func() {
		time.Sleep(timeout)
		clipboard.WriteAll("") // Очищаем буфер
		done <- true
		close(done)
	}()

	return done, nil
}

// ClearClipboard немедленно очищает буфер обмена
func ClearClipboard() error {
	return clipboard.WriteAll("")
}
