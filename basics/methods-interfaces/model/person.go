package model

import "fmt"

type Person struct {
	Name string
	Age  int
}

func NewPerson(name string, age int) Person {
	return Person{name, age}
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
}

func PrintPerson() {

	fmt.Println(" ============== Person =========")

	giba := Person{"Gilberto", 40}
	juli := Person{"Juli", 37}
	fmt.Println(giba, juli)
}

// Println formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}
