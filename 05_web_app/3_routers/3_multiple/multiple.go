package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
)

func FastRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, "\n\tRequest with high hit rate\n\n")
}

func ComplexRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "\n\tRequest with complex routing logic\n\n")
}

func RegularRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "\n\tRequest with average amount of logic\n\n")
}

func main() {
	// curl -v http://localhost:8080/fast/123
	fastApiHandler := httprouter.New()
	fastApiHandler.GET("/fast/:id", FastRequest)

	// curl -v -H "X-Requested-With: XMLHttpRequest" http://localhost:8080/complex/
	complexApiHandler := mux.NewRouter()
	complexApiHandler.HandleFunc("/complex/", ComplexRequest).
		Headers("X-Requested-With", "XMLHttpRequest") // ajax

	// curl -v http://localhost:8080/std/smth
	stdApiHandler := http.NewServeMux()
	stdApiHandler.HandleFunc("/std/", RegularRequest)

	siteMux := http.NewServeMux()
	siteMux.Handle("/fast/", fastApiHandler)
	siteMux.Handle("/complex/", complexApiHandler)
	siteMux.Handle("/std/", stdApiHandler)

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", siteMux))
}
