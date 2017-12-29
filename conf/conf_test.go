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
	assert.Equal(t, conf.Port, "8000")
	assert.Equal(t, conf.DatabaseURL, "teststring")
	assert.Equal(t, conf.RedisURL, "teststring")
	// Defaults
	assert.Equal(t, conf.MaxCacheUses, 10)
	assert.Equal(t, conf.UrlCacheExpiration, 60)
}

func TestInvalidConf(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("REDIS_URL")
	_, err := New()
	assert.NotNil(t, err)
}
