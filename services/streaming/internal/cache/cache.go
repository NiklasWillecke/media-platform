package cache

import (
	"container/list"
	"sync"
)

type item struct {
	data []byte
	key  string
}

type LRUCache struct {
	Mu       sync.Mutex
	items    map[string]*list.Element
	maxSize  int // in Bytes (z. B. 100 * 1024 * 1024 für 100 MB)
	currSize int
	order    *list.List
}

func NewLRUCache(size int) *LRUCache {
	return &LRUCache{
		items:    make(map[string]*list.Element),
		order:    list.New(),
		maxSize:  size, //In MB
		currSize: 0,
	}
}

func (c *LRUCache) Set(key string, value []byte) {

	c.Mu.Lock()
	defer c.Mu.Unlock()

	newItem := &item{
		data: value,
		key:  key,
	}

	c.currSize += len(value)

	newElement := c.order.PushFront(newItem)
	c.items[key] = newElement

	for c.maxSize < c.currSize && c.order.Len() > 0 {
		last := c.order.Back()

		lastItem := last.Value.(*item)

		delete(c.items, lastItem.key)
		c.currSize -= len(lastItem.data)

		c.order.Remove(last)
	}

}

func (c *LRUCache) Get(key string) ([]byte, bool) {
	if n, ok := c.items[key]; ok {

		c.Mu.Lock()
		c.order.MoveToFront(n)
		c.Mu.Unlock()

		return n.Value.(*item).data, true
	}
	return nil, false

}
