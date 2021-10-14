package main

import (
	"bytes"
	"example/methods-interfaces/model"
	"fmt"
	"math"
	"strconv"
)

type Vertex struct {
	X, Y float64
}

type MyFloat float64

/*
Um tipo implementa uma interface através da implementação dos métodos. Não há declaração explícita de intenções, não há palavra-chave "implements".
Interfaces implícitas dissociam-se da definição de uma interface de sua implementação, que pode então aparecer em qualquer pacote sem pré-arranjamento.
*/
type Abser interface {
	Abs() float64
}

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

func (v Vertex) ToString() string {
	var b bytes.Buffer
	b.WriteString("Vertex {\n  X:")
	b.WriteString(strconv.FormatFloat(v.X, 'f', -1, 64))
	b.WriteString("\n  Y:")
	b.WriteString(strconv.FormatFloat(v.Y, 'f', -1, 64))
	b.WriteString("\n}")
	return b.String()
}

/* Ponteiros receptores pemitino a alteração
interna dos valores da struct ao passar por parâmetros
*/
func (v *Vertex) SetValues(x float64, y float64) {
	v.X = x
	v.Y = y
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

func interfaceTest() {
	fmt.Println("====== TESTE COM INTERFACE =====")

	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := Vertex{3, 4}

	fmt.Println("F=> ", f.Abs())
	fmt.Println("V=> ", v.ToString())
	a = f
	fmt.Println("A_F=> ", a.Abs())

	a = &v
	fmt.Println("A_&V=> ", a.Abs())

	a = v //V não implementa Abs, apenas sua referência tem o método Abs
	fmt.Println("A_V=> ", a.Abs())

}

func main() {
	//Teste 1
	v := Vertex{3.5, 4.8}
	fmt.Println(" Valor ABS ", v.Abs()) //Abs está associado a Vertex
	fmt.Println(" Valor ABS ", Abs(v))  //Abs está associado a Vertex

	//Teste 2
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs())

	//Teste3
	fmt.Println("\nTeste A:\n", v.ToString())

	v.SetValues(50.6, 89.5)
	fmt.Println("\nTeste B:\n", v.ToString())

	interfaceTest()

	model.Println("Person => ", v)

	model.PrintPerson()
	p := model.NewPerson("A", 50)
	fmt.Println("Person => ", p)
}
