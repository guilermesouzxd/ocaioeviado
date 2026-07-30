package main

import "fmt"

func main() {

	var tamAr int
	var impr []string
	var cont int
	var contA int
	var backup int

	fmt.Scan(&tamAr)
	backup = tamAr
	impr = make([]string, tamAr)

	for cont = 0; cont < tamAr; cont++ {
		impr[cont] = "*"
	}

	for cont = 0; cont < backup; cont++ {
		fmt.Println()
		for contA = 0; contA < tamAr; contA++ {
			fmt.Print(impr[contA])
		}
		tamAr--
	}

	fmt.Println(" Nao odeio negros e judeus")
}
