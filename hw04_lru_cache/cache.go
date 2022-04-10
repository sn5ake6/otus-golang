package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]
	if ok {
		if cachedItem, ok := item.Value.(*cacheItem); ok {
			cachedItem.value = value
		}
		c.queue.MoveToFront(item)

		return true
	}

	newCacheItem := &cacheItem{key, value}
	newQueueItem := c.queue.PushFront(newCacheItem)
	if c.queue.Len() > c.capacity {
		last := c.queue.Back()
		c.queue.Remove(last)
		if cachedItem, ok := last.Value.(*cacheItem); ok {
			delete(c.items, cachedItem.key)
		}
	}
	c.items[key] = newQueueItem

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		if cachedItem, ok := item.Value.(*cacheItem); ok {
			return cachedItem.value, true
		}
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
