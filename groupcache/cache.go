package groupcache

import (
	"geecache-learning/groupcache/lru"
	"sync"
)

// Cache 每一个cache， 因为先后读是有影响的，所以要加锁
type Cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	CacheBytes int64
}

func (c *Cache) Add(key string, value ByteView) {
	//通杀只有一个add和get， 保证了add和get的原子性
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.CacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *Cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
