package main

import "fmt"

type Person struct {
	Id      int
	Name    string
	Address string
}

type Account struct {
	Id      int
	Cleaner func(string) string
	Owner   Person
	Person
	//Name    string
}

func main() {
	// Полное объявление структуры
	var acc Account = Account{
		Id: 1,
		Person: Person{
			Name:    "Антон",
			Address: "Москва",
		},
		//Name: "anton",
	}
	fmt.Printf("%#v\n", acc)

	// Короткое объявление структуры
	acc.Owner = Person{2, "Chumakov Anton", "Moscow"}

	fmt.Printf("%#v\n", acc)

	fmt.Println(acc.Name)
	fmt.Println(acc.Person.Name)
}
