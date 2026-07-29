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

	for cont = 0; cont < tamAr; cont++ {
		for contA = backup - 1; contA >= 0; contA++ {
			fmt.Print(impr[contA])
		}
		backup--
		fmt.Println()
	}

	fmt.Println(" Nao odeio negros e judeus")
}
