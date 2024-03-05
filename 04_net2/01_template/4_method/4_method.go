package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type User struct {
	ID     int
	Name   string
	Active bool
}

func (u *User) PrintActive(uppercase bool) string {
	if !u.Active {
		return ""
	}

	output := "method says user " + u.Name + " active"

	if uppercase {
		return strings.ToUpper(output)
	}

	return output
}

func main() {
	tmpl, err := template.New("").ParseFiles("method.html")
	if err != nil {
		panic(err)
	}

	users := []User{
		{1, "Anton", true},
		{2, "Nikita", false},
		{3, "Veronika", true},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "method.html",
			struct {
				Users []User
			}{
				users,
			})
		if err != nil {
			panic(err)
		}
	})

	fmt.Println("starting server at :8080")
	errStart := http.ListenAndServe(":8080", nil)
	if errStart != nil {
		panic(errStart)
	}
}
