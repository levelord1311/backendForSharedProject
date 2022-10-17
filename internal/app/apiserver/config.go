package apiserver


type Config struct {
	BindAddr    string `toml:"bind_addr"`
	TLSAddr     string `toml:"tls_addr"`
	Cert        string `toml:"cert"`
	Key         string `toml:"key"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`

}

func NewConfig() *Config {
	return &Config{

		BindAddr: ":80",
		TLSAddr:  ":443",
	}


}
