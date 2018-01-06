package api

import (
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

// https://docs.sentry.io/clientdev/overview/#authentication
type XSentryAuth struct {
	// sentry_version string
	// sentry_client  string
	// sentry_timestamp string
	sentry_key    string
	sentry_secret string
}

// The sentry public key can be sent in two ways, using the "X-Sentry-Auth"
// header, or as a query argument.  The header takes precendence.
func NewXSentryAuth(r *http.Request) *XSentryAuth {
	var xSentryAuth XSentryAuth
	xSentryAuth.sentry_key = r.URL.Query().Get("sentry_key")
	headerValues := header.ParseList(r.Header, "X-Sentry-Auth")
	for _, v := range headerValues {
		sp := strings.SplitN(v, "=", 2)
		if len(sp) != 2 {
			return &xSentryAuth
		}
		switch sp[0] {
		case "sentry_key":
			xSentryAuth.sentry_key = sp[1]
		case "sentry_secret":
			xSentryAuth.sentry_secret = sp[1]
		}
	}
	return &xSentryAuth
}
