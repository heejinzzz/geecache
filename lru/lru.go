package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	// Cache maximum capacity. When maxBytes is set to 0, it means that no limit is placed on the maximum Cache capacity
	maxBytes int64
	nBytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// entry is the data type of a list node
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, OnEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// Len returns the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

// Get looks up a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToBack(element)
		e := element.Value.(*entry)
		return e.value, true
	}
	return
}

// RemoveOldest removes the oldest entry
func (c *Cache) RemoveOldest() {
	target := c.ll.Front()
	if target != nil {
		c.ll.Remove(target)
		e := target.Value.(*entry)
		delete(c.cache, e.key)
		c.nBytes -= int64(len(e.key)) + int64(e.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(e.key, e.value)
		}
	}
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToBack(element)
		e := element.Value.(*entry)
		e.value = value
		c.nBytes += int64(value.Len()) - int64(e.value.Len())
	} else {
		e := &entry{key: key, value: value}
		element := c.ll.PushBack(e)
		c.cache[key] = element
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}
