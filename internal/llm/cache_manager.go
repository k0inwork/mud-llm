package llm

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type CacheManager struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		items: make(map[string]CacheItem),
	}
}

func (c *CacheManager) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(duration).UnixNano()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *CacheManager) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		// Item has expired
		delete(c.items, key)
		return nil, false
	}

	return item.Value, true
}

func (c *CacheManager) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}