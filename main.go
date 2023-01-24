package main

import "fmt"
import "strings"
import "regexp"

// split word into 3-character chunks
func getCharacterChunks(word string) []string {
	const ChunkLength = 3
	length := len(word)
	if length < ChunkLength {
		return []string{}
	}

	var result = make([]string, length-2)
	for i := 0; i < length-ChunkLength+1; i++ {
		result[i] = word[i : i+ChunkLength]
	}
	return result
}

func splitTextToLowercaseWords(text string) []string {
	wordBoundary := regexp.MustCompile(`\W+`)
	words := wordBoundary.Split(strings.ToLower(text), -1)
	if words[len(words)-1] == "" {
		return words[:len(words)-1]
	}
	return words
}

func getChunksFromText(text string) map[string]int {
	result := make(map[string]int)
	for _, word := range splitTextToLowercaseWords(text) {
		for _, chunk := range getCharacterChunks(word) {
			if chunk != "" {
				result[chunk]++
			}
		}
	}

	return result
}

func main() {
	fmt.Println("Hello, world.")
}
