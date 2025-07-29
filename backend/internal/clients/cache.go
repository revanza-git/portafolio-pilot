package clients

import (
	"sync"
	"time"
)

// CacheItem represents a cached item with expiration
type CacheItem struct {
	Quote     *Quote
	ExpiresAt time.Time
}

// MemoryCache implements Cache interface using in-memory storage
type MemoryCache struct {
	items map[string]*CacheItem
	mutex sync.RWMutex
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*CacheItem),
	}
	
	// Start cleanup routine
	go cache.cleanupExpired()
	
	return cache
}

// Get retrieves a quote from cache
func (c *MemoryCache) Get(key string) (*Quote, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.ExpiresAt) {
		// Remove expired item (do this in a separate goroutine to avoid deadlock)
		go func() {
			c.mutex.Lock()
			delete(c.items, key)
			c.mutex.Unlock()
		}()
		return nil, false
	}

	return item.Quote, true
}

// Set stores a quote in cache with TTL
func (c *MemoryCache) Set(key string, quote *Quote, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = &CacheItem{
		Quote:     quote,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a quote from cache
func (c *MemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.items, key)
}

// Clear removes all items from cache
func (c *MemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items = make(map[string]*CacheItem)
}

// Size returns the number of items in cache
func (c *MemoryCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	return len(c.items)
}

// cleanupExpired removes expired items from cache
func (c *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}