package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func main() {
	serverSideServer := &http.Server{Addr: ":8081"}
	serverSideServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "SECRET_CREDENTIALS")
	})
	go serverSideServer.ListenAndServe()

	clientSideServer := &http.Server{Addr: ":8080"}
	clientSideServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urls, ok := r.URL.Query()["url"]
		if !ok || len(urls) != 1 {
			http.Error(w, "bad url param", 500)
			return
		}

		resp, err := http.Get(urls[0])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Write out the hexdump of the bytes as plaintext.
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, hex.Dump(bytes))
	})
	clientSideServer.ListenAndServe()
}

// http://localhost:8080?url=http://www.google.com
// http://localhost:8080?url=http://localhost:8081
