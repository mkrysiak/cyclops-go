package api

import (
	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"
)

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
