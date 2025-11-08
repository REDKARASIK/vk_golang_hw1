package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var datasetPath = "dataset.xml"
var accessToken = "test"

type XMLData struct {
	Version  string `xml:"version,attr"`
	UserData []row  `xml:"row"`
}

type row struct {
	ID        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("AccessToken")
	if token != accessToken {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"Error":"invalid access token"}`))
		return
	}

	data, err := os.ReadFile(datasetPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Error":"error openning dataset file"}`))
		return
	}
	var xmlData XMLData
	if err := xml.Unmarshal(data, &xmlData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Error":"error parsing xml from dataset"}`))
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"Error": "unknown endpoint"}`))
		return
	}

	orderBy := r.URL.Query().Get("order_by")
	if orderBy == "" {
		orderBy = "0"
	}
	if orderBy != "0" && orderBy != "-1" && orderBy != "1" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error": "OrderBy invalid"}`))
		return
	}

	orderField := r.URL.Query().Get("order_field")
	if orderField == "" {
		orderField = "name"
	}
	if orderField != "id" && orderField != "age" && orderField != "name" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error": "OrderField invalid"}`))
		return
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error":"miss limit in Query"}`))
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error":"cast limit to int"}`))
		return
	}

	offset := r.URL.Query().Get("offset")
	var offsetInt int
	if offset == "" {
		offsetInt = 0
	} else {
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"Error":"cast offset to int"}`))
			return
		}
	}

	query := r.URL.Query().Get("query")

	users := make([]User, 0, len(xmlData.UserData))
	for _, u := range xmlData.UserData {
		user := User{
			ID:     u.ID,
			Name:   u.FirstName + " " + u.LastName,
			Age:    u.Age,
			About:  u.About,
			Gender: u.Gender,
		}
		if query == "" {
			users = append(users, user)
		} else if strings.Contains(user.Name, query) || strings.Contains(u.About, query) {
			users = append(users, user)
		}
	}

	switch orderBy {
	case "-1":
		sort.Slice(users, func(i, j int) bool {
			switch orderField {
			case "id":
				return users[i].ID > users[j].ID
			case "age":
				return users[i].Age > users[j].Age
			default:
				return users[i].Name > users[j].Name
			}
		})
	case "1":
		sort.Slice(users, func(i, j int) bool {
			switch orderField {
			case "id":
				return users[i].ID < users[j].ID
			case "age":
				return users[i].Age < users[j].Age
			default:
				return users[i].Name < users[j].Name
			}
		})
	}

	if offsetInt > len(users) {
		users = []User{}
	} else {
		users = users[offsetInt:]
		if limitInt > len(users) {
			limitInt = len(users)
		}
		users = users[:limitInt]
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Error":"encoding json failed"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
}
