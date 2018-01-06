package models

import (
	"time"

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

func (c *Cache) Expire(key string, expiration time.Duration) {
	err := c.Client.Expire(key, expiration).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"Key":   key,
			"Error": err,
		}).Info("Unable to set expiration on key.  Messages with this key may never get processed!")
	}
}

func (c *Cache) Ttl(key string) time.Duration {
	return c.Client.TTL(key).Val()
}

func (c *Cache) Flushdb() {
	c.Client.FlushDB()
}
