package jexl

// Cache представляет интерфейс кэша выражений.
type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
	Clear()
	Size() int
}

// CacheFactory создаёт кэш для заданного размера.
type CacheFactory func(size int) Cache[string, any]

// DefaultCacheFactory создаёт простой кэш на основе map.
var DefaultCacheFactory CacheFactory = func(size int) Cache[string, any] {
	return &simpleCache{
		data: make(map[string]any, size),
	}
}

// simpleCache простая реализация кэша на основе map.
type simpleCache struct {
	data map[string]any
}

func (c *simpleCache) Get(key string) (any, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c *simpleCache) Put(key string, value any) {
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
