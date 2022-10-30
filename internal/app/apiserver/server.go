package apiserver

import (
	"backendForSharedProject/internal/app/jwt"
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

type server struct {
	router *mux.Router
	store  store.Store
	jwtKey []byte
}

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

func newServer(store store.Store, config *Config) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
		jwtKey: config.JwtKey,
	}

	s.configureRouter()

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

func (s *server) configureRouter() {
	s.router.HandleFunc("/", s.handleDefaultPage())
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/auth", s.handleJWTCreate()).Methods("POST")
	s.router.HandleFunc("/auth/google", s.handleRedirectToGoogleLogin())
	s.router.HandleFunc("/auth/google/callback", s.handleGoogleCallback())
}

func (s *server) handleDefaultPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! This is default handler, now serving host: %s", r.Host)
	}
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

func (s *server) handleJWTCreate() http.HandlerFunc {
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

		tokenString, err := jwt.GenerateJWT(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, tokenString)

	}
}

func (s *server) handleRedirectToGoogleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := generateStateOauthCookie(w)
		url := GoogleOauthConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (s *server) handleGoogleCallback() http.HandlerFunc {
	type GoogleInfo struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Locale        string `json:"locale"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		// Read oauthState from Cookie
		oauthState, _ := r.Cookie("oauthstate")

		if r.FormValue("state") != oauthState.Value {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid oauth google state"))
			return
		}

		data, err := getUserDataFromGoogle(r.FormValue("code"))
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		googleInfo := &GoogleInfo{}
		err = json.Unmarshal(data, googleInfo)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u, err := s.store.User().FindByEmail(googleInfo.Email)
		if err == store.ErrRecordNotFound {
			if !googleInfo.EmailVerified {
				s.error(w, r, http.StatusBadRequest, errors.New("can't create new user: "+
					"google email address is not verified"))
				return
			}

			u := &model.User{
				Email:      googleInfo.Email,
				GivenName:  googleInfo.GivenName,
				FamilyName: googleInfo.FamilyName,
			}

			if err := s.store.User().CreateWithGoogle(u); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			tokenString, err := jwt.GenerateJWT(u)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			s.respond(w, r, http.StatusCreated, tokenString)

		} else if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		tokenString, err := jwt.GenerateJWT(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, tokenString)
		return

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
