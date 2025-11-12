package jexl

// Cache представляет интерфейс кэша выражений.
// Аналог org.apache.commons.jexl3.JexlCache.
type Cache[K comparable, V any] interface {
	// Get получает значение из кэша.
	Get(key K) (V, bool)
	// Put помещает значение в кэш.
	Put(key K, value V)
	// Clear очищает кэш.
	Clear()
	// Size возвращает текущий размер кэша.
	Size() int
	// Capacity возвращает максимальную ёмкость кэша.
	Capacity() int
	// Entries возвращает все записи кэша (для тестирования).
	Entries() []CacheEntry[K, V]
}

// CacheEntry представляет запись в кэше.
type CacheEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// CacheFactory создаёт кэш для заданного размера.
type CacheFactory func(size int) Cache[string, any]

// DefaultCacheFactory создаёт простой кэш на основе map.
var DefaultCacheFactory CacheFactory = func(size int) Cache[string, any] {
	return &simpleCache{
		data:     make(map[string]any, size),
		capacity: size,
	}
}

// simpleCache простая реализация кэша на основе map.
type simpleCache struct {
	data     map[string]any
	capacity int
}

func (c *simpleCache) Get(key string) (any, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c *simpleCache) Put(key string, value any) {
	// Если достигнут лимит, удаляем старые записи (простая стратегия)
	if c.capacity > 0 && len(c.data) >= c.capacity {
		// Удаляем первую запись (FIFO)
		for k := range c.data {
			delete(c.data, k)
			break
		}
	}
	c.data[key] = value
}

func (c *simpleCache) Clear() {
	for k := range c.data {
		delete(c.data, k)
	}
}

func (c *simpleCache) Size() int {
	return len(c.data)
}

func (c *simpleCache) Capacity() int {
	return c.capacity
}

func (c *simpleCache) Entries() []CacheEntry[string, any] {
	entries := make([]CacheEntry[string, any], 0, len(c.data))
	for k, v := range c.data {
		entries = append(entries, CacheEntry[string, any]{Key: k, Value: v})
	}
	return entries
}
