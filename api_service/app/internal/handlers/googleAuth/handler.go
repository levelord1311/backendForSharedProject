package googleAuth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/client/user_service"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/config"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/jwt"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"io"
	"net/http"
	"time"
)

const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="

	googleRedirectURL = "/auth/google"
	googleCallbackURL = "/auth/google/callback"
)

type Handler struct {
	Logger      logging.Logger
	JWTHelper   jwt.Helper
	UserService user_service.UserService
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, googleRedirectURL, apperror.Middleware(h.handleRedirectToGoogleLogin))
	router.HandlerFunc(http.MethodGet, googleCallbackURL, apperror.Middleware(h.handleGoogleCallback))
}

func (h *Handler) handleRedirectToGoogleLogin(w http.ResponseWriter, r *http.Request) error {
	cfg := config.GetGoogleConfig()
	state := generateStateOauthCookie(w)
	url := cfg.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func (h *Handler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) error {
	type googleInfo struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Locale        string `json:"locale"`
	}
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		return apperror.BadRequestError("invalid oauth google state")
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		return err
	}

	googleData := &googleInfo{}
	err = json.Unmarshal(data, googleData)
	if err != nil {
		return err
	}

	// TODO implement user service (find and verify user by email or create new user)

	return nil
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(5 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	cfg := config.GetGoogleConfig()

	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
