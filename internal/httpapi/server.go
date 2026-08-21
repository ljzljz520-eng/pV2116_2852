package httpapi

import (
	"net/http"
	"stickerchallenge/internal/service"
)

type Server struct {
	Service *service.Service
	Mux     *http.ServeMux
}

func New(svc *service.Service) *Server {
	s := &Server{Service: svc, Mux: http.NewServeMux()}
	s.registerRoutes()
	return s
}

func (s *Server) Handler() http.Handler { return s.Mux }

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
