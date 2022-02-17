package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"git.underland.io/ehazlett/fynca"
	api "git.underland.io/ehazlett/fynca/api/services/jobs/v1"
	"git.underland.io/ehazlett/fynca/client"
	"github.com/ehazlett/fynca-manager/version"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Config struct {
	GRPCAddress string
	// TLSCertificate is the certificate used for grpc communication
	TLSServerCertificate string
	// TLSKey is the key used for grpc communication
	TLSServerKey string
	// TLSClientCertificate is the client certificate used for communication
	TLSClientCertificate string
	// TLSClientKey is the client key used for communication
	TLSClientKey string
	// TLSInsecureSkipVerify disables certificate verification
	TLSInsecureSkipVerify bool
	// APICORSDomain is the allowed domains for requests
	APICORSDomain string
	// ListenAddr is the address for the manager to listen
	ListenAddr string
	// PublicDir is the directory of the web UI source
	PublicDir string
	// TraceEndpoint is the endpoint to send trace data
	TraceEndpoint string
	// Environment is the environment for trace identification
	Environment string
}

type Server struct {
	cfg *Config
}

// NewServer returns a new server instance
func NewServer(cfg *Config) (*Server, error) {
	srv := &Server{
		cfg: cfg,
	}

	return srv, nil
}

// Run starts the server
func (s *Server) Run() error {
	corsMiddleware := NewCORSMiddleware(s.cfg.APICORSDomain, []string{"x-session-token"})

	r := mux.NewRouter().StrictSlash(true)
	r.Methods("OPTIONS").HandlerFunc(apiOK)
	r.HandleFunc("/healthz", indexHandler)
	r.HandleFunc("/versionz", s.serverVersionHandler)
	r.Use(logMiddleware)
	r.Use(corsMiddleware.Handler)

	// auth
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Methods("OPTIONS").HandlerFunc(apiOK)
	authRouter.HandleFunc("/login", s.loginHandler).Methods("POST")

	// api
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(withTokenMiddleware)
	// v1
	v1 := apiRouter.PathPrefix("/v1").Subrouter()
	v1.Methods("OPTIONS").HandlerFunc(apiOK)
	// jobs
	v1.Handle("/jobs", otelhttp.NewHandler(http.HandlerFunc(s.jobsListHandler), "jobs.list")).Methods("GET")
	v1.Handle("/jobs", otelhttp.NewHandler(http.HandlerFunc(s.jobQueueHandler), "jobs.queue")).Methods("POST")
	v1.Handle("/jobs/{id:[0-9,a-z,-]+}", otelhttp.NewHandler(http.HandlerFunc(s.jobDetailsHandler), "job.details")).Methods("GET")
	v1.Handle("/jobs/{id:[0-9,a-z,-]+}/latest-render/{frame:[0-9]+}", otelhttp.NewHandler(http.HandlerFunc(s.jobLatestRenderHandler), "jobs.latestrender")).Methods("GET")
	v1.Handle("/jobs/{id:[0-9,a-z,-]+}", otelhttp.NewHandler(http.HandlerFunc(s.jobDeleteHandler), "job.delete")).Methods("DELETE")
	v1.Handle("/jobs/{id:[0-9,a-z,-]+}/log", otelhttp.NewHandler(http.HandlerFunc(s.jobLogHandler), "job.log")).Methods("GET")
	v1.Handle("/jobs/{id:[0-9,a-z,-]+}/archive", otelhttp.NewHandler(http.HandlerFunc(s.jobArchiveHandler), "job.archive")).Methods("GET")

	v1.Handle("/renders/{id:[0-9,a-z,-]+}/logs/{frame:[0-9]+}", otelhttp.NewHandler(http.HandlerFunc(s.renderLogHandler), "renders.log")).Methods("GET")

	v1.Handle("/workers", otelhttp.NewHandler(http.HandlerFunc(s.workersListHandler), "workers.list")).Methods("GET")
	v1.Handle("/workers/{name}/stop", otelhttp.NewHandler(http.HandlerFunc(s.workerStopHandler), "worker.stop")).Methods("POST")
	v1.Handle("/workers/{name}/update", otelhttp.NewHandler(http.HandlerFunc(s.workerUpdateHandler), "worker.update")).Methods("POST")

	v1.Handle("/accounts/profile", otelhttp.NewHandler(http.HandlerFunc(s.accountsProfileHandler), "accounts.profile")).Methods("GET")
	v1.Handle("/accounts/profile", otelhttp.NewHandler(http.HandlerFunc(s.accountsUpdateProfileHandler), "accounts.updateprofile")).Methods("POST")
	v1.Handle("/accounts/change-password", otelhttp.NewHandler(http.HandlerFunc(s.accountsChangePasswordHandler), "accounts.changepassword")).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(s.cfg.PublicDir)))

	srv := &http.Server{
		Handler:      r,
		Addr:         s.cfg.ListenAddr,
		WriteTimeout: 600 * time.Second,
		ReadTimeout:  600 * time.Second,
	}

	logrus.Infof("starting server on %s", s.cfg.ListenAddr)
	return srv.ListenAndServe()
}

func (s *Server) getClient(r *http.Request) (*client.Client, error) {
	cfg := &fynca.Config{
		GRPCAddress:           s.cfg.GRPCAddress,
		TLSServerCertificate:  s.cfg.TLSServerCertificate,
		TLSServerKey:          s.cfg.TLSServerKey,
		TLSClientCertificate:  s.cfg.TLSClientCertificate,
		TLSClientKey:          s.cfg.TLSClientKey,
		TLSInsecureSkipVerify: s.cfg.TLSInsecureSkipVerify,
	}

	opts := []client.ClientOpt{}
	token, ok := r.Context().Value(fynca.CtxTokenKey).(string)
	if ok {
		opts = append(opts, client.WithToken(token))
	}

	return client.NewClient(cfg, opts...)
}

func (s *Server) marshaler() *jsonpb.Marshaler {
	return &jsonpb.Marshaler{EmitDefaults: true}
}

func (s *Server) serverVersionHandler(w http.ResponseWriter, r *http.Request) {
	fc, err := s.getClient(r)
	if err != nil {
		logrus.WithError(err).Error("error getting fynca client")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fc.Close()

	ctx := context.Background()
	resp, err := fc.Version(ctx, &api.VersionRequest{})
	if err != nil {
		logrus.WithError(err).Error("error getting server version")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.WithError(err).Error("error encoding version")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("fynca manager " + version.BuildVersion() + "\n"))
}

func apiOK(w http.ResponseWriter, r *http.Request) {}
