package config

import (
	"os"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	JWTSecret    string
	YDBEndpoint  string
}

func LoadConfig() *Config {
	return &Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		AuthURL:      os.Getenv("AUTH_URL"),
		TokenURL:     os.Getenv("TOKEN_URL"),
		UserInfoURL:  os.Getenv("USER_INFO_URL"),
	}
}
