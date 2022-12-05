package config

import (
	"backendForSharedProject/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	JWT     struct {
		Secret string `yaml:"secret" env-required:"true"`
	}
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	Database struct {
		Url string `yaml:"url" env-required:"true"`
	}
	//GoogleClient struct {
	//	ID     string `yaml:"google_client_id" env-required:"false"`
	//	Secret string `yaml:"google_secret" env-required:"true"`
	//}
	//CategoryService struct {
	//	URL string `yaml:"url" env-required:"true"`
	//} `yaml:"category_service" env-required:"true"`
	//UserService struct {
	//	URL string `yaml:"url" env-required:"true"`
	//} `yaml:"user_service" env-required:"true"`
	//NoteService struct {
	//	URL string `yaml:"url" env-required:"true"`
	//} `yaml:"note_service" env-required:"true"`
	//TagService struct {
	//	URL string `yaml:"url" env-required:"true"`
	//} `yaml:"tag_service" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("reading application config...")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
		logger.Info("generating google config...")

	})
	return instance
}

//func NewGoogleConfig(config *Config) (*oauth2.Config, error) {
//	googleClientID := config.GoogleClient.ID
//
//	googleClientSec := config.GoogleClient.Secret
//
//	return &oauth2.Config{
//		RedirectURL:  "https://backend-server-36962.herokuapp.com/auth/google/callback",
//		ClientID:     googleClientID,
//		ClientSecret: googleClientSec,
//		Scopes: []string{
//			"https://www.googleapis.com/auth/userinfo.email",
//			"https://www.googleapis.com/auth/userinfo.profile",
//		},
//		Endpoint: google.Endpoint,
//	}, nil
//
//}
