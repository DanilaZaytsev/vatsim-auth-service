package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"vatsim-auth-service/internal/model"
	"vatsim-auth-service/pkg/logger"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// BuildVatsimAuthURL формирует ссылку на авторизацию через VATSIM
func BuildVatsimAuthURL() string {
	baseURL := strings.TrimRight(os.Getenv("VATSIM_URL"), "/")
	clientID := os.Getenv("VATSIM_CLIENT_ID")
	redirectURI := os.Getenv("VATSIM_REDIRECT_URI")
	scope := url.QueryEscape("full_name email vatsim_details country")

	return fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		baseURL,
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		scope,
	)
}

// ExchangeCodeForToken меняет авторизационный код на токен доступа
func ExchangeCodeForToken(code string) (*TokenResponse, error) {
	baseURL := strings.TrimRight(os.Getenv("VATSIM_URL"), "/")

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", os.Getenv("VATSIM_REDIRECT_URI"))
	form.Set("client_id", os.Getenv("VATSIM_CLIENT_ID"))
	form.Set("client_secret", os.Getenv("VATSIM_CLIENT_SECRET"))

	req, err := http.NewRequest("POST", baseURL+"/oauth/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debug("Token response body: " + string(body))

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	logger.Debug("Received access token: " + token.AccessToken)
	return &token, nil
}

// GetUserInfo получает информацию о пользователе по access token
func GetUserInfo(accessToken string) (*model.UserInfoResponse, error) {
	baseURL := strings.TrimRight(os.Getenv("VATSIM_URL"), "/")

	req, err := http.NewRequest("GET", baseURL+"/api/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debug("User info response body: " + string(body))

	var user model.UserInfoResponse
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
