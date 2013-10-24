package main

import (
	"fmt"
	. "./scrabble"
	"./solver"
	"strings"
	"strconv"
	"os"
	"bufio"
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

	test3(s)

	in := bufio.NewReader(os.Stdin)

	for {
		line, _, _ := in.ReadLine()
		input := string(line)
		move := MoveFromInput(input, s)
		fmt.Print("Move: ")
		PrintMove(move)
		points, _ := b.MakeMove(move)
		fmt.Print(" Results in Points: ")
		fmt.Println(points)
		PrintBoard(b)
	}
}

type tokenState int
const (
	placeMoveAcross tokenState = iota
	placeMoveDown
	makeMove
)

func MoveFromInput(input string, s solver.Solver) Move {
	tokens := strings.Split(input, " ")
	var state tokenState
	x, y := 0, 0
	var chips []Chip
	for index, token := range tokens {
		switch {
		case index == 0 && token == ":":
			state = makeMove
		case index == 0 && (token == "A" || token == "a"):
			state = placeMoveAcross
		case index == 0 && (token == "D" || token == "d"):
			state = placeMoveDown
		case index != 0 && state == makeMove:
			if token == "_" { token = " "; }
			c := NewChip().Init(rune(token[0]), 0, 0)
			chips = append(chips, c)
		case (index == 1 || index == 2) && (state == placeMoveAcross || state == placeMoveDown):
			val, _ := strconv.Atoi(token)
			if index == 1 { x = val; } else if index == 2 { y = val; }
		case index > 2 && state == placeMoveAcross:
			c := NewChip().Init(rune(token[0]), x, y)
			chips = append(chips, c)
			x += 1
		case index > 2 && state == placeMoveDown:
			c := NewChip().Init(rune(token[0]), x, y)
			chips = append(chips, c)
			y += 1
		}
	}
	var move Move
	if state == makeMove {
		move = s.SolveForChips(chips)
	} else {
		move = NewMove().InitWithChips(chips)
	}

	return move
}

func PrintMove(move Move) {
	if (move != nil) {
		for c := range move.Iter() {
			r := c.PlacedRune()
			if r != c.ChipRune() {
				fmt.Printf("(%s)", string(c.PlacedRune()))
			} else {
				fmt.Print(string(r))
			}
		}
	}
}

func PrintBoard(b Board) {
	for _, row := range b.SliceRep() {
		for _, item := range row {
			switch item.(type) {
			case string:
				fmt.Print(" _ ")
			case map[string]interface{}:
				m := item.(map[string]interface{})
				r := m["Rune"].(rune)
				p := m["Placed"].(rune)
				if r == p {
					fmt.Printf(" %s ",string(r))
				} else {
					fmt.Printf("(%s)",string(p))
				}
			}
		}
		fmt.Println("")
	}
}

func test1(s solver.Solver) {
	fmt.Println("Given Chips: FORESTRY")
	resultMove := s.SolveForRunes([]rune("FORESTRY"))
	fmt.Println("Resulting Move:")
	fmt.Print("\t")
	if (resultMove != nil) {
		for c := range resultMove.Iter() {
			fmt.Print(string(c.ChipRune()))
		}
	}
	fmt.Println("")
}

func test2(s solver.Solver) {
	fmt.Println("Given Chips: F O R E S T R _")
	resultMove := s.SolveForRunes([]rune("FORESTR "))
	fmt.Println("Resulting Move:")
	fmt.Print("\t")
	PrintMove(resultMove)
	fmt.Println("")
}

func test3(s solver.Solver) {
	fmt.Println("Given Chips: _ O R E S T R _")
	resultMove := s.SolveForRunes([]rune(" ORESTR "))
	fmt.Println("Resulting Move:")
	fmt.Print("\t")
	PrintMove(resultMove)
	fmt.Println("")
}