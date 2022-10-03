package handlers

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

type Authn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type MalformedRequest struct {
	status int
	msg    string
}

type User struct {
	id                int
	login             string
	encryptedPassword string
}
