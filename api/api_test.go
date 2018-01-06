package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/gavv/httpexpect"
	"github.com/mkrysiak/cyclops-go/conf"
	"github.com/mkrysiak/cyclops-go/models"
)

func TestAPI(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("PORT", "23343")
	os.Setenv("DATABASE_URL", "postgres://user:pw@localhost/cyclops")
	os.Setenv("REDIS_URL", "redis://:pw@localhost:23234/1")
	os.Setenv("SENTRY_URL", "http://localhost:2222")
	os.Setenv("MAX_CACHE_USES", "5")
	os.Setenv("URL_CACHE_EXPIRATION", "60")
	os.Setenv("CYCLOPS_ALLOW_ORIGIN", "localhost:23343")

	cfg, err := conf.New()
	if err != nil {
		t.Log(err)
	}

	s, err := miniredis.Run()
	if err != nil {
		t.Log(err)
	}
	redis := models.NewRedisByAddr(s.Addr())
	cache := models.NewCache(redis.Client)
	requestStorage := models.NewRequestStorage(redis.Client)
	sentryProjects := models.NewSentryProjects(cfg.GetDatabaseSchemeAndUrl())

	teardown := setup(t, sentryProjects.Db)
	defer teardown(t, sentryProjects.Db)

	sentryProjects.UpdateProjects()

	server := httptest.NewServer(NewApiRouter(cfg, cache, requestStorage, sentryProjects))
	defer server.Close()

	exception := map[string]interface{}{
		"message": "test",
		"culprit": "test",
	}

	e := httpexpect.New(t, server.URL)

	t.Run("GET /healthcheck", func(t *testing.T) {
		e.GET("/healthcheck").Expect().Status(http.StatusOK).
			Body().Equal("OK")
	})

	t.Run("POST /api", func(t *testing.T) {
		redis.Client.Del("405a671c66aefea124cc08b76ea6d30bb")
		e.POST("/api/4/store/").
			WithQuery("sentry_key", "42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
			Expect().Status(http.StatusNoContent)
	})

	t.Run("POST /api using X-Sentry-Auth", func(t *testing.T) {
		redis.Client.Del("405a671c66aefea124cc08b76ea6d30bb")
		e.POST("/api/4/store/").
			WithHeader("X-Sentry-Auth", "sentry_key=42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
			Expect().Status(http.StatusNoContent)
	})

	t.Run("POST /api 404", func(t *testing.T) {
		e.POST("/api/100/store/").
			WithQuery("sentry_key", "42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
			Expect().Status(http.StatusNotFound)
	})

	t.Run("MaxCacheUses Limit", func(t *testing.T) {
		redis.Client.Del("405a671c66aefea124cc08b76ea6d30bb")
		// cfg.MaxCacheUses should equal 5
		for i := 0; i < 5; i++ {
			e.POST("/api/4/store/").
				WithQuery("sentry_key", "42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
				Expect().Status(http.StatusNoContent).Header("X-CYCLOPS-STATUS").Equal("PROCESSED")
		}
		e.POST("/api/4/store/").
			WithQuery("sentry_key", "42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
			Expect().Status(http.StatusNoContent).Header("X-CYCLOPS-STATUS").Equal("IGNORED")

		// TODO: This is a poor test.  Provide a way to reset the stats, and verify that
		// the counters are accurate
		e.GET("/stats").Expect().Status(http.StatusOK).Body().
			Match("Processed Items: [0-9]+\nIgnored Items: [0-9]+")
		// str := "Processed Items: " + strconv.Itoa(cfg.MaxCacheUses+1) + "\nIgnored Items: 1"
		//e.GET("/stats").Expect().Status(http.StatusOK).Body().Contains(str)
	})

	t.Run("Increment Cache", func(t *testing.T) {
		redis.Client.Del("405a671c66aefea124cc08b76ea6d30bb")
		e.POST("/api/4/store/").
			WithQuery("sentry_key", "42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
			Expect().Status(http.StatusNoContent)
		cacheValue, _ := cache.Get("405a671c66aefea124cc08b76ea6d30bb")
		assert.Equal(t, int64(1), cacheValue)
	})

	t.Run("CORS Headers Exist", func(t *testing.T) {
		e.POST("/api/4/store/").
			WithQuery("sentry_key", "42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
			WithHeader("Origin", "http://test.localhost:23343").
			Expect().Status(http.StatusNoContent).Header("Access-Control-Allow-Origin").Equal("http://test.localhost:23343")
	})

	t.Run("CORS Headers Missing", func(t *testing.T) {
		e.POST("/api/4/store/").
			WithQuery("sentry_key", "42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
			WithHeader("Origin", "http://badorigin").
			Expect().Status(http.StatusNoContent).Header("Access-Control-Allow-Origin").Equal("")
	})

	t.Run("CORS OPTIONS Request", func(t *testing.T) {
		e.OPTIONS("/api/4/store/").
			WithQuery("sentry_key", "42aa6019f602d77313ec553625ecb01a").WithJSON(exception).
			WithHeader("Origin", "http://localhost:23343").
			Expect().Status(http.StatusNoContent).Header("Access-Control-Allow-Origin").Equal("http://localhost:23343")
	})

}

func setup(t *testing.T, db *sqlx.DB) func(t *testing.T, db *sqlx.DB) {

	Up(t, db)

	return func(t *testing.T, db *sqlx.DB) {
		t.Log("teardown")
		Down(t, db)
	}
}

func Up(t *testing.T, db *sqlx.DB) error {

	createTable := `CREATE TABLE sentry_projectkey (
		project_id INTEGER PRIMARY_KEY,
		public_key TEXT NOT NULL,
		secret_key TEXT NOT NULL
	);`

	db.MustExec(createTable)

	query := `INSERT INTO sentry_projectkey (project_id, public_key, secret_key)
	VALUES
	(1, '54c232ae245242e2c224a938b8ffda41', '256d07fd72870e5f4e91ed2f7f2007cd'),
	(2, '151268aad4d2ed5a8c0eded096663b85', '9026d41a8ac820ace03301bd6d47ee66'),
	(3, 'f5c3d55205ba8fff225b5b4d50451ee6', '83c356c0c25d2c8855b1cb96edf11c4b'),
	(4, '42aa6019f602d77313ec553625ecb01a', 'bc6c06e5b0f68e087bcd34e44523a55d'),
	(5, 'dee449e708ece38bb7697b8bbc7f8387', 'dae973d82aa5c183c8da4ec6283a14cf');`

	_, err := db.Exec(query)
	if err != nil {
		t.Log(err)
		return err
	}

	return nil
}

func Down(t *testing.T, db *sqlx.DB) error {
	_, err := db.Exec("DROP TABLE sentry_projectkey")
	return err
}
