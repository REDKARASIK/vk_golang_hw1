package main

import (
	"net/http"
	"sync"

	"gitlab.vk-golang.com/vk-golang/lectures/04_net2/04_api/handlers"
)

// GET - получение
// POST - добавление новых данных
// PUT - изменение данных
// DELETE - удаление

// HEAD
// PATCH
// OPTIONS

func main() {

	users := map[string]*handlers.User{
		"test": &handlers.User{
			ID:       1,
			Login:    "test",
			Password: "test",
		},
	}

	sessions := map[string]*handlers.User{
		"tokenknsjkdfklsdf": users["test"],
	}

	mu := &sync.Mutex{}

	handler := handlers.Handler{
		Sessions: sessions,
		Users:    users,
		Mu:       mu,
	}

	http.HandleFunc("/users/", handler.HandleUsers)
	http.HandleFunc("/session/", handler.HandleSession)

	http.ListenAndServe(":8080", nil)
}
