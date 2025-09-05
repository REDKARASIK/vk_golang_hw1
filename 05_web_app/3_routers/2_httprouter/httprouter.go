package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// curl -v -H "Content-Type: application/json" http://localhost:8080/
func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { // NOTE: 3 params for handler
	fmt.Fprint(w, "\n\tYou see user list\n\n")
}

// curl -v -H "Content-Type: application/json" http://localhost:8080/users/rvasily
func GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "\n\tYou try to see user %q\n\n", ps.ByName("id"))
}

// curl -v -X PUT -H "Content-Type: application/json" http://localhost:8080/users
func CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "\n\tYou try to create new user\n\n")
}

// curl -v -X POST -H "Content-Type: application/json" http://localhost:8080/users/rvasily
func UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "\n\tYou try to update %q\n\n", ps.ByName("login"))
}

// curl -v -X POST -H "Content-Type: application/json" http://localhost:8080/users/rvasily
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	fmt.Fprintf(w, "\n\tYou try to update %q with std handler\n\n", params.ByName("login"))
}

func main() {
	router := httprouter.New()

	router.GET("/", List)
	router.GET("/users", List)
	router.PUT("/users", CreateUser)
	router.GET("/users/:id", GetUser)

	router.POST("/users/:login", UpdateUser)
	// router.HandlerFunc(http.MethodPost, "/users/:login", UpdateUserHandler)

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
