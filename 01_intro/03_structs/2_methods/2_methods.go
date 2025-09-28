package main

import "fmt"

type Person struct {
	Id   int
	Name string
}

// Не изменит оригинальной структуры, для который вызван
func (p Person) UpdateName(name string) {
	p.Name = name
}

// Изменяет оригинальную структуру
func (p *Person) SetName(name string) {
	p.Name = name
}

type Account struct {
	Id   int
	Name string
	Person
}

// А если этого метода нет?
func (p *Account) SetName(name string) {
	p.Name = name
}

type MySlice []int

func (sl *MySlice) Add(val int) {
	*sl = append(*sl, val)
}

func (sl *MySlice) Count() int {
	return len(*sl)
}

func main() {
	// Вариант 1 - Создать объект с заданными вручную полями
	// pers := &Person{1, "Vasily"}

	// Вариант 2 - Создать пустой объект и задать поля через методы
	pers := new(Person)
	fmt.Printf("empty person: %#v\n", pers)
	pers.SetName("Vasily Romanov")
	fmt.Printf("updated person: %#v\n", pers)

	acc := Account{
		Id:   1,
		Name: "vasya",
		Person: Person{
			Id:   2,
			Name: "Vasily Romanov",
		},
	}

	acc.SetName("romanov.vasily")
	acc.Person.SetName("Test")

	fmt.Printf("%#v \n", acc)

	sl := MySlice([]int{1, 2})
	sl.Add(5)
	fmt.Println(sl.Count(), sl)
}
