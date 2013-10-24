package dictionary

import "sync"

// this is a special wildcard character
const WILDCARD rune = ' '

type RunePair struct {
	Word 		rune
	Combination rune
}

type RuneResult struct {
	Word 		[]rune
	Combination []rune
}

func (result RuneResult) Unzip() []RunePair {
	var runepairs []RunePair
	for i, _ := range result.Word {
		runepairs = append(runepairs, RunePair{result.Word[i], result.Combination[i]})
	}
	return runepairs
}

type StringResult struct {
	Word 		string
	Combination string
}

// represents the interface for a dictionary trie node
type DictNode interface {
	// add a word to the node of a dictionary trie
	Add(word string)

	// find a word starting from the node
	// returns an array of all resulting strings and a success bool
	Find(word string) ([]string, bool)

	// find all possible dictionary word matches for the given slice of runes
	// returns a channel for all possible rune combinations
	FindAllRunes(runes []rune) <-chan RuneResult

	FindAllStrings(input string) <-chan StringResult

	FindAllRunesFromRunePairs(runepairs []RunePair) <-chan RuneResult
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

func (n *dictNode) Find(word string) ([]string, bool) {
	if word == "" { return []string{}, n.wordEnd; }

	add := func(r rune, results []string, adding []string) []string {
		for _, word := range adding {
			str := string(append([]rune{r}, []rune(word)...))
			results = append(results, str)
		}
		if len(adding) == 0 {
			results = append(results, string([]rune{r}))
		}
		return results
	}

	r := rune(word[0])
	if r == WILDCARD {
		results := []string{}
		success := false
		for l, node := range n.nodes {
			if result, ok := node.Find(word[1:]); ok {
				results = add(l, results, result)
				success = true
			}
		}
		if success {
			return results, true
		}
		return results, false
	}

	if node, ok := n.nodes[r]; ok {
		if result, success := node.Find(word[1:]); success {
			return add(r, []string{}, result), true
		}
	}

	return []string{}, false
}

func (dict *dictNode) FindAllRunesFromRunePairs(runes []RunePair) <-chan RuneResult {
	results := make(chan RuneResult, 100)

	var waitGroup sync.WaitGroup

	// make sure the results channel is eventually closed
	go func() {
		waitGroup.Wait()
		close(results)
	}()

	placeRuneIntoResult := func(runes []RunePair, i int, r rune, result RuneResult) (leftoverRunes []RunePair, resulting RuneResult) {
		leftoverRunes = append(append(leftoverRunes, runes[:i]...), runes[i+1:]...)

		resulting.Word = append(append(resulting.Word, result.Word...), r)
		resulting.Combination = append(append(resulting.Combination, result.Combination...), runes[i].Combination)

		return leftoverRunes, resulting
	}

	var find func(n *dictNode, runes []RunePair, result RuneResult)
	find = func(n *dictNode, runes []RunePair, result RuneResult) {

		if n.wordEnd {
			results <- result;
		}

		for i, r := range runes {

			if r.Word == WILDCARD {
				for r, node := range n.nodes {
					waitGroup.Add(1)
					remaining, result := placeRuneIntoResult(runes, i, r, result)
					go find(node, remaining, result)
				}
			} else {
				remaining, result := placeRuneIntoResult(runes, i, r.Word, result)
				if node, nodeExists := n.nodes[r.Word]; nodeExists {
					waitGroup.Add(1)
					go find(node, remaining, result)
				}
			}
		}

		waitGroup.Done()
	}

	waitGroup.Add(1)
	var result RuneResult
	go find(dict, runes, result)

	return results
}

func (n *dictNode) FindAllRunes(runes []rune) <-chan RuneResult {
	var runepairs []RunePair
	for _, r := range runes {
		runepairs = append(runepairs, RunePair{r,r})
	}
	return n.FindAllRunesFromRunePairs(runepairs);
}

func (n *dictNode) FindAllStrings(input string) <-chan StringResult {
	results := make(chan StringResult, 100)
	go func() {
		for result := range n.FindAllRunes([]rune(input)) {
			results <- StringResult{string(result.Word), string(result.Combination)}
		}
		close(results)
	}()
	return results;
}
