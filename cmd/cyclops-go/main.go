package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkrysiak/cyclops-go/api"
	"github.com/mkrysiak/cyclops-go/conf"
	"github.com/mkrysiak/cyclops-go/models"
	"github.com/mkrysiak/cyclops-go/tasks"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg, err := conf.New()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.WithFields(log.Fields{
		"Port":               cfg.Port,
		"DatabaseUrl":        cfg.DatabaseURL,
		"RedisUrl":           cfg.RedisURL,
		"UrlCacheExpiration": cfg.UrlCacheExpiration,
		"MaxCacheUses":       cfg.MaxCacheUses,
	}).Info("Loaded Config:")

	redis := models.NewRedisClient(cfg.RedisURL)
	cache := models.NewCache(redis.Client)
	requestStorage := models.NewRequestStorage(redis.Client)
	sentryProjects := models.NewSentryProjects(cfg.GetDatabaseSchemeAndUrl())

	tasks.StartBackgroundTaskRunners(sentryProjects, requestStorage)

	// Exit cleanly
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		sentryProjects.Shutdown()
		redis.Shutdown()
		os.Exit(0)
	}()

	http.ListenAndServe(":"+cfg.Port, api.NewApiRouter(cfg, cache, requestStorage, sentryProjects))

}
