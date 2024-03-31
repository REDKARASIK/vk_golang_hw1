package main

import (
	"fmt"
)

func deferTest() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic happend FIRST:", err)
		}
	}()
	fmt.Println("Some userful work")
	panic("something bad happend")
	return
}

func main() {
	deferTest()
	println("good")
	return
}
