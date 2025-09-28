package main

// Объявление нового типа
type UserID int

// Объявление псевдонима (алиаса)
type Temperature = float32

func main() {
	idx := 1
	var uid UserID = 42

	// Даже если базовый тип одинаковый, разные типы несовместимы
	// cannot use idx (type int) as type UserID in assignment
	// uid := idx

	println(uid, idx)
	uid = UserID(idx)
	println(uid, idx)

	// ----------------------------------------------------------

	illManTemperature := float32(37.8)

	var manTemperature Temperature = 36.6 // Надо ли явно указывать тип, если Temperature будет алиасом на float64?

	println(manTemperature, illManTemperature)
	manTemperature = illManTemperature
	println(manTemperature, illManTemperature)
}
