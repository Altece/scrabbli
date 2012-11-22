package scrabble

// a struct to represent a space on the board
type space struct {
	chip 		Chip 		// a chip occupying the space (if there is one)
	multiplier 	Multiplier 	// the multiplier for this space
}

// a struct to represent the board
type board struct {
	grid			[][]*space 	// a 2D array to represent the board
	height, width	int 		// the height and width of the board
}

func NewBoard() Board {
	return new(board)
}

// check to see if the given x y coordinates are within the bounds of the board
func (b *board) WithinBounds(x, y int) bool {
	if x >= 0 && y >= 0 {
		if x < b.width && y < b.height {
			return true
		}
	}
	return false
}

// get the multiplier for the space at the given x y coordinates
func (b *board) MultiplierForSpace(x, y int) Multiplier {
	if b.WithinBounds(x, y) {
		return b.grid[x][y].multiplier
	}
	return normal
}

func (b *board) Init() Board {
	// set up the width
	height, width := 15, 15
	b.height, b.width = height, width

	// initialize the grid
	b.grid = make([][]*space, width)
	for x, _ := range b.grid {
		b.grid[x] = make([]*space, height)
		for y, _ := range b.grid[x] {
			b.grid[x][y] = new(space)
			b.grid[x][y].chip = nil
			b.grid[x][y].multiplier = b.multiplier(x, y)
		}
	}

	return b
}

func (b *board) ChipAtSpace(x, y int) Chip {
	if b.WithinBounds(x, y) {
		return b.grid[x][y].chip
	}
	return nil
}

func (b *board) Center() (x, y int) {
	x = b.width / 2
	y = b.height / 2
	return x, y
}

func (b *board) PlacedChipIter() <-chan Chip {
	chips := make(chan Chip, 10)

	go func() {
		for x := 0; x < b.width; x+=1 {
			for y := 0; y < b.height; y+=1 {
				if c := b.grid[x][y].chip; c != nil {
					chips <- c
				}
			}
		}
		close(chips)
	}()

	return chips
}

func (b *board) IsMoveValid(m Move) bool {
	for chip := range m.Iter() {
		print("checking chip: ")
		println(chip)
		print("\t")
		print(string(chip.Rune()))
		x, y := chip.Position()
		print(" ")
		print(x)
		print(" ")
		println(y)
		if !b.WithinBounds(x, y) {
			return false
		} else if other := b.ChipAtSpace(x, y); other != nil {
			if other.Rune() != chip.Rune() {
				return false
			}
		}
	}
	return true
}

func (b *board) PointsForMove(m Move) int {
	points := 0
	multipliers := make([]int, 0)

	for chip := range m.Iter() {
		x, y := chip.Position()
		if b.grid[x][y].chip != nil {
			switch b.MultiplierForSpace(x, y) {
			case normal:
				points += chip.Points()
			case double_letter_score:
				points += chip.Points() * 2
			case triple_letter_score:
				points += chip.Points() * 3
			case double_word_score:
				points += chip.Points()
				multipliers = append(multipliers, 2)
			case triple_word_score:
				points += chip.Points()
				multipliers = append(multipliers, 3)
			}
		} else {
			points += chip.Points()
		}
	}

	for i := range multipliers {
		points *= i
	}

	return points
}

func (b *board) MakeMove(m Move) (int, bool) {
	if b.IsMoveValid(m) {
		points := b.PointsForMove(m)

		for chip := range m.Iter() {
			x, y := chip.Position()
			b.grid[x][y].chip = chip
		}

		return points, true
	}
	return 0, false
}

func (b *board) SliceRep() [][]interface{} {
	representation := make([][]interface{}, b.height)

	for y := 0; y < b.height; y+=1 {

		representation[y] = make([]interface{}, b.width)
		for x := 0; x < b.width; x+=1 {

			if b.grid[x][y].chip == nil {
				representation[y][x] = b.grid[x][y].chip.MapRep()
			} else {
				representation[y][x] = b.grid[x][y].multiplier.string()
			}
		}
	}

	return representation
}