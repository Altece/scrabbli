package dictionary

import "sync"

// this is a special wildcard character
const WILDCARD rune = ' '

// represents the interface for a dictionary trie node
type DictNode interface {
	// add a word to the node of a dictionary trie
	Add(word string)

	// find a word starting from the node
	Find(word string) bool

	// find all possible dictionary word matches for the given slice of runes
	// returns a channel of all possible rune combinations
	FindAllPossible(runes []rune) <-chan []rune
}

// a node data structure for a dictionary trie
type dictNode struct {
	nodes	map[rune]*dictNode 	// a map of rune to node mapings
	wordEnd bool 				// a boolean flag to indicate that this 
								// node can be the end of a word
}

// make a new dictionary trie nde
func NewNode() *dictNode {
	n := new(dictNode)
	n.nodes = make(map[rune]*dictNode)
	n.wordEnd = false
	return n
}

func (n *dictNode) Add(word string) {
	var r rune
	var rest string

	switch len(word) {
	case 0:
		n.wordEnd = true
		return
	default:
		r = rune(word[0])
		rest = word[1:]
	}

	if node, ok := n.nodes[r]; ok {
		node.Add(rest)
	} else {
		node := NewNode()
		n.nodes[r] = node
		node.Add(rest)
	}
}

func (n *dictNode) Find(word string) bool {
	if word == "" { return n.wordEnd; }

	r := rune(word[0])

	if node, ok := n.nodes[r]; ok {
		return node.Find(word[1:])
	}

	return false
}

func (n *dictNode) FindAllPossible(runes []rune) <-chan []rune {
	results := make(chan []rune, 100)

	var waitGroup sync.WaitGroup

	// make sure the results channel is eventually closed
	go func() {
		waitGroup.Wait()
		close(results)
	}()

	resulting := func(runes, word []rune, i int) (resultingRunes, resultingWord []rune) {
		resultingRunes = append(append(resultingRunes, runes[:i]...), runes[i+1:]...)

		resultingWord = append(resultingWord, word...)

		return resultingRunes, resultingWord
	}

	var find func(n *dictNode, runes, word []rune)
	find = func(n *dictNode, runes, word []rune) {

		if n.wordEnd { results <- word; }

		for i, r := range runes {
			

			if r == WILDCARD {
				for _, node := range n.nodes {
					resultingRunes, resultingWord := resulting(runes, word, i)

					waitGroup.Add(1)
					resultingWord = append(resultingWord, r)
					go find(node, resultingRunes, resultingWord)
				}
			} else {
				resultingRunes, resultingWord := resulting(runes, word, i)

				if node, nodeExists := n.nodes[r]; nodeExists {
					waitGroup.Add(1)
					resultingWord = append(resultingWord, r)
					go find(node, resultingRunes, resultingWord)
				}
			}
		}

		waitGroup.Done()
	}

	waitGroup.Add(1)
	var word []rune
	go find(n, runes, word)

	return results
}