package models

import (
	"bytes"
	"strconv"
	"sync"
)

type Counter struct {
	mu     sync.Mutex
	values map[string]uint64
}

func NewCounter() *Counter {
	return &Counter{
		values: make(map[string]uint64),
	}
}

func (c *Counter) Get(key string) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.values[key]
}

func (c *Counter) Incr(key string) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key]++
	return c.values[key]
}

func (c *Counter) GetCountersString() string {
	var buffer bytes.Buffer
	for k, v := range c.values {
		buffer.WriteString(k)
		buffer.WriteString(": ")
		buffer.WriteString(strconv.FormatUint(v, 10))
		buffer.WriteString("\n")
	}
	return buffer.String()
}
