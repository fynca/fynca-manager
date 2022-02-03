package server

import (
	"context"
	"net/http"
	"strconv"

	api "git.underland.io/ehazlett/fynca/api/services/render/v1"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (s *Server) renderLogHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fr := params["frame"]
	if id == "" || fr == "" {
		http.Error(w, "id and frame must be specified", http.StatusBadRequest)
		return
	}
	slice := r.URL.Query().Get("slice")
	sliceIndex := int64(-1)
	if slice != "" {
		idx, err := strconv.Atoi(slice)
		if err != nil {
			http.Error(w, "slice must be a number", http.StatusBadRequest)
			return
		}
		sliceIndex = int64(idx)
	}

	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	frame, err := strconv.Atoi(fr)
	if err != nil {
		http.Error(w, "frame must be a number", http.StatusBadRequest)
		return
	}

	logrus.Debugf("getting render log frame=%d slice=%d", frame, sliceIndex)

	ctx := context.Background()
	resp, err := fc.RenderLog(ctx, &api.RenderLogRequest{
		ID:    id,
		Frame: int64(frame),
		Slice: sliceIndex,
	})
	if err != nil {
		logrus.WithError(err).Error("error getting render log")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.marshaler().Marshal(w, resp); err != nil {
		logrus.WithError(err).Error("error marshaling job")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
