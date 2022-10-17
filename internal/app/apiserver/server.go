package apiserver

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	"log"
	"net"
	"net/http"
)

type server struct {
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
}

const (
	sessionName = "default_session"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

func newTLSServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureTLSRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func redirectToTls(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	u.Scheme = "https"
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		u.Host = r.Host
		log.Printf("Error splitting host from port: %s, result in url: %s\n", err, u)
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		return
	}

	// нужно придумать динамическую подгрузку из конфига
	u.Host = net.JoinHostPort(host, "443") // !!! порт захардкоден !!!

	log.Printf("%s is redirected to %s", r.Host, u.String())
	http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
}

func (s *server) configureTLSRouter() {
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
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
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
