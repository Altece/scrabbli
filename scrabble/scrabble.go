package scrabble

// an interface for interacting with a playing chip used in a game of scrabble
type Chip interface {
	// initialize the chip with the given rune and x y coordinate positions
	Init(r rune, x, y int) Chip

	// get thr rune from the chip
	Rune() rune

	// get the position of the chip
	Position() (int, int)

	// get the points for the chip's rune
	Points() int
}

// an interface for interacting with a move for a scrabble board
type Move interface {
	// initialize the move
	Init() Move

	// initialize the move with the given chips
	InitWithChips(chips []Chip) Move

	// get an iterator channel to iterate over the chips in the move
	Iter() <-chan Chip

	// add a chip to the move
	Add(chip Chip)
}

// an interface for interacting with a scrabble board
type Board interface {
	// initialize the board
	Init() Board

	// get the rune at the given x y coordinates
	ChipAtSpace(x, y int) Chip

	// checks to see if the given move is valid for the game board
	IsMoveValid(m Move) bool

	// get the points that would be gained by placing a rune at the specified position
	PointsForMove(m Move) int

	// apply the given move to the game board
	// returns an int and a bool
	// the int is the amount of points earned by the move
	// the bool is true if successful, false otherwise
	// if the move is not successful, the int value is 0
	MakeMove(m Move) (int, bool)
}

func WorkingMessage() string {
	return "Scrabble Working!"
}