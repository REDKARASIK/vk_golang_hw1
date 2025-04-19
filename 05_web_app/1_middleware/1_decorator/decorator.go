package main

import "fmt"

type printer func(string)

func slashLineDecorator(in printer) printer {
	return func(str string) {
		fmt.Println("//////// start ////////")
		in(str)
		fmt.Println("//////// end ////////")
	}
}

func hyphenLineDecorator(in printer) printer {
	return func(str string) {
		fmt.Println("------- start -------")
		in(str)
		fmt.Println("------- end -------")
	}
}

func main() {
	var basePrinter printer = func(str string) {
		fmt.Println("hello from printer:", str)
	}

	basePrinter("base printer")
	fmt.Println()

	printer1 := slashLineDecorator(basePrinter)
	printer1("base printer + slashLineDecorator")
	fmt.Println()

	printer1 = hyphenLineDecorator(printer1)
	printer1("base printer + slashLineDecorator + hyphenLineDecorator")
	fmt.Println()

	printer2 := hyphenLineDecorator(basePrinter)
	printer2 = slashLineDecorator(printer2)
	printer2("base printer + hyphenLineDecorator + slashLineDecorator")
	fmt.Println()
}
