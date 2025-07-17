package dal

import (
	"sync"
	"time"
)

// CacheItem represents an item stored in the cache with its expiration time.
type CacheItem struct {
	Value      interface{}
	Expiration int64 // Unix timestamp
}

// Cache is a simple in-memory key-value store with optional TTL.
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]CacheItem),
	}
	// Start a goroutine to clean up expired items periodically
	go c.cleanupLoop()
	return c
}

// Set adds an item to the cache with a given key and TTL.
// ttl is in seconds. If ttl is 0, the item never expires.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := int64(0)
	if ttl > 0 {
		expiration = time.Now().Add(ttl).Unix()
	}

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves an item from the cache. Returns nil if not found or expired.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		// Item expired, remove it
		delete(c.items, key)
		return nil, false
	}

	return item.Value, true
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

// cleanupLoop periodically removes expired items from the cache.
func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute) // Check every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, item := range c.items {
			if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// SetMany adds multiple items to the cache.
func (c *Cache) SetMany(items map[string]interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := int64(0)
	if ttl > 0 {
		expiration = time.Now().Add(ttl).Unix()
	}

	for key, value := range items {
		c.items[key] = CacheItem{
			Value:      value,
			Expiration: expiration,
		}
	}
}
