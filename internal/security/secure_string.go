package security

import (
	"crypto/rand"
	"runtime"
	"runtime/debug"
	"unsafe"
)

// SecureString представляет безопасную строку для хранения чувствительных данных
type SecureString struct {
	data []byte
	size int
}

// NewSecureString создает новую безопасную строку из обычной строки
func NewSecureString(s string) *SecureString {
	data := make([]byte, len(s))
	copy(data, []byte(s))
	return &SecureString{
		data: data,
		size: len(s),
	}
}

// NewSecureStringFromBytes создает новую безопасную строку из байтов
func NewSecureStringFromBytes(b []byte) *SecureString {
	data := make([]byte, len(b))
	copy(data, b)
	return &SecureString{
		data: data,
		size: len(b),
	}
}

// String возвращает строковое представление (небезопасно, использовать осторожно)
func (s *SecureString) String() string {
	if s == nil || s.data == nil {
		return ""
	}
	return string(s.data[:s.size])
}

// Bytes возвращает копию байтов
func (s *SecureString) Bytes() []byte {
	if s == nil || s.data == nil {
		return nil
	}
	result := make([]byte, s.size)
	copy(result, s.data[:s.size])
	return result
}

// Len возвращает длину строки
func (s *SecureString) Len() int {
	if s == nil {
		return 0
	}
	return s.size
}

// IsEmpty проверяет, пуста ли строка
func (s *SecureString) IsEmpty() bool {
	return s == nil || s.size == 0
}

// Clear безопасно очищает содержимое строки
func (s *SecureString) Clear() {
	if s == nil || s.data == nil {
		return
	}
	
	// Перезаписываем данные случайными байтами несколько раз
	for i := 0; i < 3; i++ {
		rand.Read(s.data)
	}
	
	// Затем заполняем нулями
	for i := range s.data {
		s.data[i] = 0
	}
	
	s.size = 0
	s.data = nil
	
	// Принудительно запускаем сборщик мусора
	runtime.GC()
	debug.FreeOSMemory()
}

// SecureCompare безопасно сравнивает две SecureString за константное время
func (s *SecureString) SecureCompare(other *SecureString) bool {
	if s == nil && other == nil {
		return true
	}
	if s == nil || other == nil {
		return false
	}
	if s.size != other.size {
		return false
	}
	
	return constantTimeCompare(s.data[:s.size], other.data[:other.size])
}

// constantTimeCompare выполняет сравнение за константное время
func constantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	
	return result == 0
}

// ZeroMemory безопасно обнуляет память
func ZeroMemory(data []byte) {
	if len(data) == 0 {
		return
	}
	
	// Используем unsafe для гарантированной очистки
	ptr := unsafe.Pointer(&data[0])
	for i := 0; i < len(data); i++ {
		*(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i))) = 0
	}
}

// SecureWipe безопасно очищает буфер перезаписью случайными данными
func SecureWipe(buf []byte) {
	if len(buf) == 0 {
		return
	}
	
	// Перезаписываем случайными данными
	rand.Read(buf)
	
	// Затем заполняем нулями
	ZeroMemory(buf)
}

// SecureRandom генерирует криптографически стойкие случайные байты
func SecureRandom(size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
