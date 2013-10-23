package solver

import (
	. "../scrabble"
	"../dictionary"
)

// an interface for a solver
type Solver interface {
	// initialize the solver with the given board and dictionary file name
	// returns the solver and a bool flag for file success
	Init(board Board, dictionaryFilename string) (Solver, bool)

	// find the next best move for the given set of chips
	// and all the available chips on the baord
	SolveForChips(chips []Chip) Move
}

// data for a solver
type solver struct {
	board 	Board 		// the board for the solver
	dict 	dictionary.Dictionary 	// the dictionary for the solver
}

// create a new solver object
func NewSolver() *solver {
	s := new(solver)
	return s
}

func (s *solver) Init(board Board, dictionaryFilename string) (Solver, bool) {
	s.board = board
	var success bool
	s.dict, success = dictionary.NewDictionary().Init(dictionaryFilename)
	return s, success
}

func (s *solver) SolveForChips(chips []Chip) Move {
	// create the rune slice
	runes := make([]rune, len(chips))
	for index, c := range chips {
		runes[index] = c.ChipRune()
	}

	// all possible solutions
	moves := make([]Move, 0)

	// function to generate moves for the given runes around the pivot chip
	// the pivot chip's rune must be in the rues slice
	genMoves := func(runes []rune, pivotChip Chip) []Move {
		for possible := range s.dict.FindAllRunes(runes) {
			moves = append(moves, s.createMoves(possible, pivotChip)...)
		}
		return moves
	}

	// for all the available pieces
	for availablePiece := range s.availableMoveSpaces() {
		if availablePiece != nil { // if there are pieces

			runes = append(runes, availablePiece.PlacedRune())
			moves = genMoves(runes, availablePiece)

		} else { // if there are no pieces
			// make a potential pivot piece from every rune
			x, y := s.board.Center()
			for _, r := range runes {
				chip := NewChip().Init(r, x, y)
				moves = genMoves(runes, chip)
			}
		}
	}

	return s.findBestMove(moves)
}

// a private enumeration for directionality of moves
type move_direction int
const (
	none	move_direction = iota
	down
	across
)

// get the direction that the move around the given chip should be
func (s *solver) moveDirection(pivotChip Chip) move_direction {
	chipX, chipY := pivotChip.Position()
	above 	:= s.board.ChipAtSpace(chipX, chipY-1)
	below 	:= s.board.ChipAtSpace(chipX, chipY+1)
	left 	:= s.board.ChipAtSpace(chipX-1, chipY)
	right 	:= s.board.ChipAtSpace(chipX+1, chipY)
	if above == nil && below == nil {
		return down
	} else if left == nil && right == nil {
		return across
	}
	return none
}

// create a slics of chips with positions around the supplied pivot chip
// head is a collection of runes that preceed the pivot chip
// tail it a collection of runes that follow the pivot chip
// the pivot chip has a defined position and a rune
func (s *solver) chipSliceAroundPiviot(head []dictionary.RunePair, pivotChip Chip, tail []dictionary.RunePair) []Chip {
	var chips []Chip

	x, y := pivotChip.Position()
	switch s.moveDirection(pivotChip) {
	case down:
		for index, runepair := range head {
			newChip := NewChip().InitPlaced(runepair.Word, runepair.Combination, x, y-(len(head)-index))
			chips = append(chips, newChip)
		}
		chips = append(chips, pivotChip)
		for index, runepair := range tail {
			newChip := NewChip().InitPlaced(runepair.Word, runepair.Combination, x, y+1+index)
			chips = append(chips, newChip)
		}
	case across:
		for index, runepair := range head {
			newChip := NewChip().InitPlaced(runepair.Word, runepair.Combination, x-(len(head)-index), y)
			chips = append(chips, newChip)
		}
		chips = append(chips, pivotChip)
		for index, runepair := range tail {
			newChip := NewChip().InitPlaced(runepair.Word, runepair.Combination, x+1+index, y)
			chips = append(chips, newChip)
		}
	}

	return chips
}

// create a move using the given runes
// the pivot chip has a rune that exists somewhere in the runes slice
func (s *solver) createMoves(word dictionary.RuneResult, pivotChip Chip) []Move {
	moves := make([]Move, 0)

	runepairs := dictionary.UnzipRuneResult(word)

	for index, runepair := range runepairs {
		if runepair.Word == pivotChip.PlacedRune() {
			chips := s.chipSliceAroundPiviot(runepairs[:index], pivotChip, runepairs[index+1:])
			newMove := NewMove().InitWithChips(chips)
			if s.board.IsMoveValid(newMove) {
				moves = append(moves, newMove)
			}
		}
	}
	return moves
}

// if this is the first move of the game, nil is passes through the chan
// else, all chips that are available for in-play are passed
func (s *solver) availableMoveSpaces() <-chan Chip {
	spaces := make(chan Chip, 10)

	go func() {
		counter := 0
		for chip := range s.board.PlacedChipIter() {
			switch s.moveDirection(chip) {
			case down, across:
				spaces <- chip
			}
			counter += 1
		}
		if counter == 0 { spaces <- nil; }
		close(spaces)
	}()

	return spaces
}

// find the best move, given a slice of available moves
// returns either the best move or nil if no moves are possible
func (s *solver) findBestMove(moves []Move) Move {
	var bestMove Move
	bestPoints := 0
	for _, m := range moves {
		points := s.board.PointsForMove(m)
		if points > bestPoints {
			bestMove = m
			bestPoints = points
		}
	}
	return bestMove
}