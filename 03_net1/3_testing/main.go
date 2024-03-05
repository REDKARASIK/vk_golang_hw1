package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID int
}

var data = map[string][]byte{
	"ok":   []byte(`{"ID": 27}`),
	"fail": []byte(`{"ID": 27`),
}

func GetUser(key string) (*User, error) {
	jsonStr, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("user doesnt exist")
	}

	user := &User{}
	err := json.Unmarshal(jsonStr, user)
	if err != nil {
		return nil, fmt.Errorf("cant decode json")
	}

	return user, nil
}
