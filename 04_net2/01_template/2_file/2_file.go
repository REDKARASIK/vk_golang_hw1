package main

import (
	"fmt"
	"html/template"
	"net/http"
	// "text/template"
)

type User struct {
	ID     int
	Name   string
	Active bool
}

func main() {
	tmpl := template.Must(template.ParseFiles("users.html"))

	users := []User{
		User{1, "Anton", true},
		// xss
		User{2, "<script>alert(document.cookie)</script>", false},
		User{3, "Nikita", true},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w,
			struct {
				Users []User
			}{
				users,
			})
	})

	fmt.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
