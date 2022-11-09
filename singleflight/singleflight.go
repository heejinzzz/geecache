package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mutex sync.Mutex // protect m
	m     map[string]*call
}

func NewGroup() *Group {
	return &Group{
		m: make(map[string]*call),
	}
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mutex.Lock()
	if c, ok := g.m[key]; ok {
		g.mutex.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := &call{}
	c.wg.Add(1)
	g.m[key] = c
	g.mutex.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mutex.Lock()
	delete(g.m, key)
	g.mutex.Unlock()

	return c.val, c.err
}
