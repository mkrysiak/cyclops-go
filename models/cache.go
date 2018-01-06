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
	_, err := c.Client.Expire(key, expiration).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"Key":   key,
			"Error": err,
		}).Info("Unable to set expiration on key")
	}
}

// func (c *Cache) Set(key string, expiration time.Duration) error {
// 	// When running multiple instances of cyclops-go, lock during the set operation so two different
// 	// instances don't attempt to set the key to 0, which would reset the counter, losing count
// 	lock, err := lock.Obtain(c.Client, "cyclops:lock:", &lock.Options{LockTimeout: 300 * time.Second})
// 	if err != nil {
// 		// It's normal to not obtain a lock, so don't log it as an error
// 		// log.Error(err)
// 		return err
// 	}
// 	defer lock.Unlock()
// 	err = c.Client.Set(key, 0, expiration).Err()
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return err
// }

func (c *Cache) Ttl(key string) time.Duration {
	return c.Client.TTL(key).Val()
}

func (c *Cache) Flushdb() {
	c.Client.FlushDB()
}
