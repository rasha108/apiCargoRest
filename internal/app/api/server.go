package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rasha108/apiCargoRest.git/internal/app/model"
	"github.com/rasha108/apiCargoRest.git/internal/app/store"

	"github.com/gorilla/sessions"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

const (
	sessionName = "gopherschool"
)

var (
	errIncorrectEmailOrPassword = errors.New("Incorrect email or password")
)

type server struct {
	router       *chi.Mux
	logger       *logrus.Logger
	config       *Config
	store        store.Store
	sessionStore sessions.Store
}

func NewServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       chi.NewRouter(),
		logger:       logrus.New(),
		config:       NewConfig(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configRouter() {
	router := s.router
	conf := s.config
	basePath := conf.BasePath
	router.Route(basePath, func(scope chi.Router) {
		// public routers
		scope.Group(func(public chi.Router) {
			public.Post("/users", s.HandleUserCreate)
			public.Post("/sessions", s.HandleSessionsCreate)
		})
	})
}

func (s *server) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"email"`
	}

	req := request{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := s.store.User().Create(u); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	u.Saintize()
	s.respond(w, r, http.StatusCreated, u)
}

func (s *server) HandleSessionsCreate(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req := request{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	u, err := s.store.User().FindByEmail(req.Email)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	session.Values["user_id"] = u.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusOK, nil)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	// Записываем оишбку в лог
	s.logger.Fatal(err)
	// Вывод ошибки
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
