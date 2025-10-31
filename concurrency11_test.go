package main

import (
	"container/list"
	"sync"
	"testing"
)

// 设计一个并发安全的LRU缓存，支持并发读写操作，并在测试中验证其正确性和性能。
func TestConcurrency11(t *testing.T) {
	lru := NewLRUCache(3)
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		lru.Put(1, 10)
		lru.Put(2, 20)
		lru.Put(3, 30)
	}()
	go func() {
		defer wg.Done()
		if val, ok := lru.Get(1); ok {
			t.Logf("Get key 1: %d", val)
		}
	}()
	go func() {
		defer wg.Done()
		lru.Put(4, 40) // 这将淘汰key 2
	}()
	go func() {
		defer wg.Done()
		if val, ok := lru.Get(2); ok {
			t.Logf("Get key 2: %d", val)
		} else {
			t.Log("Key 2 not found")
		}
	}()
	go func() {
		defer wg.Done()
		lru.Put(5, 50) // 这将淘汰key 3
	}()
	go func() {
		defer wg.Done()
		if val, ok := lru.Get(3); ok {
			t.Logf("Get key 3: %d", val)
		}
	}()
	wg.Wait()

}

type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	mu       sync.RWMutex
	lrulist  *list.List
}

func NewLRUCache(capacity int) *LRUCache {
	lruCache := &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		lrulist:  list.New(), // list的作用在于删除最久未使用的元素
	}
	return lruCache
}

type Entry struct {
	key   int
	value int
}

func (lru *LRUCache) Get(key int) (int, bool) {
	lru.mu.RLock()
	defer lru.mu.RUnlock()

	if elem, ok := lru.cache[key]; ok {
		lru.lrulist.MoveToFront(elem)
		return elem.Value.(*Entry).value, true
	}
	return 0, false
}
func (lru *LRUCache) Put(key int, value int) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if elem, ok := lru.cache[key]; ok {
		elem.Value.(*Entry).value = value
		lru.lrulist.MoveToFront(elem)
		return
	}
	// miss
	elem := lru.lrulist.PushFront(&Entry{key: key, value: value})
	lru.cache[key] = elem
	if lru.lrulist.Len() > lru.capacity {
		back := lru.lrulist.Back()
		lru.lrulist.Remove(back)
		delete(lru.cache, back.Value.(*Entry).key)
	}
}
