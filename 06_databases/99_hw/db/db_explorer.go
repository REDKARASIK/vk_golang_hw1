package main

import (
	"database/sql"
	"net/http"
)

func NewDBExplorer(_ *sql.DB) (http.Handler, error) {
	// тут вы пишете код
	// обращаю ваше внимание - в этом задании запрещены глобальные переменные
	return nil, nil
}
