package dictionary

// this is a special wildcard character
const WILDCARD := '*'

// represents the interface for a dictionary trie node
type DictNode interface {
	// add a word to the node of a dictionary trie
	Add(word string)

	// find a word starting from the node
	Find(word string) bool

	// find all possible dictionary word matches for the given slice of runes
	// returns a channel of all possible words
	FindAllPossible(runes []rune) <-chan string
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
	var rune rune
	var rest string

	switch x := len(word) {
	case 0:
		n.wordEnd = true
		return
	default:
		rune = word[0]
		rest = word[1:]
	}

	if node, ok := n.nodes[r]; ok {
		node.Add(word[1:])
	} else {
		node := NewNode()
		n.nodes[r] = node
		node.Add(word[1:])
	}
}

func (n *dictNode) Find(word string) bool {
	if word == "" { return n.wordEnd; }

	r := word[0]

	if node, ok := n.nodes[r]; ok {
		return node.Find(word[1:])
	}

	return false
}

func (n *dictNode) FindAllPossible(runes []rune) <-chan string {
	results := make(chan string, 100)
	searchStart := make(chan int, 1)
	searchFinish := make(chan int, 1)

	// make sure the results channel is eventually closed
	go func() {
		count := 0
		for {
			select {
			case s := <-searchStart:
				count += 1
			case s := <-searchFinish:
				count -= 1
			}

			if count == 0 {
				close(results)
				return
			}
		}
	}

	var find func(n *dictNode, runes []rune, word string)
	find = func(n *dictNode, runes []rune, word string) {
		searchStart <- 1

		if n.wordEnd { results <- word; }
		
		visitedRunes := make(map[rune]bool)
		for index, rune := range runes {
			if _, runeIsVisited := visitedRunes[rune]; !runeIsVisited {
				newRunes := append(runes[:index], runes[index+1:]...)

				if rune == SPECIAL {
					for rune, node := range n.nodes {
						go find(node, newRunes, string(append([]rune(word), rune)))
					}
				} else {
					if node, nodeExists := n.nodes[rune]; nodeExists {
						go find(node, newRunes, string(append([]rune(word), rune)))
					}
				}
			}			
		}

		searchFinish <- -1
	}

	go find(n, runes, "")

	return results
}