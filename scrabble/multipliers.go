package scrabble

// define the multiplier enum type identifiers
type Multiplier int
const (
	// represents a normal space
	normal	Multiplier = iota

	// represents a letter score bonus
	double_letter_score
	triple_letter_score

	// represents a word score bonus
	double_word_score
	triple_word_score
)

func (m Multiplier) string() string {
	switch m {
	case normal:
		return "NORMAL"
	case double_letter_score:
		return "DOUBLE LETTER SCORE"
	case triple_letter_score:
		return "TRIPLE LETTER SCORE"
	case double_word_score:
		return "DOUBLE WORD SCORE"
	case triple_word_score:
		return "TRIPLE WORD SCORE"
	}
	return ""
}

// this is for a 15 x 15 board only (for now)
// idk how to make scrabble boards of indefinite size...
func (b *board) multiplier(x, y int) Multiplier {
	max_x, max_y := b.width-1, b.height-1

	outside := func(a, b int) bool {
		xInside := x >= a && x <= b
		yInside := y >= a && y <= b
		return !(xInside && yInside)
	}

	// double word score
	// the center spot
	if (x == y) && (x == max_x/2) {
		return double_word_score
	}


	// triple word score
	// outer grid of halfs
	if (x % (max_x/2) == 0) && (y % (max_y/2) == 0) {
		return triple_word_score
	}


	// double word score
	// the diagnols that aren't already something and are not in the center
	if (x == y || max_x-x == y || max_y-y == x) && outside(max_x-2, max_x+2) {
		return double_word_score
	}


	// triple letter score
	// inner grid of fourths that aren't already something
	if (x-1 == (max_x-2)/4) && (y-1 == (max_y-2)/4) {
		return triple_letter_score
	}


	// double letter score
	// the following is very 15x15 board specific...
	boxAroundCenter := func(a, max int) bool {
		return (a == (max/2)-1 || a == (max/2)+1)
	}
	// a box around the center
	if boxAroundCenter(x, max_x) && boxAroundCenter(y, max_y) {
		return double_letter_score
	}
	// on the outside border
	if outside(1, 13) {
		if (x == 3 || x == 11) || (y == 3 || y == 11) {
			return double_letter_score
		}
	}
	// those triangle things
	if !outside(3, 12) {
		if (x == 3 || x == 11) && (y == max_y/2) {
			return double_letter_score
		}
		if (y == 3 || y == 11) && (x == max_x/2) {
			return double_letter_score
		}
		if (x == 2 || x == 12) && (y == max_y-1 || y == max_y+1) {
			return double_letter_score
		}
		if (y == 2 || y == 12) && (x == max_x-1 || x == max_x+1) {
			return double_letter_score
		}
	}

	return normal
}