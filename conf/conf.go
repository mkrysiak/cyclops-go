package conf

import (
	"net/url"
	"os"

	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Environment        string `default:"development" env:"ENVIRONMENT"`
	Port               string `required:"true" env:"PORT"`
	DatabaseURL        string `required:"true" env:"DATABASE_URL"`
	RedisURL           string `required:"true" env:"REDIS_URL"`
	SentryURL          string `required:"true" env:"SENTRY_URL"`
	MaxCacheUses       int64  `default:"10" env:"MAX_CACHE_USES"`
	UrlCacheExpiration int    `default:"60" env:"URL_CACHE_EXPIRATION"`
}

//TODO: Input validation
func New() (*Config, error) {
	config := new(Config)
	os.Setenv("CONFIGOR_ENV_PREFIX", "-")
	err := configor.New(&configor.Config{Debug: false}).Load(config)
	if err != nil {
		return nil, err
	}
	return config, err
}

func (c *Config) GetDatabaseSchemeAndUrl() (string, string) {

	if c.Environment == "test" {
		return "sqlite3", ":memory:"
	}

	u, err := url.Parse(c.DatabaseURL)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	return u.Scheme, c.DatabaseURL

}
