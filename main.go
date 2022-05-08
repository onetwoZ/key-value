package main

import (
	"container/list"
	"fmt"
	"sync"
)

type Data struct {
	data interface{}
	e    *list.Element
}

type Cache struct {
	mu       *sync.Mutex
	capacity int
	data     map[interface{}]*Data
	list     *list.List
}

func NewCache(capacity int) *Cache {
	return &Cache{
		mu:       new(sync.Mutex),
		capacity: capacity,
		data:     make(map[interface{}]*Data),
		list:     list.New(),
	}
}

func (c *Cache) dlt(e *list.Element) (i *list.Element) {
	c.list.Remove(e)
	delete(c.data, e.Value)
	return
}

func (c *Cache) upd(key interface{}) (i *list.Element) {
	i = c.list.PushBack(key)
	return
}

func (c *Cache) Put(key interface{}, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if node, ok := c.data[key]; ok {
		c.dlt(node.e)
		node.e = c.upd(key)
		node.data = data
		c.data[key] = node
	} else {
		var e *list.Element
		if len(c.data) < c.capacity {
			e = c.upd(key)
		} else {
			c.dlt(c.list.Front())
			e = c.upd(key)
		}
		node := &Data{
			data: data,
			e:    e,
		}
		c.data[key] = node
	}
}

func (c *Cache) Get(key interface{}) (exists bool, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if node, ok := c.data[key]; ok {
		data = node.data
		c.dlt(node.e)
		node.e = c.upd(key)
		c.data[key] = node
		exists = true
	}
	return
}
