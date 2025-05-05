package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var privateAlbums = map[string]bool{
	"selfies": true,
	"secret": true,
}

var privateAlbumsACL = map[string]map[int]bool{
	"selfies": {
		1: true,
		2: true,
		3: true,
	},
	"secret": {
		2: true,
	},
}

func isAllowedByACL(albumName string, uid int) bool {
	if !privateAlbums[albumName] {
		log.Println(albumName, "not private")
		return true
	}
	allowedUsers, ok := privateAlbumsACL[albumName]
	if !ok {
		log.Println(albumName, "not in privateAlbumsACL")
		return false
	}
	return allowedUsers[uid]
}

func getSession(r string) int {
	uid, _ := strconv.Atoi(r)
	return uid
}

func main() {

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AUTH request: %+v", r)

		req, _ := url.Parse(r.Header.Get("X-Original-URI"))
		uid := getSession(req.Query().Get("user_id"))

		str := req.Path
		albumName := strings.ReplaceAll(str, "/albums/", "")
		log.Println("PARAMS", albumName, uid)

		if !isAllowedByACL(albumName, uid) {
			log.Println("ACL failed", albumName, uid)
			http.Error(w, "", http.StatusForbidden)
		}

		log.Println("ACL OK", albumName, uid)
	})

	http.HandleFunc("/albums/{album}", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("incoming request: %+v", r)
		w.Write(fmt.Appendf(nil,
			"hi, %s. You are in album %s",
			r.Header.Get("X-User-ID"),
			r.PathValue("album")))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ROOT incoming request: %+v", r)
		user := r.Header.Get("X-User-ID")
		if user == "" {
			user = "anonimous"
		}
		w.Write(fmt.Appendf(nil, "hi, %s", user))
	})

	log.Println("start server at :8080")
	http.ListenAndServe(":8080", nil)
}
