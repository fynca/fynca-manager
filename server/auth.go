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
	"encoding/json"
	"net/http"

	api "git.underland.io/ehazlett/fynca/api/services/accounts/v1"
	"github.com/sirupsen/logrus"
)

type credentials struct {
	Username string
	Password string
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds *credentials
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		logrus.WithError(err).Error("error decoding credentials")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	ctx := context.Background()
	resp, err := fc.Authenticate(ctx, &api.AuthenticateRequest{
		Username: creds.Username,
		Password: []byte(creds.Password),
	})
	if err != nil {
		logrus.WithError(err).Error("invalid username or password")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// clear password crypt
	resp.Account.PasswordCrypt = nil

	logrus.Debugf("%+v", resp.Account)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.WithError(err).Error("error encoding user account")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
