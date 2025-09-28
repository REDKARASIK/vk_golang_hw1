package main

import "fmt"

func main() {
	// Инициализация при создании
	var user map[string]string = map[string]string{
		"name":     "Anton",
		"lastName": "Сhumakov",
	}

	// Создание сразу с нужной ёмкостью
	profile := make(map[string]string, 10)

	// Количество элементов
	mapLength := len(user)

	fmt.Printf("%d %+v\n", mapLength, profile)

	// Если ключа нет - вернёт значение по умолчанию для типа
	mName := user["middleName"]
	fmt.Println("mName:", mName)

	// Проверка на существование ключа
	mName, mNameExist := user["middleName"]
	fmt.Println("mName:", mName, "mNameExist:", mNameExist)

	// Пустая переменная - только проверяем что ключ есть
	_, mNameExist2 := user["middleName"]
	fmt.Println("mNameExist2", mNameExist2)

	// Удаление ключа
	delete(user, "lastName")
	fmt.Printf("%#v\n", user)
}
