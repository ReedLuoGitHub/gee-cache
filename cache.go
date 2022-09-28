package gee_cache

import (
	"gee-cache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

// add 添加缓存，通过互斥锁保证并发性
func (cache *cache) add(key string, value ByteView) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// 当lru为空值，初始化lru
	if cache.lru == nil {
		cache.lru = lru.New(cache.cacheBytes, nil)
	}

	cache.lru.Add(key, value)
}

// get ...
func (cache *cache) get(key string) (ByteView, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.lru == nil {
		return ByteView{}, false
	}

	if v, ok := cache.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return ByteView{}, false
}
