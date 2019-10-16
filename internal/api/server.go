package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *chi.Mux
	logger *logrus.Logger
}

func NewServer() *server {
	s := &server{
		router: chi.NewRouter(),
		logger: logrus.New(),
	}

	s.configRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configRouter() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
}
