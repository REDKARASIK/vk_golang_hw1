package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	client = http.Client{Timeout: time.Duration(time.Millisecond)}

	ErrResource = errors.New("resource error")
)

func getRemoteResource() error {
	url := "http://127.0.0.1:9999/pages?id=123"
	_, err := client.Get(url)
	if err != nil {
		return ErrResource
	}
	return nil
}

func getRemoteResourceAndSomeWork() error {
	err := getRemoteResource()
	if err != nil {
		return fmt.Errorf("getRemoteResource: %w", err)
	}
	// do work with it
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := getRemoteResource()
	// err := getRemoteResourceAndSomeWork()
	if err != nil {
		fmt.Printf("error happend: %+v\n", err)
		switch err {
		case ErrResource:
			http.Error(w, "remote resource error", 500)
		default:
			http.Error(w, "internal error", 500)
		}
		return
	}
	w.Write([]byte("all is OK"))
}

func handlerIs(w http.ResponseWriter, r *http.Request) {
	err := getRemoteResourceAndSomeWork()
	if errors.Is(err, ErrResource) {
		http.Error(w, "remote resource error", 500)
		return
	}
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}

	w.Write([]byte("all is OK"))
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/is", handlerIs)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
