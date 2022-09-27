package lru

import (
	"container/list"
)

// Cache LRU 缓存，并发访问不安全
type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	nBytes    int64                         // 已经使用了的内存
	ll        *list.List                    // Go语言标准库实现的双向链表
	cache     map[string]*list.Element      // value为双向链表对应的节点指针
	OnEvicted func(key string, value Value) // 当一条记录被删除后执行的回调函数，可以为空
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// entry 双向链表的节点
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// Get 查找指定的节点
func (cache *Cache) Get(key string) (Value, bool) {
	// ok表示缓存中存在指定的数据
	if ele, ok := cache.cache[key]; ok {

		// 将该节点移动到链表队头
		cache.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)

		// 返回查询到的数据
		return kv.value, true
	}

	return nil, false
}

// RemoveOldest 缓存淘汰，移除最近最少使用的数据
func (cache *Cache) RemoveOldest() {
	ele := cache.ll.Back()
	if ele != nil {
		cache.ll.Remove(ele)
		kv := ele.Value.(*entry)

		// 删除map中的kv键值对
		delete(cache.cache, kv.key)

		// 修改缓存占用内存
		cache.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())

		// 如果存在回调函数，则执行
		if cache.OnEvicted != nil {
			cache.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 添加缓存
func (cache *Cache) Add(key string, value Value) {
	// 如果key存在，表示更新节点的值
	if ele, ok := cache.cache[key]; ok {
		kv := ele.Value.(*entry)

		// 更新缓存后缓存的大小
		var tempBytes = int64(len(kv.key)) + int64(kv.value.Len()) - int64(value.Len()) - int64(len(key))
		// 如果该大小超过maxBytes，则移除最近最少使用的缓存
		for cache.maxBytes != 0 && cache.maxBytes < tempBytes {
			cache.RemoveOldest()
		}

		cache.ll.MoveToFront(ele)
		cache.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		for cache.maxBytes != 0 && cache.maxBytes < cache.nBytes+int64(value.Len())+int64(len(key)) {
			cache.RemoveOldest()
		}

		// 新增节点
		ele = cache.ll.PushFront(&entry{key, value})
		cache.cache[key] = ele
		cache.nBytes += int64(len(key)) + int64(value.Len())
	}

	// 如果添加缓存后缓存容量超过 maxBytes，则移除最少访问的节点
	//for cache.maxBytes != 0 && cache.maxBytes < cache.nBytes {
	//	cache.RemoveOldest()
	//}
}

// Len 返回节点的个数
func (cache *Cache) Len() int {
	return cache.ll.Len()
}
