package main

import "fmt"

// Обычное объявление
func singleIn(in int) int {
	return in
}

// Множество параметров
func multIn(a, b int, c int) int {
	return a + b + c
}

// Без результата
func withoutReturn(s string) {
	for _, sym := range s {
		fmt.Print(sym)
	}
}

// Именованный результат
func namedReturn() (out int) {
	out = 2
	return
}

// Несколько результатов
func multipleReturn(in int) (int, error) {
	if in > 2 {
		return 0, fmt.Errorf("some error happend")
	}
	return in, nil
}

// Несколько именованных результатов
func multipleNamedReturn(ok bool) (rez int, err error) {
	rez = 1
	if ok {
		err = fmt.Errorf("some error happend")
		// Аналогично return rez, err
		return
	}
	rez = 2
	return
}

// Нефиксированное количество параметров
func sum(in ...int) (result int) {
	fmt.Printf("in := %#v \n", in)
	for _, val := range in {
		result += val
	}
	return
}

func main() {
	nums := []int{1, 2, 3, 4}
	fmt.Println(nums, sum(nums...))
	return
}
