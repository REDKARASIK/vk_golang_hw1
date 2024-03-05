package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type tplParams struct {
	URL     string
	Browser string
}

const EXAMPLE = `
Browser {{.Browser}}

you at {{.URL}}
`

// var someHtml = "<html><body>Hello, <b>{{.Name}}</b>!</body></html>"

var tmpl = template.New("123")

func handle(w http.ResponseWriter, r *http.Request) {
	params := tplParams{
		URL:     r.URL.String(),
		Browser: r.UserAgent(),
	}

	tmpl.Execute(w, params)
}

func main() {
	tmpl, _ = tmpl.Parse(EXAMPLE)

	http.HandleFunc("/", handle)

	fmt.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
