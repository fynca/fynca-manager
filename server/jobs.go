package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	api "git.underland.io/ehazlett/fynca/api/services/jobs/v1"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (s *Server) jobsListHandler(w http.ResponseWriter, r *http.Request) {
	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	ctx := r.Context()
	resp, err := fc.ListJobs(ctx, &api.ListJobsRequest{})
	if err != nil {
		logrus.WithError(err).Error("error getting jobs")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// remove the frame jobs for a lighter client payload
	for _, job := range resp.Jobs {
		job.FrameJobs = []*api.FrameJob{}
	}

	if err := s.marshaler().Marshal(w, resp); err != nil {
		logrus.WithError(err).Error("error marshaling jobs")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) jobArchiveHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "id must be specified", http.StatusBadRequest)
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
	resp, err := fc.GetJobArchive(ctx, &api.GetJobArchiveRequest{
		ID: id,
	})
	if err != nil {
		logrus.WithError(err).Error("error getting job archive")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp.JobArchive.ArchiveUrl))
}

func (s *Server) jobDeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "id must be specified", http.StatusBadRequest)
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
	if _, err := fc.DeleteJob(ctx, &api.DeleteJobRequest{ID: id}); err != nil {
		logrus.WithError(err).Error("error deleting job")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) jobLatestRenderHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fr := params["frame"]
	if id == "" || fr == "" {
		http.Error(w, "id and frame must be specified", http.StatusBadRequest)
		return
	}

	frame, err := strconv.Atoi(fr)
	if err != nil {
		http.Error(w, "frame must be a number", http.StatusBadRequest)
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

	req := &api.GetLatestRenderRequest{
		ID:    id,
		Frame: int64(frame),
	}
	// ttl if specified
	if v := r.URL.Query().Get("ttl"); v != "" {
		ttl, err := time.ParseDuration(v)
		if err != nil {
			logrus.WithError(err).Error("error setting ttl for latest render")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.TTL = ttl
	}

	resp, err := fc.GetLatestRender(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("error getting latest render for job")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.marshaler().Marshal(w, resp); err != nil {
		logrus.WithError(err).Error("error marshaling latest render")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) jobDetailsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "id must be specified", http.StatusBadRequest)
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
	resp, err := fc.GetJob(ctx, &api.GetJobRequest{
		ID: id,
	})
	if err != nil {
		logrus.WithError(err).Error("error getting jobs")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.marshaler().Marshal(w, resp); err != nil {
		logrus.WithError(err).Error("error marshaling job")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) jobQueueHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(128 << 20)

	logrus.Debugf("%+v", r.Form)
	resX, err := strconv.Atoi(r.FormValue("renderResolutionX"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resY, err := strconv.Atoi(r.FormValue("renderResolutionY"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resolutionScale, err := strconv.Atoi(r.FormValue("resolutionScale"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderStartFrame, err := strconv.Atoi(r.FormValue("renderStartFrame"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderEndFrame, err := strconv.Atoi(r.FormValue("renderEndFrame"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderSamples, err := strconv.Atoi(r.FormValue("renderSamples"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderUseGPU, err := strconv.ParseBool(r.FormValue("renderUseGPU"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderSlices, err := strconv.Atoi(r.FormValue("renderSlices"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderEngine, err := parseRenderEngine(r.FormValue("renderEngine"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jobPriority, err := parseJobPriority(r.FormValue("priority"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jobRequest := &api.JobRequest{
		Name:             r.FormValue("name"),
		ResolutionX:      int64(resX),
		ResolutionY:      int64(resY),
		ResolutionScale:  int64(resolutionScale),
		RenderStartFrame: int64(renderStartFrame),
		RenderEndFrame:   int64(renderEndFrame),
		RenderSamples:    int64(renderSamples),
		RenderUseGPU:     renderUseGPU,
		RenderSlices:     int64(renderSlices),
		RenderEngine:     renderEngine,
		Priority:         jobPriority,
	}

	projectFile, _, err := r.FormFile("project")
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cBuf := make([]byte, 512)
	if _, err := projectFile.Read(cBuf); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	contentType := http.DetectContentType(cBuf)
	if _, err := projectFile.Seek(0, 0); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jobRequest.ContentType = contentType

	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	stream, err := fc.QueueJob(context.Background())
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := stream.Send(&api.QueueJobRequest{
		Data: &api.QueueJobRequest_Request{
			Request: jobRequest,
		},
	}); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rdr := bufio.NewReader(projectFile)
	buf := make([]byte, 4096)
	for {
		n, err := rdr.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req := &api.QueueJobRequest{
			Data: &api.QueueJobRequest_ChunkData{
				ChunkData: buf[:n],
			},
		}

		if err := stream.Send(req); err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(res.GetUUID()))
}

func (s *Server) jobLogHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "id and frame must be specified", http.StatusBadRequest)
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
	resp, err := fc.JobLog(ctx, &api.JobLogRequest{
		ID: id,
	})
	if err != nil {
		logrus.WithError(err).Error("error getting job log")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(resp.JobLog.Log))
}

func parseRenderEngine(e string) (api.RenderEngine, error) {
	switch strings.ToLower(e) {
	case "cycles":
		return api.RenderEngine_CYCLES, nil
	case "eevee":
		return api.RenderEngine_BLENDER_EEVEE, nil
	}
	return api.RenderEngine_UNKNOWN, fmt.Errorf("unsupported render engine: %q", e)
}

func parseJobPriority(p string) (api.JobPriority, error) {
	switch strings.ToLower(strings.TrimSpace(p)) {
	case "urgent":
		return api.JobPriority_URGENT, nil
	case "normal":
		return api.JobPriority_NORMAL, nil
	case "low":
		return api.JobPriority_LOW, nil
	}

	return api.JobPriority_NORMAL, fmt.Errorf("unknown priority %s", p)
}
