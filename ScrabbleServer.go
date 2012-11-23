package main

import (
	"fmt"
	. "./scrabble"
	"./solver"
)

func main() {
	fmt.Println("Creating Borad")
	b := NewBoard().Init()

	filename := "dictionary.txt"

	fmt.Println("Creating Solver")
	s, success := solver.NewSolver().Init(b, filename)
	if !success { fmt.Println("Dictionary Creation Failed"); }

	test1(s)

	test2(s)
}

func test1(s solver.Solver) {
	fmt.Println("Given Chips: PENIS")
	firstChips := make([]Chip, 5)
	for i, r := range []rune{'P', 'E', 'N', 'I', 'S'} {
		firstChips[i] = NewChip().Init(r, 0, 0)
	}
	resultMove := s.SolveForChips(firstChips)
	fmt.Println("Resulting Move:")
	fmt.Print("\t")
	for c := range resultMove.Iter() {
		fmt.Print(string(c.Rune()))
	}
	fmt.Println("")
}

func test2(s solver.Solver) {
	fmt.Println("Given Chips: P E N _ _")
	firstChips := make([]Chip, 5)
	for i, r := range []rune{'P', 'E', 'N', ' ', ' '} {
		firstChips[i] = NewChip().Init(r, 0, 0)
	}
	resultMove := s.SolveForChips(firstChips)
	fmt.Println("Resulting Move:")
	fmt.Print("\t")
	for c := range resultMove.Iter() {
		r := c.Rune()
		if r == ' ' { r = '_'; }
		fmt.Print(string(r))
	}
	fmt.Println("")
}