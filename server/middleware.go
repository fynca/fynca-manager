package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type ctxKey string

var authTokenKey ctxKey

type CORSMiddleware struct {
	domain  string
	headers []string
}

func NewCORSMiddleware(domain string, headers []string) *CORSMiddleware {
	return &CORSMiddleware{
		domain:  domain,
		headers: headers,
	}
}

func (m *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", m.domain)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		headers := "Accept, Access-Control-Allow-Origin, DNT, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, " + strings.Join(m.headers, ",")
		w.Header().Set("Access-Control-Allow-Headers", headers)
		next.ServeHTTP(w, r)
	})
}

func withTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("X-Session-Token")
		// return access denied if no token and not an OPTIONS request
		if token == "" && r.Method != http.MethodOptions {
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}

		// add auth token to context
		ctx = context.WithValue(r.Context(), authTokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		fields := logrus.Fields{
			"addr":   r.RemoteAddr,
			"method": r.Method,
			"uri":    r.RequestURI,
		}
		if token != "" {
			fields["token"] = token
		}
		logrus.WithFields(fields).Info("api request")
		next.ServeHTTP(w, r)
	})
}
