package api

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (a *Api) OptionsHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	origin := r.Header.Get("Origin")
	u, err := url.Parse(origin)
	if err != nil {
		log.Error(err)
	}
	if (a.cfg.CorsExternalDNS == "" || !strings.HasSuffix(u.Host, a.cfg.CorsExternalDNS)) &&
		strings.HasSuffix(u.Host, a.cfg.CorsOriginSuffix) {

		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Expose-Headers",
			"Cache-Control,Content-Encoding,Content-Length,Content-Type,Date,ETag,Expires,Pragma,Server,Vary,X-CYCLOPS-CACHE-COUNT,X-CYCLOPS-STATUS")

		if r.Method == "OPTIONS" {
			rw.Header().Set("Access-Control-Max-Age", "86400")
			rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		}
	}
	// do some stuff before
	next(rw, r)
	// do some stuff after

}

func (a *Api) LoggingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// do some stuff before
	next(rw, r)
	// do some stuff after

	var requestLogLine bytes.Buffer
	requestLogLine.WriteString(r.RemoteAddr)
	requestLogLine.WriteString(" ")
	requestLogLine.WriteString(r.Method)
	requestLogLine.WriteString(" ")
	requestLogLine.WriteString(r.URL.Path)
	log.Info(requestLogLine.String())
}
