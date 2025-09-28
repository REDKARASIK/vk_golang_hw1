package main

import "fmt"

// Одиночная константа
const pi = 3.141

// Блок констант
const (
	hello = "Привет"
	e     = 2.718
)

// iota. Базовый пример
const (
	zero = iota
	_    // Пустая переменная, пропуск iota
	two
	three // = 3
)

// iota. Продвинутый пример
const (
	_         = iota             // Пропускаем первое значение
	KB uint64 = 1 << (10 * iota) // 1 << (10 * 1) = 1024
	MB                           // 1 << (10 * 2) = 1048576
)

const (
	// Нетипизированная константа
	year = 2017
	// Типизированная константа
	yearTyped int = 2017
)

func main() {
	var month int32 = 13
	fmt.Println(month + year)

	// month + yearTyped // <- Ошибка mismatched types int32 and int
	// fmt.Println( month + yearTyped )
}
