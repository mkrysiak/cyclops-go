package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidConf(t *testing.T) {
	os.Setenv("PORT", "8000")
	os.Setenv("DATABASE_URL", "teststring")
	os.Setenv("REDIS_URL", "teststring")
	os.Setenv("SENTRY_URL", "http://localhost:2222")
	conf, err := New()
	if err != nil {
		t.Log(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "8000", conf.Port)
	assert.Equal(t, "teststring", conf.DatabaseURL)
	assert.Equal(t, "teststring", conf.RedisURL)
	// Defaults
	assert.Equal(t, int64(10), conf.MaxCacheUses)
	assert.Equal(t, 60, conf.UrlCacheExpiration)
}

func TestInvalidConf(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("REDIS_URL")
	_, err := New()
	assert.NotNil(t, err)
}
