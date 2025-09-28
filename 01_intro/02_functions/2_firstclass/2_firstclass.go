package main

import "fmt"

// Обычная функция
func doNothing() {
	fmt.Println("i'm regular function")
}

func main() {
	// Анонимная функция
	func(in string) {
		fmt.Println("anon func out:", in)
	}("nobody")

	// Присванивание анонимной функции в переменную
	printer := func(in string) {
		fmt.Println("printer outs:", in)
	}
	printer("as variable")

	// определяем тип функции
	type strFuncType func(string)

	// функция принимает коллбек (другую функцию)
	worker := func(callback strFuncType) {
		callback("as callback")
	}
	worker(printer)

	// Функция возвращает замыкание
	prefixer := func(prefix string) strFuncType {
		return func(in string) {
			fmt.Printf("[%s] %s\n", prefix, in)
		}
	}
	successLogger := prefixer("SUCCESS")
	successLogger("expected behaviour")
}
