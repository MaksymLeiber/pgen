package security

import (
	"sync"
)

// SecureMemoryPool управляет пулом безопасной памяти
type SecureMemoryPool struct {
	pool chan []byte
	size int
	mu   sync.Mutex
}

// NewSecureMemoryPool создает новый пул безопасной памяти
func NewSecureMemoryPool(poolSize, bufferSize int) *SecureMemoryPool {
	pool := make(chan []byte, poolSize)
	
	// Предварительно заполняем пул
	for i := 0; i < poolSize; i++ {
		pool <- make([]byte, bufferSize)
	}
	
	return &SecureMemoryPool{
		pool: pool,
		size: bufferSize,
	}
}

// Get получает буфер из пула
func (p *SecureMemoryPool) Get() []byte {
	select {
	case buf := <-p.pool:
		// Очищаем буфер перед использованием
		ZeroMemory(buf)
		return buf
	default:
		// Если пул пуст, создаем новый буфер
		return make([]byte, p.size)
	}
}

// Put возвращает буфер в пул после безопасной очистки
func (p *SecureMemoryPool) Put(buf []byte) {
	if buf == nil || len(buf) != p.size {
		return
	}
	
	// Безопасно очищаем буфер
	SecureWipe(buf)
	
	select {
	case p.pool <- buf:
		// Буфер успешно возвращен в пул
	default:
		// Пул полон, буфер будет собран сборщиком мусора
	}
}

// Size возвращает размер буферов в пуле
func (p *SecureMemoryPool) Size() int {
	return p.size
}

// Available возвращает количество доступных буферов в пуле
func (p *SecureMemoryPool) Available() int {
	return len(p.pool)
}

// SecureBuffer представляет безопасный буфер с автоматической очисткой
type SecureBuffer struct {
	data []byte
	pool *SecureMemoryPool
}

// NewSecureBuffer создает новый безопасный буфер
func NewSecureBuffer(pool *SecureMemoryPool) *SecureBuffer {
	return &SecureBuffer{
		data: pool.Get(),
		pool: pool,
	}
}

// Data возвращает данные буфера
func (sb *SecureBuffer) Data() []byte {
	return sb.data
}

// Write записывает данные в буфер
func (sb *SecureBuffer) Write(data []byte) int {
	if sb.data == nil {
		return 0
	}
	
	n := copy(sb.data, data)
	return n
}

// WriteString записывает строку в буфер
func (sb *SecureBuffer) WriteString(s string) int {
	return sb.Write([]byte(s))
}

// Clear очищает буфер
func (sb *SecureBuffer) Clear() {
	if sb.data != nil {
		ZeroMemory(sb.data)
	}
}

// Close освобождает буфер и возвращает его в пул
func (sb *SecureBuffer) Close() {
	if sb.data != nil && sb.pool != nil {
		sb.pool.Put(sb.data)
		sb.data = nil
		sb.pool = nil
	}
}

// SecureAllocator управляет выделением безопасной памяти
type SecureAllocator struct {
	pools map[int]*SecureMemoryPool
	mu    sync.RWMutex
}

// NewSecureAllocator создает новый аллокатор безопасной памяти
func NewSecureAllocator() *SecureAllocator {
	return &SecureAllocator{
		pools: make(map[int]*SecureMemoryPool),
	}
}

// GetPool возвращает пул для указанного размера буфера
func (sa *SecureAllocator) GetPool(size int) *SecureMemoryPool {
	sa.mu.RLock()
	pool, exists := sa.pools[size]
	sa.mu.RUnlock()
	
	if exists {
		return pool
	}
	
	sa.mu.Lock()
	defer sa.mu.Unlock()
	
	// Проверяем еще раз под write lock
	if pool, exists := sa.pools[size]; exists {
		return pool
	}
	
	// Создаем новый пул
	pool = NewSecureMemoryPool(10, size) // 10 буферов по умолчанию
	sa.pools[size] = pool
	return pool
}

// Allocate выделяет безопасный буфер указанного размера
func (sa *SecureAllocator) Allocate(size int) *SecureBuffer {
	pool := sa.GetPool(size)
	return NewSecureBuffer(pool)
}

// Stats возвращает статистику использования памяти
func (sa *SecureAllocator) Stats() map[int]int {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	
	stats := make(map[int]int)
	for size, pool := range sa.pools {
		stats[size] = pool.Available()
	}
	return stats
}

// DefaultAllocator - глобальный экземпляр аллокатора
var DefaultAllocator = NewSecureAllocator()

// AllocateSecure выделяет безопасный буфер через глобальный аллокатор
func AllocateSecure(size int) *SecureBuffer {
	return DefaultAllocator.Allocate(size)
}
