package models

import (
	"github.com/go-redis/redis"
	"github.com/mkrysiak/cyclops-go/conf"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	Client *redis.Client
}

func NewRedisClient(cfg *conf.Config) *Redis {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Error(err)
	}

	log.WithFields(log.Fields{
		"Address": opt.Addr,
		"DB":      opt.DB,
	}).Info("Redis Connection:")

	return &Redis{
		Client: redis.NewClient(opt),
	}
}

// For use by testing with miniredis
// TODO: Try mocking Redis instead?
// https://github.com/go-redis/redis/issues/332
func NewRedisByAddr(add string) *Redis {

	return &Redis{
		Client: redis.NewClient(&redis.Options{Addr: add}),
	}
}

func (r *Redis) Shutdown() {
	log.Info("Closing Redis Connection")
	r.Client.Close()
}
