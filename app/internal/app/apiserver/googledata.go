package apiserver

//type GoogleClaims struct {
//	Email         string `json:"email"`
//	EmailVerified bool   `json:"email_verified"`
//	FirstName     string `json:"given_name"`
//	LastName      string `json:"family_name"`
//	jwt.StandardClaims
//}
//
//const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="
//
//var GoogleOauthConfig, errGoogleConfigLoad = config.NewGoogleConfig(config.Config)
//
//func init() {
//	if errGoogleConfigLoad != nil {
//		log.Fatal("error loading Google Config:", errGoogleConfigLoad)
//	}
//}
//
//func generateStateOauthCookie(w http.ResponseWriter) string {
//	var expiration = time.Now().Add(5 * time.Minute)
//
//	b := make([]byte, 16)
//	rand.Read(b)
//	state := base64.URLEncoding.EncodeToString(b)
//	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
//	http.SetCookie(w, &cookie)
//
//	return state
//}
//
//func getUserDataFromGoogle(code string) ([]byte, error) {
//	// Use code to get token and get user info from Google.
//
//	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
//	if err != nil {
//		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
//	}
//	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
//	if err != nil {
//		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
//	}
//	defer response.Body.Close()
//	contents, err := io.ReadAll(response.Body)
//	if err != nil {
//		return nil, fmt.Errorf("failed read response: %s", err.Error())
//	}
//	return contents, nil
//}
