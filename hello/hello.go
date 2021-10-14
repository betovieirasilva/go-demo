package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	//formula => z - = (z * z - x) / (2 * z)
	z := float64(1) //or 1.0 para forçar a entrada como float64
	raizProxima := float64(0)
	//testa a xpoximação de z2 == x
	for i := 0; i < 15; i++ {
		z -= (z*z - x) / (2 * z)
		fmt.Println("\tZ", z)
		raizProxima = z
		if (z * 2) <= x {
			break
		}

	}
	return raizProxima

}

func main() {
	//fmt.Println(Sqrt(2))

	//obtém a raiz quadrada de 0..20
	for i := 1; i < 21; i++ {
		valorBase := float64(i)
		raizQuadrada := math.Sqrt(valorBase) //função oficial
		fmt.Println("Raiz de ", i, " é => ", raizQuadrada, " | func => ", Sqrt(valorBase))
	}

}
