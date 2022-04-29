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
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (s *Server) proxyMediaHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	urlContent := params["content"]
	if urlContent == "" {
		return
	}

	data, err := base64.StdEncoding.DecodeString(urlContent)
	if err != nil {
		logrus.WithError(err).Error("error during content decode for %s", urlContent)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mediaUrl := string(data)
	logrus.Debugf("proxying media from %s", mediaUrl)

	resp, err := http.Get(mediaUrl)
	if err != nil {
		logrus.WithError(err).Error("error during media proxy for %s", mediaUrl)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	logrus.Debugf("proxy response: %+v", resp)

	// copy headers
	for name, headers := range resp.Header {
		for _, value := range headers {
			w.Header().Set(name, value)
		}
	}
	// check for extra params to add headers
	if v := r.URL.Query().Get("filename"); v != "" {
		logrus.Debugf("proxy: setting content disposition filename to %s", v)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", v))
	}
	if v := r.URL.Query().Get("type"); v != "" {
		logrus.Debugf("proxy: setting content type to %s", v)
		w.Header().Set("Content-Type", v)
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
