package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"vatsim-auth-service/internal/db"
	"vatsim-auth-service/internal/jwt"
	"vatsim-auth-service/internal/repository"
	"vatsim-auth-service/internal/service"
	"vatsim-auth-service/pkg/logger"
)

func VatsimLoginHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("VATSIM_CLIENT_ID")
	redirectURI := os.Getenv("VATSIM_REDIRECT_URI")
	baseURL := os.Getenv("VATSIM_URL")
	scope := url.QueryEscape("full_name email vatsim_details country")
	authURL := fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		baseURL, url.QueryEscape(clientID), url.QueryEscape(redirectURI), scope)

	log.Debug().Str("url", authURL).Msg("Redirecting user to VATSIM")
	http.Redirect(w, r, authURL, http.StatusFound)
}

func VatsimCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		logger.Info("No code provided in callback")
		return
	}
	token, err := service.ExchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		logger.Error(err, "Token exchange failed")
		return
	}

	user, err := service.GetUserInfo(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		logger.Error(err, "Fetching user info failed")
		return
	}

	cidUint, _ := strconv.ParseUint(user.Data.CID, 10, 64)
	repo := repository.NewUserRepository(db.Table)

	role := "pilot"
	existingRole, err := repo.GetUserRole(r.Context(), cidUint)
	if err == nil && existingRole != "" {
		role = existingRole
	}
	err = repo.UpsertUser(
		r.Context(),
		cidUint,
		user.Data.Personal.Email,
		user.Data.Personal.FullName,
		token.RefreshToken,
		role,
	)
	if err != nil {
		logger.Error(err, "Failed to save user in YDB")
	}

	jwtToken, err := jwt.GenerateToken(cidUint, user.Data.Personal.Email, role, user.Data.Personal.Country.Name, user.Data.VATSIM.Division.Name)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		logger.Error(err, "JWT generation failed")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	log.Info().
		Str("cid", user.Data.CID).
		Str("email", user.Data.Personal.Email).
		Str("country_id", user.Data.Personal.Country.ID).
		Str("country_name", user.Data.Personal.Country.Name).
		Msg("User authenticated successfully")

	fmt.Fprintf(w, "Привет, %s! Токен установлен в cookie 🚀", user.Data.Personal.FullName)
}
