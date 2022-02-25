// Copyright 2022 Evan Hazlett
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package server

import (
	"context"
	"net/http"
	"strings"

	"git.underland.io/ehazlett/fynca"
	"github.com/sirupsen/logrus"
)

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
		ctx = context.WithValue(r.Context(), fynca.CtxTokenKey, token)
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
