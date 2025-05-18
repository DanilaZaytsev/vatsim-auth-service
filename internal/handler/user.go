package handler

import (
	"encoding/json"
	"net/http"
	"vatsim-auth-service/internal/db"
	"vatsim-auth-service/internal/middleware"
	"vatsim-auth-service/internal/repository"
)

type MeResponse struct {
	CID          uint64 `json:"cid"`
	Email        string `json:"email"`
	Roles        string `json:"roles"`
	CountryName  string `json:"country_name"`
	DivisionName string `json:"division_name"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	repo := repository.NewUserRepository(db.Table)
	role, err := repo.GetUserRole(r.Context(), claims.CID)
	if err != nil {
		http.Error(w, "failed to get user role", http.StatusInternalServerError)
		return
	}

	resp := MeResponse{
		CID:          claims.CID,
		Email:        claims.Email,
		Roles:        role,
		CountryName:  claims.CountryName,
		DivisionName: claims.DivisionName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/auth/vatsim/login", http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"token":"Bearer ` + cookie.Value + `"}`))
}

type UpdateRolesRequest struct {
	CID   uint64 `json:"cid"`
	Roles string `json:"roles"`
}

func UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil || claims.Roles != "admin" {
		http.Error(w, "forbidden — admin only", http.StatusForbidden)
		return
	}

	var req struct {
		CID  uint64 `json:"cid"`
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	repo := repository.NewUserRepository(db.Table)
	err := repo.UpdateUserRole(r.Context(), req.CID, req.Role)
	if err != nil {
		http.Error(w, "failed to update role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Role updated"))
}
