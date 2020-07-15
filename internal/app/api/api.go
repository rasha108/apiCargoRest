package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/rasha108/apiCargoRest.git/internal/app/rabbitmq"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/rasha108/apiCargoRest.git/internal/app/db"
	"github.com/rasha108/apiCargoRest.git/internal/app/model"
	logger "github.com/sirupsen/logrus"
)

const (
	sessionName        = "gopherschool"
	cxtKeyUser  cxtKey = iota
	cxtKeyRequestID
)

type cxtKey int8

var (
	errIncorrectEmailOrPassword = errors.New("Incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	router       *chi.Mux
	logger       *logger.Logger
	config       *Config
	store        db.Store
	sessionStore sessions.Store
	mqConnection *rabbitmq.Connection
}

func NewServer(store db.Store, sessionStore sessions.Store, config *Config, mqConnection *rabbitmq.Connection) *server {
	s := &server{
		router:       chi.NewRouter(),
		logger:       logger.New(),
		config:       config,
		store:        store,
		sessionStore: sessionStore,
		mqConnection: mqConnection,
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
	basePath := conf.APIConfig.BasePath
	router.Route(basePath, func(scope chi.Router) {
		scope.Use(s.setRequestID)
		scope.Use(s.logRequest)
		// public routers
		scope.Group(func(public chi.Router) {
			public.Post("/users", s.HandleUserCreate)
			public.Post("/sessions", s.HandleSessionsCreate)

			public.Route("/organization", func(send chi.Router) {
				send.Get("/sendToEmail", s.HandleSendEmail)
			})

			// private routers
			scope.Group(func(private chi.Router) {
				private.Route("/private", func(r chi.Router) {
					private.Use(s.authenticatedUser)
					private.Get("/whoami", s.HandleWhoami)
				})
			})
		})
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), cxtKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logger.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(cxtKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		rw := &responseWriter{
			w,
			http.StatusOK,
		}
		next.ServeHTTP(rw, r)
		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	u := &model.User{}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := s.store.User().Create(u); err != nil {
		logger.
			WithError(err).
			Error("sql failed users")
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	// Чет не работает, должна прятать пароль при запросе, но не прячет возвращается password = ""
	u.Saintize()
	s.respond(w, r, http.StatusCreated, u)
}

func (s *server) HandleSessionsCreate(w http.ResponseWriter, r *http.Request) {
	req := &model.User{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	u, err := s.store.User().FindByEmail(req.Email)
	if err != nil || !u.ComparePassword(req.Password) {
		logger.
			WithError(err).
			Error("sql failed users")
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

func (s *server) authenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			logger.
				WithError(err).
				Error("sql failed organizations")
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), cxtKeyUser, u)))
	})
}

func (s *server) HandleWhoami(w http.ResponseWriter, r *http.Request) {
	s.respond(w, r, http.StatusOK, r.Context().Value(cxtKeyUser).(*model.User))
}

func (s *server) HandleSendEmail(w http.ResponseWriter, r *http.Request) {
	rabbitMQ := s.mqConnection

	id := uuid.MustParse("12312312321")
	organizations, err := s.store.User().Organizations(id)
	if err != nil {
		logger.
			WithError(err).
			Error("sql failed organizations")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	To := []string{"rashid1995bik@mail.ru"}
	mail, err := rabbitmq.NewSimpleMail("rashid1995bik@mail.ru", To, "123", "nananananannan")
	if err != nil {
		logger.
			WithError(err).
			Error("sql failed organizations")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	err = rabbitMQ.Send(mail)
	if err != nil {
		logger.
			WithError(err).
			Error("sql failed organizations")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	s.respond(w, r, http.StatusOK, organizations)
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
