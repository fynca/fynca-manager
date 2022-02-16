package server

import (
	"context"
	"net/http"

	api "git.underland.io/ehazlett/fynca/api/services/jobs/v1"
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
