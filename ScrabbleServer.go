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
	fmt.Println("Given Chips: A B C D E")
	firstChips := make([]Chip, 5)
	for i, r := range []rune{'A', 'B', 'C', 'D', 'E'} {
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
	fmt.Println("Given Chips: A B _ D E")
	firstChips := make([]Chip, 5)
	for i, r := range []rune{'A', 'B', ' ', 'D', 'E'} {
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