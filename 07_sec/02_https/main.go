package main

import "net/http"

func main() {
	http.HandleFunc("/some_api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("123"))
	})
	// http.ListenAndServe(":80", nil)
	http.ListenAndServeTLS(":80", "localhost.crt", "localhost.key", nil)
}

//* PLAIN
// sudo tcpdump -i lo0 'port 80'
//
// curl http://localhost:80/some_api

//* SSL
// openssl req -new -subj "/C=RU/ST=Msk/CN=localhost" -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr
// openssl x509 -req -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt
//
// sudo tcpdump -i lo0 'port 80'
//
// curl -k https://localhost:80/some_api
