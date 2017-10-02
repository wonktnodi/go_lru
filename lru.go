// This package provides a simple LRU cache. It is based on the
// LRU implementation in groupcache:
// https://github.com/golang/groupcache/tree/master/lru
package go_lru

import (
	"sync"
	"time"
)

// Cache is a thread-safe fixed size LRU cache.
type Cache struct {
	lru  *BASELRU
	lock sync.RWMutex
}

// New creates an LRU of the given size
func New(size int, defaultExpiration time.Duration) (*Cache, error) {
	return NewWithEvict(size, defaultExpiration,nil)
}

// NewWithEvict constructs a fixed size cache with the given eviction
// callback.
func NewWithEvict(size int, defaultExpiration time.Duration, onEvicted func(key string, value interface{})) (*Cache, error) {
	lru, err := NewBaseLRU(size, EvictCallback(onEvicted), defaultExpiration)
	if err != nil {
		return nil, err
	}
	c := &Cache{
		lru: lru,
	}
	return c, nil
}

// Purge is used to completely clear the cache
func (c *Cache) Purge() {
	c.lock.Lock()
	c.lru.Purge()
	c.lock.Unlock()
}

// Add adds a value to the cache with expiration.  Returns true if an eviction occurred.
func (c *Cache) AddWithExpire(key string, value interface{}, d time.Duration) bool {
    c.lock.Lock()
    defer c.lock.Unlock()
    return c.lru.AddWithExpire(key, value, d)
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *Cache) Add(key string, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Add(key, value)
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Get(key)
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *Cache) Contains(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Contains(key)
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *Cache) Peek(key string) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Peek(key)
}

// ContainsOrAdd checks if a key is in the cache  without updating the
// recent-ness or deleting it for being stale,  and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache) ContainsOrAdd(key string, value interface{}) (ok, evict bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru.Contains(key) {
		return true, false
	} else {
		evict := c.lru.Add(key, value)
		return false, evict
	}
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) {
	c.lock.Lock()
	c.lru.Remove(key)
	c.lock.Unlock()
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache) RemoveOldest() {
	c.lock.Lock()
	c.lru.RemoveOldest()
	c.lock.Unlock()
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *Cache) Keys() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Keys()
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Len()
}
