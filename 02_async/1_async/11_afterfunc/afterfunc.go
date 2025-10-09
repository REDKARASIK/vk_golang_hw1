package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello World")
}

func main() {
	timer := time.AfterFunc(1*time.Second, sayHello)

	fmt.Scanln()
	// Если успеем остановить таймер - функция не выполнится.
	timer.Stop()

	fmt.Scanln()
}
