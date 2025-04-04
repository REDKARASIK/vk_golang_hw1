package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) HandleUsers(w http.ResponseWriter, _ *http.Request) {
	usersList := make([]User, 0)

	for _, user := range h.Users {
		usersList = append(usersList, *user)
	}

	bytes, _ := json.Marshal(&usersList)
	w.Write(bytes)
}
