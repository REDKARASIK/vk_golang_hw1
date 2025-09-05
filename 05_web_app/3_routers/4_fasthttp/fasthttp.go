package main

import (
	"encoding/json"
	"fmt"
	"log"

	fasthttprouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// curl -v http://localhost:8080/
func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	users := []string{"rvasily", "bob"}
	body, _ := json.Marshal(users)

	ctx.Write([]byte("\n\tUsers list: "))
	ctx.Write(body)
	ctx.Write([]byte("\n\n"))
}

// curl -v http://localhost:8080/users/100500
func GetUser(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "\n\tYou try to see user %q\n\n", ctx.UserValue("id"))
}

func main() {
	router := fasthttprouter.New()

	router.GET("/", Index)
	router.GET("/users/{id}", GetUser)

	fmt.Println("starting server at :8080")
	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
