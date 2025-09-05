package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Подборка лучших роутеров для go:
// https://github.com/avelino/awesome-go?tab=readme-ov-file#routers

// curl -v -H "Content-Type: application/json" http://localhost:8080/
func UserList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "\n\tYou see user list\n\n")
}

// curl -v -H "Content-Type: application/json" --resolve qq.vk.ru:8080:127.0.0.1 http://qq.vk.ru:8080/users
func UserListVK(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "\n\tYou see user list for subdomain %q for host %q\n\n", vars["subdomain"], r.Host)
}

// curl -v -X PUT -H "Content-Type: application/json" http://localhost:8080/users/rvasily
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "\n\tYou try to update user %q\n\n", vars["login"])
}

// curl -v -H "Content-Type: application/json" http://localhost:8080/users/100500
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "\n\tYou try to see user with id=%q\n\n", vars["id"])
}

// curl -v -X POST -H "Content-Type: application/json" -H "X-Auth: test" http://localhost:8080/users/rvasily
func CreateTestUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "\n\tYou try to create new user %q\n\n", vars["login"])
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", UserList)

	r.HandleFunc("/users", UserList).Host("localhost")
	r.HandleFunc("/users", UserListVK).Host("{subdomain}.vk.ru")

	// See anything wrong?:
	r.HandleFunc("/users/{login:[0-9a-z]+}", UpdateUser)
	r.HandleFunc("/users/{id:[0-9]+}", GetUser)
	r.HandleFunc("/users/{login}", CreateTestUser).Methods("POST").Headers("X-Auth", "test")

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
