package pokecache

import (
    "sync"
    "time"
)

type cacheEntry struct {
    createdAt time.Time
    value []byte
}

type Cache struct {
    cache map[string]cacheEntry
    mu sync.RWMutex 
}

// Create a new cache that periodically reaps stale entries at the given interval. Duration <= 0
// means no reaping.
func NewCache(interval time.Duration) *Cache {
    cache := Cache { cache: make(map[string]cacheEntry, 0) }

    if interval > 0 {
        go cache.reapLoop(interval)
    }

    return &cache
}

func (c *Cache) Add(key string, value []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()

    now := time.Now()
    c.cache[key] = cacheEntry { createdAt: now, value: value }
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    entry, exists := c.cache[key]
    return entry.value, exists
}

func (c *Cache) reapLoop(interval time.Duration) {
    for tick := range time.Tick(interval) {
        c.mu.Lock()

        for key, entry := range c.cache {
            if tick.Sub(entry.createdAt).Milliseconds() >= interval.Milliseconds() {
                delete(c.cache, key)
            }
        }

        c.mu.Unlock()
    }
}
