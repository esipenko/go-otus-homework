package hw04lrucache

import (
	"sync"
)

type Key string

type mapValue struct {
	Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	mx       sync.Mutex
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	defer c.mx.Unlock()
	c.mx.Lock()

	myVal, ok := c.items[key]

	// Иначе не знаю как за O(1) удалить значение из мапы
	newValue := mapValue{key, value}

	if ok {
		myVal.Value = newValue

		c.queue.MoveToFront(myVal)
		return true
	}

	newLi := c.queue.PushFront(newValue)
	c.items[key] = newLi

	if c.queue.Len() > c.capacity {
		lastEl := c.queue.Back()
		c.queue.Remove(lastEl)
		delete(c.items, lastEl.Value.(mapValue).Key)
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	defer c.mx.Unlock()
	c.mx.Lock()

	val, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(val)
	return val.Value.(mapValue).value, true
}

func (c *lruCache) Clear() {
	defer c.mx.Unlock()
	c.mx.Lock()

	c.items = make(map[Key]*ListItem)
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem),
	}
}
