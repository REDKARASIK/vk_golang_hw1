package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (h *Handler) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")

	for session, user := range h.Sessions {
		if user.Login == login {
			resp := Response{"session_id - " + session}
			bytes, _ := json.Marshal(&resp)
			w.Write(bytes)
			return
		}
	}

	log.Printf("session for user %s not found", login)
	w.Header().Set("Content-Type", "application/json")
	resp := Response{"no session"}
	bytes, _ := json.Marshal(&resp)
	w.Write(bytes)
}

func (h *Handler) HandleDeleteSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("cookie err: %s", err)
		w.Header().Set("Content-Type", "application/json")
		resp := Response{"no cookie"}
		bytes, _ := json.Marshal(&resp)
		w.Write(bytes)
		return
	}

	cookie.Expires = time.Now().Add(-1)

	http.SetCookie(w, cookie)

	h.Mu.Lock()
	delete(h.Sessions, cookie.Value)
	h.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (h *Handler) HandleSession(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandleGetSession(w, r)
		return
	case http.MethodDelete:
		h.HandleDeleteSession(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	resp := Response{"wrong method"}
	bytes, _ := json.Marshal(&resp)
	w.Write(bytes)
}
