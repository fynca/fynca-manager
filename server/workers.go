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
	"fmt"
	"net/http"
	"net/url"
	"strings"

	api "github.com/fynca/fynca/api/services/workers/v1"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (s *Server) workersListHandler(w http.ResponseWriter, r *http.Request) {
	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	ctx := context.Background()
	resp, err := fc.ListWorkers(ctx, &api.ListWorkersRequest{})
	if err != nil {
		logrus.WithError(err).Error("error getting workers")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := &jsonpb.Marshaler{}
	if err := m.Marshal(w, resp); err != nil {
		logrus.WithError(err).Error("error marshaling workers")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) workerStopHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	workerName := params["name"]
	if workerName == "" {
		http.Error(w, "name must be specified", http.StatusBadRequest)
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
	if _, err := fc.ControlWorker(ctx, &api.ControlWorkerRequest{
		WorkerID: workerName,
		Message: &api.ControlWorkerRequest_Stop{
			Stop: &api.WorkerStop{},
		},
	}); err != nil {
		logrus.WithError(err).Errorf("error stopping worker %s", workerName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) workerUpdateHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	workerName := params["name"]
	if workerName == "" {
		http.Error(w, "name must be specified", http.StatusBadRequest)
		return
	}

	updateURL := r.FormValue("url")
	if updateURL == "" {
		http.Error(w, "url must be specified", http.StatusBadRequest)
		return
	}
	if !strings.HasSuffix(updateURL, "/") {
		updateURL += "/"
	}

	if _, err := url.Parse(updateURL); err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid update url", updateURL), http.StatusBadRequest)
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
	if _, err := fc.ControlWorker(ctx, &api.ControlWorkerRequest{
		WorkerID: workerName,
		Message: &api.ControlWorkerRequest_Update{
			Update: &api.WorkerUpdate{
				URL: updateURL,
			},
		},
	}); err != nil {
		logrus.WithError(err).Errorf("error updating worker %s", workerName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) workerPauseHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	workerName := params["name"]
	if workerName == "" {
		http.Error(w, "name must be specified", http.StatusBadRequest)
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
	if _, err := fc.ControlWorker(ctx, &api.ControlWorkerRequest{
		WorkerID: workerName,
		Message: &api.ControlWorkerRequest_Pause{
			Pause: &api.WorkerPause{},
		},
	}); err != nil {
		logrus.WithError(err).Errorf("error pausing worker %s", workerName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) workerResumeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	workerName := params["name"]
	if workerName == "" {
		http.Error(w, "name must be specified", http.StatusBadRequest)
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
	if _, err := fc.ControlWorker(ctx, &api.ControlWorkerRequest{
		WorkerID: workerName,
		Message: &api.ControlWorkerRequest_Resume{
			Resume: &api.WorkerResume{},
		},
	}); err != nil {
		logrus.WithError(err).Errorf("error resuming worker %s", workerName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
