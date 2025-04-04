package main

import (
	"fmt"
	"net/http"
	"text/template"
	// "html/template"
)

type tplParams struct {
	URL     string
	Browser string
}

type htmlParams struct {
	Name string
}

const EXAMPLE = `
Browser {{.Browser}}

you at {{.URL}}
`

// const SOME_HTML = "<html><body>Hello, <b>{{.Name}}</b>!</body></html>"

var tmpl = template.New("123")

func handle(w http.ResponseWriter, r *http.Request) {
	params := tplParams{
		URL:     r.URL.String(),
		Browser: r.UserAgent(),
	}

	//hParams := htmlParams{
	//	Name: "world",
	//}

	tmpl.Execute(w, params)
	// tmpl.Execute(w, hParams)
}

func main() {
	tmpl, _ = tmpl.Parse(EXAMPLE)
	// tmpl, _ = tmpl.Parse(SOME_HTML)

	http.HandleFunc("/", handle)

	fmt.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
