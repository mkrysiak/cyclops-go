package models

import (
	"bytes"
	"time"

	"github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		Client: client,
	}
}

func (c *Cache) Get(key string) (int64, error) {
	value, err := c.Client.Get(key).Int64()
	if err != nil {
		log.Error(err)
	}
	return value, err
}

func (c *Cache) Incr(key string) (int64, error) {
	value, err := c.Client.Incr(key).Result()
	if err != nil {
		log.Error(err)
	}
	return value, err
}

func (c *Cache) Set(key string, expiration time.Duration) error {
	var lockKey bytes.Buffer
	lockKey.WriteString("cycleops:lock:")
	lockKey.WriteString(key)

	lock, err := lock.Obtain(c.Client, lockKey.String(), &lock.Options{LockTimeout: 300 * time.Second})
	if err != nil {
		log.Error(err)
		return err
	}
	defer lock.Unlock()
	log.Info(expiration)
	err = c.Client.Set(key, 0, expiration).Err()
	if err != nil {
		log.Error(err)
	}
	return err
}

func (c *Cache) Ttl(key string) time.Duration {
	return c.Client.TTL(key).Val()
}

func (c *Cache) Flushdb() {
	c.Client.FlushDB()
}
