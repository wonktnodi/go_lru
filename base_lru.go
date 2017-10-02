package go_lru

import (
    "container/list"
    "errors"
    "time"
)

const (
    // For use with functions that take an expiration time.
    NoExpiration time.Duration = -1
    // For use with functions that take an expiration time. Equivalent to
    // passing in the same expiration duration as was given to New() or
    // NewFrom() when the cache was created (e.g. 5 minutes.)
    DefaultExpiration time.Duration = 0
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback func(key string, value interface{})

// LRU implements a non-thread safe fixed size LRU cache
type BASELRU struct {
    size              int
    evictList         *list.List
    items             map[string]*list.Element
    onEvict           EvictCallback
    defaultExpiration time.Duration
}

// entry is used to hold a value in the evictList
type entry struct {
    key        string
    value      interface{}
    Expiration int64
}

// Returns true if the item has expired.
func (item entry) Expired() bool {
    if item.Expiration == 0 {
        return false
    }
    return time.Now().UnixNano() > item.Expiration
}

// NewLRU constructs an LRU of the given size
func NewBaseLRU(size int, onEvict EvictCallback, defaultExpiration time.Duration) (*BASELRU, error) {
    if size <= 0 {
        return nil, errors.New("Must provide a positive size")
    }
    c := &BASELRU{
        size:      size,
        evictList: list.New(),
        items:     make(map[string]*list.Element),
        onEvict:   onEvict,
        defaultExpiration: defaultExpiration,
    }
    return c, nil
}

// Purge is used to completely clear the cache
func (c *BASELRU) Purge() {
    for k, v := range c.items {
        if c.onEvict != nil {
            c.onEvict(k, v.Value.(*entry).value)
        }
        delete(c.items, k)
    }
    c.evictList.Init()
}

func (c *BASELRU) Add(key string, value interface{}) bool {
    return c.AddWithExpire(key, value, NoExpiration)
}

func (c *BASELRU) AddWithTimeout(key string, value interface{}, timeout int64) bool {
    // Check for existing item
    if ent, ok := c.items[key]; ok {
        c.evictList.MoveToFront(ent)
        ent.Value.(*entry).value = value
        return false
    }

    // Add new item
    ent := &entry{key, value, timeout}
    entry := c.evictList.PushFront(ent)
    c.items[key] = entry

    evict := c.evictList.Len() > c.size
    // Verify size not exceeded
    if evict {
        c.removeOldest()
    }
    return evict
}


// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *BASELRU) AddWithExpire(key string, value interface{}, d time.Duration) bool {
    var e int64
    if d == DefaultExpiration {
        d = c.defaultExpiration
    }
    if d > 0 {
        e = time.Now().Add(d).UnixNano()
    }
    return c.AddWithTimeout(key, value, e)
}

// Get looks up a key's value from the cache.
func (c *BASELRU) Get(key string) (value interface{}, ok bool) {
    ent, ok := c.items[key]
    if !ok {
        return nil, false
    }

    item := ent.Value.(*entry)

    if item.Expiration > 0 {
        if time.Now().UnixNano() > item.Expiration {
            return nil, false
        }
    }
    c.evictList.MoveToFront(ent)

    return ent.Value.(*entry).value, true
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *BASELRU) Contains(key string) (ok bool) {
    item, ok := c.items[key]
    if ok && item.Value.(*entry).Expired() {
        return false
    }
    return ok
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *BASELRU) Peek(key string) (value interface{}, ok bool) {
    if ent, ok := c.items[key]; ok {
        return ent.Value.(*entry).value, true
    }
    return nil, ok
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *BASELRU) PeekWithExpire(key string) (value interface{}, ok bool, ts int64) {
    if ent, ok := c.items[key]; ok {
        val := ent.Value.(*entry)
        return val.value, true, val.Expiration
    }
    return nil, ok, 0
}


// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *BASELRU) Remove(key string) bool {
    if ent, ok := c.items[key]; ok {
        c.removeElement(ent)
        return true
    }
    return false
}

// RemoveOldest removes the oldest item from the cache.
func (c *BASELRU) RemoveOldest() (string, interface{}, bool) {
    ent := c.evictList.Back()
    if ent != nil {
        c.removeElement(ent)
        kv := ent.Value.(*entry)
        return kv.key, kv.value, true
    }
    return "", nil, false
}

// GetOldest returns the oldest entry
func (c *BASELRU) GetOldest() (string, interface{}, bool) {
    ent := c.evictList.Back()
    if ent != nil {
        kv := ent.Value.(*entry)
        return kv.key, kv.value, true
    }
    return "", nil, false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *BASELRU) Keys() []string {
    keys := make([]string, len(c.items))
    i := 0
    for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
        keys[i] = ent.Value.(*entry).key
        i++
    }
    return keys
}

// Len returns the number of items in the cache.
func (c *BASELRU) Len() int {
    return c.evictList.Len()
}

// removeOldest removes the oldest item from the cache.
func (c *BASELRU) removeOldest() {
    ent := c.evictList.Back()
    if ent != nil {
        c.removeElement(ent)
    }
}

// removeElement is used to remove a given list element from the cache
func (c *BASELRU) removeElement(e *list.Element) {
    c.evictList.Remove(e)
    kv := e.Value.(*entry)
    delete(c.items, kv.key)
    if c.onEvict != nil {
        c.onEvict(kv.key, kv.value)
    }
}
