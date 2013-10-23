package main

import (
	"fmt"
	"./dictionary"
)

func hasLetters(str []rune, ls []rune) bool {
	for _, s := range str {
		cont := false;
		for _, l := range ls {
			if s == l {
				cont = true
				break
			}
		}
		if cont == false {
			return false
		}
	}
	return true
}

func main() {
	dict := dictionary.NewDictionary()
	dict.Init("dictionary.txt")
	// results, _ := dict.FindAllStrings("FOREST")

	for result := range dict.FindAllStrings("FOR EST") {
		if _, ok := dict.Find(result.Word); !ok {
			fmt.Println("failure with", result.Word, result.Combination)
		} else {
			fmt.Printf("%s:(%s)\n", result.Word, result.Combination)
		}
	}
	/*
	for s := range dict.FindAllPossible([]rune("  SOEK")) {
		if len(s) == 5 && 
			s[0] == rune('S') && s[4] == rune('E') && s[2] == rune('O') {
			fmt.Println(string(s));
		}
	}
	*/
}