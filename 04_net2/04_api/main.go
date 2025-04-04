package main

import (
	"log"
	"net/http"
	"sync"

	"gitlab.vk-golang.ru/vk-golang/lectures/04_net2/04_api/handlers"
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
		"test": {
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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("start server error: %s", err.Error())
	}
}
