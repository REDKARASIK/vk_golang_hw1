package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

/*

curl -X POST -H "Content-Type: application/json" -d '{"login": "rvasily", "password": "love"}' http://localhost:8080/login

curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NDUwLCJsb2dpbiI6InJ2YXNpbHkiLCJuYW1lIjoiVmFzaWx5IFJvbWFub3YiLCJyb2xlIjoidXNlciJ9.Y0FJFm8fSbjc4nzBa1LHJSxNRRYp-chOZLr26sOJSgo" http://localhost:8080/profile

*/

type User struct {
	ID       int
	FullName string
	Role     string
}

type Claims struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	ID    int    `json:"id"`
	jwt.RegisteredClaims
}

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	users = map[string]User{
		"rvasily":        {450, "Vasily Romanov", "user"}, // a tribute to the original course author
		"romanov.vasily": {42, "Василий Романов", "admin"},
	}

	ExamplePassword    = "love"
	ExampleTokenSecret = []byte("супер секретный ключ")
)

func profilePage(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	// Extract token from "Bearer <token>" format
	inToken := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		inToken = authHeader[7:]
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(inToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad sign method")
		}
		return ExampleTokenSecret, nil
	})
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "bad token")
		return
	}
	if !token.Valid {
		jsonError(w, http.StatusUnauthorized, "bad token")
		return
	}

	resp, _ := json.Marshal(map[string]interface{}{
		"status": http.StatusOK,
		"data": map[string]interface{}{
			"login": claims.Login,
			"name":  claims.Name,
			"role":  claims.Role,
			"id":    claims.ID,
		},
	})
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	body, _ := io.ReadAll(r.Body)
	r.Body.Close()

	fd := &LoginForm{}
	err := json.Unmarshal(body, fd)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	user, exist := users[fd.Login]
	if !exist || fd.Password != ExamplePassword {
		jsonError(w, http.StatusUnauthorized, "bad login or password")
		return
	}

	// Create token with custom claims
	claims := &Claims{
		Login:            fd.Login,
		Name:             user.FullName,
		Role:             user.Role,
		ID:               user.ID,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ExampleTokenSecret)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp, _ := json.Marshal(map[string]interface{}{
		"status": http.StatusOK,
		"data": map[string]interface{}{
			"token": tokenString,
		},
	})
	w.Write(resp)
	w.Write([]byte("\n\n"))
}

func main() {
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/profile", profilePage)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func jsonError(w io.Writer, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	w.Write(resp)
}
