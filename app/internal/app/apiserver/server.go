package apiserver

import (
	"backendForSharedProject/internal/app/model"
	"backendForSharedProject/internal/app/store"
	"backendForSharedProject/pkg/jwt"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	router *mux.Router
	store  store.Store
	jwtKey []byte
}

var (
	errIncorrectUsernameEmailOrPassword = errors.New("incorrect username/email or password")
)

func newServer(store store.Store, jwtKey []byte) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
		jwtKey: jwtKey,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/", s.handleDefaultPage())
	s.router.HandleFunc("/users/create", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/auth", s.handleJWTCreate()).Methods("POST")
	//s.router.HandleFunc("/auth/google", s.handleRedirectToGoogleLogin())
	//s.router.HandleFunc("/auth/google/callback", s.handleGoogleCallback())
	s.router.HandleFunc("/estate_lots/create", s.handleEstateLotsCreate()).Methods("POST")
	s.router.HandleFunc("/estate_lots/get_fresh", s.handleGetFreshEstateLots())
}

func (s *server) handleDefaultPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! This is default handler, now serving host: %s", r.Host)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
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
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().CreateUser(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleEstateLotsCreate() http.HandlerFunc {
	type request struct {
		TypeOfEstate string `json:"type_of_estate"`
		Rooms        int    `json:"rooms"`
		Area         int    `json:"area"`
		Floor        int    `json:"floor"`
		MaxFloor     int    `json:"max_floor"`
		City         string `json:"city"`
		District     string `json:"district"`
		Street       string `json:"street"`
		Building     string `json:"building"`
		Price        int    `json:"price"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		lot := &model.EstateLot{
			TypeOfEstate: req.TypeOfEstate,
			Rooms:        req.Rooms,
			Area:         req.Area,
			Floor:        req.Floor,
			MaxFloor:     req.MaxFloor,
			City:         req.City,
			District:     req.District,
			Street:       req.Street,
			Building:     req.Building,
			Price:        req.Price,
		}
		if err := s.store.EstateLot().CreateEstateLot(lot); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, lot)
	}
}

func (s *server) handleGetFreshEstateLots() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lots, err := s.store.EstateLot().GetFreshEstateLots()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusCreated, &lots)
	}
}

func (s *server) handleJWTCreate() http.HandlerFunc {
	// TODO store all DTO's in separate model file
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		//validate both login and password are present
		if err := validation.ValidateStruct(req,
			validation.Field(&req.Login, validation.Required),
			validation.Field(&req.Password, validation.Required)); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u := &model.User{}
		if validation.Validate(req.Login, is.Email) == nil {
			u, err := s.store.User().FindByEmail(req.Login)
			if err != nil || !u.ComparePassword(req.Password) {
				s.error(w, r, http.StatusUnauthorized, errIncorrectUsernameEmailOrPassword)
				return
			}
		} else {
			u, err := s.store.User().FindByUsername(req.Login)
			if err != nil || !u.ComparePassword(req.Password) {
				s.error(w, r, http.StatusUnauthorized, errIncorrectUsernameEmailOrPassword)
				return
			}
		}
		tokenString, err := jwt.GenerateJWT(u, s.jwtKey)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, tokenString)

	}
}

//func (s *server) handleRedirectToGoogleLogin() http.HandlerFunc {
//
//	return func(w http.ResponseWriter, r *http.Request) {
//		state := generateStateOauthCookie(w)
//		url := GoogleOauthConfig.AuthCodeURL(state)
//		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
//	}
//}
//
//func (s *server) handleGoogleCallback() http.HandlerFunc {
//	type GoogleInfo struct {
//		Sub           string `json:"sub"`
//		Name          string `json:"name"`
//		GivenName     string `json:"given_name"`
//		FamilyName    string `json:"family_name"`
//		Picture       string `json:"picture"`
//		Email         string `json:"email"`
//		EmailVerified bool   `json:"email_verified"`
//		Locale        string `json:"locale"`
//	}
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		// Read oauthState from Cookie
//		oauthState, _ := r.Cookie("oauthstate")
//
//		if r.FormValue("state") != oauthState.Value {
//			s.error(w, r, http.StatusBadRequest, errors.New("invalid oauth google state"))
//			return
//		}
//
//		data, err := getUserDataFromGoogle(r.FormValue("code"))
//		if err != nil {
//			log.Println(err.Error())
//			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
//			return
//		}
//
//		googleInfo := &GoogleInfo{}
//		err = json.Unmarshal(data, googleInfo)
//		if err != nil {
//			s.error(w, r, http.StatusInternalServerError, err)
//			return
//		}
//
//		u, err := s.store.User().FindByEmailGoogle(googleInfo.Email)
//		if err == store.ErrRecordNotFound {
//			if !googleInfo.EmailVerified {
//				s.error(w, r, http.StatusBadRequest, errors.New("can't create new user: "+
//					"google email address is not verified"))
//				return
//			}
//
//			u := &model.User{
//				Email:      googleInfo.Email,
//				GivenName:  googleInfo.GivenName,
//				FamilyName: googleInfo.FamilyName,
//			}
//
//			if err := s.store.User().CreateUserWithGoogle(u); err != nil {
//				s.error(w, r, http.StatusInternalServerError, err)
//				return
//			}
//
//			tokenString, err := jwt.GenerateJWT(u, s.jwtKey)
//			if err != nil {
//				s.error(w, r, http.StatusInternalServerError, err)
//				return
//			}
//
//			s.respond(w, r, http.StatusCreated, tokenString)
//			return
//
//		} else if err != nil {
//			s.error(w, r, http.StatusInternalServerError, err)
//			return
//		}
//
//		tokenString, err := jwt.GenerateJWT(u, s.jwtKey)
//		if err != nil {
//			s.error(w, r, http.StatusInternalServerError, err)
//			return
//		}
//
//		s.respond(w, r, http.StatusOK, tokenString)
//		return
//
//	}
//}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
