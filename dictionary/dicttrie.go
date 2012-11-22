package dictionary

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
	signalchan := make(chan int, 100)

	// make sure the results channel is eventually closed
	go func() {
		count := <-signalchan
		for {
			count += <-signalchan

			if count == 0 {
				close(results)
				return
			}
		}
	}()

	var find func(n *dictNode, runes []rune, word []rune)
	find = func(n *dictNode, runes []rune, word []rune) {
		signalchan <- 1

		if n.wordEnd { results <- word; }
		
		visitedRunes := make(map[rune]bool)
		for index, rune := range runes {
			if _, runeIsVisited := visitedRunes[rune]; !runeIsVisited {
				newRunes := append(runes[:index], runes[index+1:]...)

				if rune == WILDCARD {
					for _, node := range n.nodes {
						find(node, newRunes, append(word, rune))
					}
				} else {
					if node, nodeExists := n.nodes[rune]; nodeExists {
						find(node, newRunes, append(word, rune))
					}
				}
			}			
		}

		signalchan <- -1
	}

	go find(n, runes, make([]rune, 0))

	return results
}