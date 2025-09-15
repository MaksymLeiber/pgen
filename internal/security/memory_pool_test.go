package security

import (
	"sync"
	"testing"
)

func TestSecureMemoryPool(t *testing.T) {
	t.Run("Создание пула памяти", func(t *testing.T) {
		pool := NewSecureMemoryPool(10, 1024)
		
		if pool == nil {
			t.Fatal("Пул памяти не должен быть nil")
		}
		
		if pool.Size() != 1024 {
			t.Errorf("Ожидался размер буфера 1024, получено %d", pool.Size())
		}
		
		if pool.Available() != 10 {
			t.Errorf("Ожидалось 10 доступных буферов, получено %d", pool.Available())
		}
	})
	
	t.Run("Получение и возврат буфера", func(t *testing.T) {
		pool := NewSecureMemoryPool(5, 256)
		
		// Получаем буфер
		buffer := pool.Get()
		if buffer == nil {
			t.Fatal("Буфер не должен быть nil")
		}
		
		if len(buffer) != 256 {
			t.Errorf("Ожидался размер буфера 256, получено %d", len(buffer))
		}
		
		availableBefore := pool.Available()
		
		// Возвращаем буфер
		pool.Put(buffer)
		
		availableAfter := pool.Available()
		if availableAfter <= availableBefore {
			t.Error("Количество доступных буферов должно увеличиться после возврата")
		}
	})
	
	t.Run("Переиспользование буферов", func(t *testing.T) {
		pool := NewSecureMemoryPool(3, 128)
		
		// Получаем буфер
		buffer1 := pool.Get()
		
		// Записываем данные
		copy(buffer1, []byte("тестовые данные"))
		
		// Возвращаем буфер
		pool.Put(buffer1)
		
		// Получаем буфер снова
		buffer2 := pool.Get()
		
		// Проверяем, что данные очищены
		for i, b := range buffer2 {
			if b != 0 {
				t.Errorf("Буфер не очищен в позиции %d: получено %d", i, b)
			}
		}
	})
	
	t.Run("Поведение при пустом пуле", func(t *testing.T) {
		pool := NewSecureMemoryPool(1, 64)
		
		// Получаем все буферы из пула
		buffer1 := pool.Get()
		buffer2 := pool.Get() // Этот буфер должен быть создан новый
		
		if len(buffer1) != 64 || len(buffer2) != 64 {
			t.Error("Все буферы должны иметь правильный размер")
		}
		
		// Возвращаем буферы
		pool.Put(buffer1)
		pool.Put(buffer2)
	})
}

func TestSecureBuffer(t *testing.T) {
	t.Run("Создание безопасного буфера", func(t *testing.T) {
		pool := NewSecureMemoryPool(5, 512)
		buffer := NewSecureBuffer(pool)
		
		if buffer == nil {
			t.Fatal("Буфер не должен быть nil")
		}
		
		if len(buffer.Data()) != 512 {
			t.Errorf("Ожидался размер 512, получено %d", len(buffer.Data()))
		}
	})
	
	t.Run("Запись и очистка буфера", func(t *testing.T) {
		pool := NewSecureMemoryPool(3, 64)
		buffer := NewSecureBuffer(pool)
		
		// Записываем данные
		testData := []byte("конфиденциальные данные")
		n := buffer.Write(testData)
		
		if n != len(testData) {
			t.Errorf("Ожидалось записать %d байт, записано %d", len(testData), n)
		}
		
		// Проверяем, что данные записались
		hasData := false
		for _, b := range buffer.Data()[:len(testData)] {
			if b != 0 {
				hasData = true
				break
			}
		}
		if !hasData {
			t.Error("Данные должны быть записаны в буфер")
		}
		
		// Очищаем буфер
		buffer.Clear()
		
		// Проверяем, что данные очищены
		for i, b := range buffer.Data() {
			if b != 0 {
				t.Errorf("Буфер не очищен в позиции %d: получено %d", i, b)
			}
		}
	})
	
	t.Run("Запись строки в буфер", func(t *testing.T) {
		pool := NewSecureMemoryPool(2, 32)
		buffer := NewSecureBuffer(pool)
		
		testString := "тестовая строка"
		n := buffer.WriteString(testString)
		
		if n != len(testString) {
			t.Errorf("Ожидалось записать %d байт, записано %d", len(testString), n)
		}
	})
	
	t.Run("Закрытие буфера", func(t *testing.T) {
		pool := NewSecureMemoryPool(2, 128)
		availableBefore := pool.Available()
		
		buffer := NewSecureBuffer(pool)
		availableAfterGet := pool.Available()
		
		if availableAfterGet >= availableBefore {
			t.Error("Количество доступных буферов должно уменьшиться после получения")
		}
		
		buffer.Close()
		availableAfterClose := pool.Available()
		
		if availableAfterClose <= availableAfterGet {
			t.Error("Количество доступных буферов должно увеличиться после закрытия")
		}
	})
}

func TestSecureAllocator(t *testing.T) {
	t.Run("Выделение буферов разных размеров", func(t *testing.T) {
		allocator := NewSecureAllocator()
		
		// Выделяем буферы разных размеров
		buffer1 := allocator.Allocate(256)
		buffer2 := allocator.Allocate(512)
		buffer3 := allocator.Allocate(256) // Тот же размер, что и buffer1
		
		if buffer1 == nil || buffer2 == nil || buffer3 == nil {
			t.Fatal("Буферы не должны быть nil")
		}
		
		if len(buffer1.Data()) != 256 {
			t.Errorf("Ожидался размер 256, получено %d", len(buffer1.Data()))
		}
		
		if len(buffer2.Data()) != 512 {
			t.Errorf("Ожидался размер 512, получено %d", len(buffer2.Data()))
		}
		
		// Проверяем статистику
		stats := allocator.Stats()
		if len(stats) == 0 {
			t.Error("Статистика не должна быть пустой")
		}
		
		// Закрываем буферы
		buffer1.Close()
		buffer2.Close()
		buffer3.Close()
	})
	
	t.Run("Получение пулов для разных размеров", func(t *testing.T) {
		allocator := NewSecureAllocator()
		
		pool1 := allocator.GetPool(128)
		pool2 := allocator.GetPool(256)
		pool3 := allocator.GetPool(128) // Тот же размер
		
		if pool1 == nil || pool2 == nil || pool3 == nil {
			t.Fatal("Пулы не должны быть nil")
		}
		
		// Пулы для одинакового размера должны быть одним и тем же объектом
		if pool1 != pool3 {
			t.Error("Пулы для одинакового размера должны быть одним объектом")
		}
		
		// Пулы для разных размеров должны быть разными
		if pool1 == pool2 {
			t.Error("Пулы для разных размеров должны быть разными")
		}
	})
	
	t.Run("Глобальный аллокатор", func(t *testing.T) {
		buffer := AllocateSecure(64)
		
		if buffer == nil {
			t.Fatal("Буфер не должен быть nil")
		}
		
		if len(buffer.Data()) != 64 {
			t.Errorf("Ожидался размер 64, получено %d", len(buffer.Data()))
		}
		
		buffer.Close()
	})
}

func TestConcurrentAccess(t *testing.T) {
	t.Run("Конкурентный доступ к пулу памяти", func(t *testing.T) {
		pool := NewSecureMemoryPool(20, 256)
		
		var wg sync.WaitGroup
		numGoroutines := 10
		operationsPerGoroutine := 100
		
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				for j := 0; j < operationsPerGoroutine; j++ {
					buffer := pool.Get()
					if buffer != nil {
						// Записываем данные
						copy(buffer, []byte("тест"))
						// Возвращаем буфер
						pool.Put(buffer)
					}
				}
			}()
		}
		
		wg.Wait()
		
		// Проверяем, что пул все еще работает
		finalBuffer := pool.Get()
		if finalBuffer == nil {
			t.Error("Пул должен продолжать работать после конкурентного доступа")
		} else {
			pool.Put(finalBuffer)
		}
	})
}

// Бенчмарки
func BenchmarkПулПамяти_Получение(b *testing.B) {
	pool := NewSecureMemoryPool(100, 1024)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := pool.Get()
		pool.Put(buffer)
	}
}

func BenchmarkПулПамяти_БезПула(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer := make([]byte, 1024)
		SecureWipe(buffer)
	}
}

func BenchmarkБуфер_Очистка(b *testing.B) {
	pool := NewSecureMemoryPool(10, 1024)
	buffer := NewSecureBuffer(pool)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Clear()
	}
}
