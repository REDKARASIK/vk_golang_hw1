package main

import "fmt"

func main() {
	// Объявление переменной через var.
	// Используется значение по умолчанию
	var num0 int

	// Объявление переменной через var.
	// Указание значения при инициализации
	var num1 int = 1

	// Пропуск типа при объявлении через var с указанием значения
	var num2 = 20

	fmt.Println(num0, num1, num2)

	// Короткое объявление переменной
	num := 30

	// Только для новых переменных
	// no new variables on left side of :=
	// num := 31 // <- Здесь будет ошибка

	num += 1
	fmt.Println("+=", num)

	// ++num нету
	num++
	fmt.Println("++", num)

	// camelCase - принятый стиль
	userIndex := 10
	// under_score - не принято
	user_index := 10
	fmt.Println(userIndex, user_index)

	// объявление нескольких переменных
	var weight, height = 10, 20

	// присваивание в существующие переменные
	weight, height = 11, 21

	// короткое присваивание
	// хотя-бы одна переменная должна быть новой!
	weight, age := 12, 22

	fmt.Println(weight, height, age)
}
