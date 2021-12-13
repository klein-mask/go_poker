package main

import (
	"go_poker/poker"
	"os"
	"strconv"
)

func main() {
	bigBlind := 200
	smallBilnd := 100
	playerInitMoney := 3000
	if len(os.Args) > 3 {
		bigBlind, _ = strconv.Atoi(os.Args[1])
		smallBilnd, _ = strconv.Atoi(os.Args[2])
		playerInitMoney, _ = strconv.Atoi(os.Args[3])
	}

	p := poker.NewPoker(bigBlind, smallBilnd, playerInitMoney)
	p.InitSetUp()
}
