package main

import (
	"strings"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	wordsMap := make(map[string]int)
	
	for _,word := range words {
		elem, found := wordsMap[word]
		if (found) {
			wordsMap[word] += elem
		} else {
			wordsMap[word] = 1
		}
	}

	return wordsMap
}

//func main() {
//	wc.Test(WordCount)
//}


/*

pointer receivers anvenders til at gøre noget ved structen i stedet for at bruge den til en beregning

*/