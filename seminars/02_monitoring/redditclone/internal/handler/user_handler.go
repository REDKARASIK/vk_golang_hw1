package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"redditclone/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	repo user.Repo
}

func NewUserHandler(repo user.Repo) *UserHandler {
	return &UserHandler{repo: repo}
}

var jwtSecret = []byte("supersecretkey") // In production, use env/config

type jwtResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "invalid request"}`, http.StatusBadRequest)
		return
	}
	u, err := h.repo.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, `{"message": "user exists"}`, http.StatusConflict)
		return
	}
	token, err := generateJWT(u)
	if err != nil {
		http.Error(w, `{"message": "could not create token"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(jwtResponse{Token: token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "invalid request"}`, http.StatusBadRequest)
		return
	}
	u, err := h.repo.Authorize(req.Username, req.Password)
	if err != nil || u == nil {
		http.Error(w, `{"message": "invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	token, err := generateJWT(u)
	if err != nil {
		http.Error(w, `{"message": "could not create token"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jwtResponse{Token: token})
}

func generateJWT(u *user.User) (string, error) {
	claims := jwt.MapClaims{
		"user": map[string]string{
			"id":       u.ID,
			"username": u.Username,
		},
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
