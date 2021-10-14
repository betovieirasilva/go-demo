package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
}

func printPerson() {

	fmt.Println(" ============== Person =========")

	giba := Person{"Gilberto", 40}
	juli := Person{"Juli", 37}
	fmt.Println(giba, juli)
}
