package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"gitlab.vk-golang.ru/vk-golang/lectures/04_net2/02_rest/storage"

	"github.com/gorilla/mux"
)

// GET - получение
// POST - добавление новых данных
// PUT - изменение данных
// DELETE - удаление

// HEAD - получить только заголовки (metadata)
// PATCH
// OPTIONS

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type BooksHandler struct {
	store *storage.BookStore
}

func (api *BooksHandler) List(w http.ResponseWriter, _ *http.Request) {

	books, err := api.store.GetBooks()
	if err != nil {
		http.Error(w, `{"error":"db"}`, 500)
		return
	}

	body := map[string]interface{}{
		"books": books,
	}
	json.NewEncoder(w).Encode(&Result{Body: body})
}

// POST http://127.0.0.1:8080/book/ with form title=test&price=123

func (api *BooksHandler) Add(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	price, _ := strconv.Atoi(r.FormValue("price"))

	in := storage.Book{
		Title: title,
		Price: uint(price),
	}

	addedBook, err := api.store.AddBook(in)
	if err != nil {
		http.Error(w, `{"error":"db"}`, 500)
		return
	}

	body := map[string]interface{}{
		"book": addedBook,
	}
	json.NewEncoder(w).Encode(&Result{Body: body})
}

func (api *BooksHandler) BookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		return
	}

	book, err := api.store.GetBook(id)
	if errors.Is(err, storage.NotFound) {
		http.Error(w, `{"error":"not found"}`, 404)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"db"}`, 500)
		return
	}

	body := map[string]interface{}{
		"book": book,
	}

	json.NewEncoder(w).Encode(&Result{Body: body})
}

func (api *BooksHandler) Change(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"bad id"}`, 400)
		return
	}

	title := r.FormValue("title")
	price, _ := strconv.Atoi(r.FormValue("price"))

	in := storage.Book{
		ID:    id,
		Title: title,
		Price: uint(price),
	}

	addedBook, err := api.store.Change(in)
	if errors.Is(err, storage.NotFound) {
		http.Error(w, `{"error":"not found"}`, 404)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"db"}`, 500)
		return
	}

	body := map[string]interface{}{
		"book": addedBook,
	}

	json.NewEncoder(w).Encode(&Result{Body: body})
}

func main() {
	r := mux.NewRouter()

	api := &BooksHandler{
		store: storage.NewBookStore(),
	}

	r.HandleFunc("/book/", api.List).Methods("GET")
	r.HandleFunc("/book/", api.Add).Methods("POST")
	r.HandleFunc("/book/{id:[0-9]+}", api.BookByID).Methods("GET")
	r.HandleFunc("/book/{id:[0-9]+}", api.Change).Methods("PUT")

	log.Println("start serving :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Printf("start server error: %s", err.Error())
	}
}
