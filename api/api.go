package api

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mkrysiak/cyclops-go/conf"
	"github.com/mkrysiak/cyclops-go/hash"
	"github.com/urfave/negroni"

	"github.com/mkrysiak/cyclops-go/models"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Api struct {
	cfg                *conf.Config
	cache              *models.Cache
	requestStorage     *models.RequestStorage
	projects           *models.SentryProjects
	urlCacheExpiration time.Duration
	ignoredItems       uint64
	processedItems     uint64
}

func NewApiRouter(cfg *conf.Config, cache *models.Cache, requestStorage *models.RequestStorage,
	projects *models.SentryProjects) *negroni.Negroni {
	api := &Api{
		cfg:                cfg,
		requestStorage:     requestStorage,
		projects:           projects,
		cache:              cache,
		urlCacheExpiration: time.Duration(cfg.UrlCacheExpiration) * time.Second,
		ignoredItems:       0,
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/{projectId:[0-9]+}/store/", api.apiHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/healthcheck", api.healthcheckHandler).Methods("GET")
	//TODO: Restrict access to /stats.  Ideally, it should not be public.
	r.HandleFunc("/stats", api.statsHandler).Methods("GET")

	// Middleware
	n := negroni.New()
	n.Use(negroni.HandlerFunc(api.LoggingMiddleware))
	n.UseHandler(r)

	return n
}

func (a *Api) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (a *Api) apiHandler(w http.ResponseWriter, r *http.Request) {

	a.addCorsHeaders(w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	projectId, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	xSentryAuth := NewXSentryAuth(r)
	if !a.projects.IsValidProjectAndPublicKey(projectId, xSentryAuth.sentry_key) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The request body could be plain JSON or Base64 encoded JSON
	bodyBytes, err := getRequestBody(r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Calculate a hash that identifies a unique body, and use it as a cache key in Redis
	exceptionHash, err := hash.HashForGrouping(bodyBytes)
	if err != nil {
		log.Errorf("Unable to calculate a hash for the request: %s", err)
	}

	var cacheKey bytes.Buffer
	cacheKey.WriteString(vars["projectId"])
	cacheKey.WriteString(exceptionHash)
	log.Debugf("Cache Key: %s", cacheKey.String())

	var originUrl bytes.Buffer
	originUrl.WriteString(a.cfg.SentryURL)
	originUrl.WriteString(r.RequestURI)
	log.Debugf("Origin URL: %s", originUrl.String())

	count := a.validateCache(cacheKey.String())
	if count > a.cfg.MaxCacheUses {
		w.Header().Set("X-CYCLOPS-CACHE-COUNT", strconv.FormatInt(count, 10))
		w.Header().Set("X-CYCLOPS-STATUS", "IGNORED")
		atomic.AddUint64(&a.ignoredItems, 1)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("X-CYCLOPS-CACHE-COUNT", strconv.FormatInt(count, 10))
	w.Header().Set("X-CYCLOPS-STATUS", "PROCESSED")
	atomic.AddUint64(&a.processedItems, 1)

	a.processRequest(r, projectId, originUrl.String(), bodyBytes)

	w.WriteHeader(http.StatusNoContent)

}

func (a *Api) statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Processed Items: "))
	w.Write([]byte(strconv.FormatUint(a.processedItems, 10)))
	w.Write([]byte("\n"))
	w.Write([]byte("Ignored Items: "))
	w.Write([]byte(strconv.FormatUint(a.ignoredItems, 10)))
}

func (a *Api) validateCache(cacheKey string) int64 {
	var count int64
	if a.urlCacheExpiration > 0 {
		count, _ = a.cache.Incr(cacheKey)
		if count == 1 {
			a.cache.Expire(cacheKey, a.urlCacheExpiration)
		}
	}
	return count
}

func (a *Api) processRequest(r *http.Request, projectId int, originUrl string, body []byte) {

	// Headers is a map[string][]string

	m := &models.Message{
		ProjectId:     projectId,
		RequestMethod: r.Method,
		Headers:       r.Header,
		OriginUrl:     originUrl,
		RequestBody:   body,
	}

	a.requestStorage.Put(projectId, m)
}

func (a *Api) addCorsHeaders(rw http.ResponseWriter, r *http.Request) {
	origin, validOrigin := a.isValidOrigin(r.Header.Get("Origin"))
	if validOrigin {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Headers", "X-Sentry-Auth, X-Requested-With, Origin, Accept, Content-Type, Authentication")
		rw.Header().Set("Access-Control-Expose-Headers",
			"Cache-Control,Content-Encoding,Content-Length,Content-Type,Date,ETag,Expires,Pragma,Server,Vary,X-CYCLOPS-CACHE-COUNT,X-CYCLOPS-STATUS")
		rw.Header().Set("Access-Control-Max-Age", "86400")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	}
}

func (a *Api) isValidOrigin(origin string) (string, bool) {
	if a.cfg.CyclopsAllowOrigin == "" {
		return "", false
	}
	if a.cfg.CyclopsAllowOrigin == "*" {
		return "*", true
	}
	u, err := url.Parse(origin)
	if err != nil {
		log.Error(err)
	}
	if strings.HasSuffix(u.Host, a.cfg.CyclopsAllowOrigin) {
		return origin, true
	}
	return "", false
}

func getRequestBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return body, err
	}

	// This seems to be the best way to test if a byte array is encoded.
	// If it fails, it's not encoded.
	b64decodedBody, err := base64.StdEncoding.DecodeString(string(body))
	if err == nil {
		body = b64decodedBody
	}
	return body, nil
}
