package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	config "github.com/spf13/viper"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	router *mux.Router
	server *http.Server
}

func New() (*Server, error) {

	r := mux.NewRouter()
	r.Use(httpLoggerMiddleware())

	r = r.PathPrefix("/api/v1").Subrouter()

	s := &Server{
		logger: zap.S().With("package", "server"),
		router: r,
	}

	return s, nil
}

func (s *Server) ListenAndServe() error {
	s.server = &http.Server{
		Addr:    net.JoinHostPort(config.GetString("server.host"), config.GetString("server.port")),
		Handler: s.router,
	}

	// Listen
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("could not listen on %s: %v", s.server.Addr, err)
	}
	s.logger.Infow("API Listening", "address", s.server.Addr)

	s.router.NotFoundHandler = http.HandlerFunc(s.notFound)
	s.router.HandleFunc("/version", s.getVersion()).Methods(http.MethodGet)

	if err = s.server.Serve(listener); err != nil {
		s.logger.Fatalw("API Listen error", "error", err, "address", s.server.Addr)
	}

	return nil
}

// Router returns the router
func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	s.logger.Infow("not found", "path", r.RequestURI)
	w.WriteHeader(http.StatusNotFound)
	type err struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
	json.NewEncoder(w).Encode(err{Message: "not found", Code: http.StatusNotFound})
}

// getVersion returns version
func (s *Server) getVersion() http.HandlerFunc {

	// Simple version struct
	type version struct {
		Version string `json:"version"`
	}
	var v = &version{Version: "v1.0.0"} // TODO

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(v)
	}
}
