package dictionary

import (
	"os"
	"bufio"
	"strings"
)

// an interface for a Dictionary of words
type Dictionary interface {
	// initialize the dictionary with the given file
	// returns the dictionary and a bool flag for file success
	Init(filename string) (Dictionary, bool)

	// a dictionary is also a dictionary node
	DictNode
}

// get a new dictionary
func NewDictionary() *dictNode {
	return NewNode()
}

func (d *dictNode) Init(filename string) (Dictionary, bool) {
	if file, error := os.Open(filename); error == nil {
		defer func() { file.Close() }()
		reader := bufio.NewReader(file)

		for {
			line, error := reader.ReadString('\n')
			if len(line) != 0 { line = line[:len(line)-1]; }

			if line != "" {
				d.Add(strings.ToUpper(line))
			}

			if error != nil { return d, true; }
		}
	}
	return d, false
}