package main

import (
	"fmt"
)

/*
	go run -gcflags -m main.go
	go run -gcflags '-m -m' main.go
*/

type User struct {
	ID    int
	Login string
}

func (u *User) GetID() int {
	return u.ID
}

var test = ""

func newUser(login string) *User {
	return &User{123, login}
}

func setToZero(in *int) {
	for i := 0; i < 3; i++ {
		*in = 1
	}
	// *in = 0
}

type Config struct {
	Tmp1 int
}

func tmp(c *Config) {}

func main() {
	cfg := &Config{}

	u := newUser("test")
	u.ID = 1

	tmp(cfg)
	println(u.GetID())

	data := make([]string, 0, 20)
	data = append(data, "test")

	i := 1
	setToZero(&i)

	_ = fmt.Sprint(data)
	// _ = fmt.Sprint(u)
	fmt.Println("test")
}
