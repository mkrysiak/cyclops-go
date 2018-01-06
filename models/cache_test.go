package models

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Log(err)
	}
	redis := NewRedisByAddr(s.Addr())
	cache := NewCache(redis.Client)

	teardown := setup(t, cache)
	defer teardown(t, cache)

	t.Run("Get Empty", func(t *testing.T) {
		v, err := cache.Get("testkey")
		assert.Equal(t, int64(0), v, "Getting a non-existant key")
		assert.NotNil(t, err, "Getting a non-existant key should error")
	})

	t.Run("SetAndGet", func(t *testing.T) {
		cache.Flushdb()
		cache.Incr("testkey")
		cache.Expire("testkey", 97*time.Second)
		v, err := cache.Get("testkey")
		if err != nil {
			t.Log(err)
		}
		ttl := cache.Ttl("testkey").Seconds()
		assert.Equal(t, int64(1), v, "Redis key should equal \"1\"")
		assert.Condition(t, func() bool {
			return ttl > 94 && ttl <= 97
		}, "TTL should only be slightly less than 97")
	})

	t.Run("Increment", func(t *testing.T) {
		cache.Incr("ikey")
		cache.Incr("ikey")
		v, err := cache.Incr("ikey")
		if err != nil {
			t.Log(err)
		}
		assert.Equal(t, int64(3), v, "Redis key should equal \"3\"")
	})
}

func setup(t *testing.T, c *Cache) func(t *testing.T, c *Cache) {
	c.Flushdb()
	return func(t *testing.T, c *Cache) {
		c.Flushdb()
	}
}
