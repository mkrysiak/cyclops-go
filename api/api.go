package api

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/mkrysiak/cyclops-go/conf"
	"github.com/mkrysiak/cyclops-go/hash"

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
	projects *models.SentryProjects) *mux.Router {
	api := &Api{
		cfg:                cfg,
		requestStorage:     requestStorage,
		projects:           projects,
		cache:              cache,
		urlCacheExpiration: time.Duration(cfg.UrlCacheExpiration) * time.Second,
		ignoredItems:       0,
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/{projectId:[0-9]+}/store/", api.apiHandler).Methods("POST")
	r.HandleFunc("/healthcheck", api.healthcheckHandler).Methods("GET")
	//TODO: Restrict access to /stats.  It should not be public.
	r.HandleFunc("/stats", api.statsHandler).Methods("GET")
	return r
}

func (a *Api) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
	log.Info(r.RemoteAddr + " " + r.Method + " " + r.URL.Path)
}

func (a *Api) apiHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

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

	//TODO: Make URL configurable
	var originUrl bytes.Buffer
	originUrl.WriteString("http://localhost:2222")
	originUrl.WriteString(r.RequestURI)
	log.Debugf("Origin URL: %s", originUrl.String())

	// TODO: It's bad practice to return headers that can identify the product that's in use if
	// this proxy is externally exposed.
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

func (a *Api) validateCache(url string) int64 {
	var count int64
	if a.urlCacheExpiration > 0 {
		count, _ = a.cache.Get(url)
		if count == 0 {
			a.cache.Set(url, a.urlCacheExpiration)
		}
		count, _ = a.cache.Incr(url)
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

func logRequest(r *http.Request) {
	var requestLogLine bytes.Buffer
	requestLogLine.WriteString(r.RemoteAddr)
	requestLogLine.WriteString(" ")
	requestLogLine.WriteString(r.Method)
	requestLogLine.WriteString(" ")
	requestLogLine.WriteString(r.URL.Path)
	log.Info(requestLogLine.String())
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
