// Тут менять вам может потребоваться только коннект к базе
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// DSN это соединение с базой
	// Вы можете изменить его на то, которое вам нужно
	// docker run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
	DSN = "root:1234@tcp(localhost:3306)/golang?charset=utf8"
)

func main() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	err = db.Ping() // Тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	handler, err := NewDBExplorer(db) //nolint:typecheck
	if err != nil {
		panic(err)
	}

	fmt.Println("starting server at :8082")
	if err := http.ListenAndServe(":8082", handler); err != nil {
		log.Printf("error listenAndServer: %v", err)
	}
}
