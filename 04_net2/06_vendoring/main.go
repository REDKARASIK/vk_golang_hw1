package main

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
)

// go run -mod=vendor main.go

func main() {
	p := bluemonday.UGCPolicy()
	res := p.Sanitize("some test")
	fmt.Println(res)
}
