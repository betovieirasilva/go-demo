package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

type MyFloat float64

/**
Go não tem classes. No entanto, você pode definir métodos em tipos.
O método é uma função com um argumento receptor especial.
O receptor aparece em sua lista de argumentos entre a própria palavra-chave func e o nome do método.
Neste exemplo, o método Abs tem um receptor do tipo Vertex chamado v.
*/
//este é um método
func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.Y + v.Y*v.Y)
}

//Aqui está Abs escrito como uma função regular, sem qualquer alteração na funcionalidade.
//esta é uma function
func Abs(v Vertex) float64 {
	return math.Sqrt(v.X*v.Y + v.Y*v.Y)
}

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func main() {
	//Teste 1
	v := Vertex{3, 4}
	fmt.Println(" Valor ABS ", v.Abs()) //Abs está associado a Vertex
	fmt.Println(" Valor ABS ", Abs(v))  //Abs está associado a Vertex

	//Teste 2
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())
}
