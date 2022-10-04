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
	Status int
	Msg    string
}
