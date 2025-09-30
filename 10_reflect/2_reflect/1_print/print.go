package main

import (
	"fmt"
	"reflect"
)

type Login string

type User struct {
	ID       int
	RealName string `unpack:"-"`
	Login    Login
	Flags    int
}

func PrintReflect(u interface{}) error {
	val := reflect.ValueOf(u)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fmt.Println(val.Type(), val.Kind())


	fmt.Printf("%T have %d fields:\n", u, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		fmt.Printf("\tname=%v, type=%v, value=%v, tag=`%v`\n",
			typeField.Name,
			typeField.Type,
			valueField,
			typeField.Tag,
		)
	}
	return nil
}

func main() {
	u := User{
		ID:       42,
		RealName: "rvasily",
		Flags:    32,
	}
	fmt.Printf("%#v\n", u)
	err := PrintReflect(u)
	if err != nil {
		panic(err)
	}
}
