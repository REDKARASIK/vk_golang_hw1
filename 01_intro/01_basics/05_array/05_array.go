package main

import "fmt"

func main() {
	// Размер массива является частью его типа

	// Инициализация значениями по-умолчанию
	var a1 [3]int                   // [0,0,0]
	fmt.Println("a1 short", a1)     // a1 short [0 0 0]
	fmt.Printf("a1 short %v\n", a1) // a1 short [0 0 0]
	fmt.Printf("a1 full %#v\n", a1) // a1 full [3]int{0, 0, 0}

	const size = 2
	var a2 [2 * size]bool // [false,false,false,false]
	fmt.Println("a2", a2)

	// Определение размера при объявлении
	a3 := [...]int{1, 2, 3}
	fmt.Println("a2", a3)

	// Проверка при компиляции или при выполнении
	// invalid array index 4 (out of bounds for 3-element array)
	// a3[idx] = 12
}
