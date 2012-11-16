package scrabble

// a struict representation of a move
type move struct {
	chips	[]Chip 	// an array of chips to be used for the move
}

// make a new move instance
func NewMove() Move {
	return new(move)
}

func (m *move) Init() Move {
	m.chips = make([]Chip, 2)
	return m
}

func (m *move) InitWithChips(chips []Chip) Move {
	m.chips = chips
	return m
}

func (m *move) Iter() <-chan Chip {
	chips_chan := make(chan Chip, len(m.chips))

	go func() {
		for _, chip := range m.chips {
			chips_chan <- chip
		}
	}()

	return chips_chan
}

func (m *move) Add(chip Chip) {
	m.chips = append(m.chips, chip)
}