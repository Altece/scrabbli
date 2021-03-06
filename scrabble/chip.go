package scrabble

import (
	"sync"
	"strings"
)

// a singleton representation of the points map
// this map is guaranteed to only be created once
type points_singleton struct {
	sync.Once
	// initialized 	bool
	points 			map[rune]int
}

// the singleton instance of the points_singleton struct
var points_holder points_singleton

// get the initialized points map
func points() map[rune]int {
	points_holder.Do( func(){
	// if !points_holder.initialized {
		points_holder.points = map[rune]int{
			// 0 points
			' ': 0,
			
			// 1 point
			'E': 1,
			'A': 1,
			'I': 1,
			'O': 1,
			'N': 1,
			'R': 1,
			'T': 1,
			'L': 1,
			'S': 1,
			'U': 1,

			// 2 points
			'D': 2,
			'G': 2,

			// 3 points
			'B': 3,
			'C': 3,
			'M': 3,
			'P': 3,

			// 4 points
			'F': 4,
			'H': 4,
			'V': 4,
			'W': 4,
			'Y': 4,

			// 5 points
			'K': 5,

			// 8 points
			'J': 8,
			'X': 8,

			// 10 points
			'Q': 10,
			'Z': 10	}
	} )
	// }
	return points_holder.points
}

// a struct that represents a chip on a scrabble board
type chip struct {
	r					rune 	// the character represented by the chip
	placedRune			rune	// the character represented by the chip after placed
	x, y 				int 	// the x y coordinate position of the chip
}

// return either the capital version of the given rune 
// (if it is not already capital), or a space
func capitalRune(r rune) rune {
	str := string(r)
	upper := strings.ToUpper(str)
	return rune(upper[0])
}

// make a new chip instance
func NewChip() Chip {
	return new(chip)
}

func (c *chip) Init(r rune, x, y int) Chip {
	c.r = r
	c.placedRune = r;
	c.x, c.y = x, y
	return c
}

func (c *chip) InitPlaced(r, p rune, x, y int) Chip {
	c.Init(r, x, y)
	c.placedRune = p;
	return c
}

func (c *chip) ChipRune() rune {
	return c.r
}

func (c *chip) PlacedRune() rune {
	return c.placedRune
}

func (c *chip) Position() (int, int) {
	return c.x, c.y
}

func (c *chip) Points() int {
	points := points()
	num_points, ok := points[capitalRune(c.r)]
	if ok {
		return num_points
	}
	return 0;
}

func (c *chip) MapRep() map[string]interface{} {
	representation := map[string]interface{} {
		"Rune" 		: c.r,
		"Placed"	: c.placedRune,
		"Points" 	: c.Points(),
		"X" 		: c.x,
		"Y" 		: c.y 	}
	return representation
}