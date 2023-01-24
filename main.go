package main

import "fmt"
import "math/rand"
import "strings"
import "regexp"
import "errors"

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
	// TODO investigate removing stop words like "the", "not", etc to avoid skewing the corpus
	if words[len(words)-1] == "" {
		return words[:len(words)-1]
	}
	return words
}

type ChunkCounts map[string]int

func getChunksFromText(text string) ChunkCounts {
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

type AdjacentChunkMap map[string][]string

func getAdjacentChunkMap(chunks []string) AdjacentChunkMap {
	result := make(map[string][]string)
	for _, chunk := range chunks {
		startChars := chunk[0:2]
		result[startChars] = append(result[startChars], chunk)
	}
	return result
}

func getNextChar(currentChunk string, allChunks ChunkCounts, adjacentChunks AdjacentChunkMap) (string, error) {
	lastChars := currentChunk[1:]
	candidateChunks, ok := adjacentChunks[lastChars]
	if !ok {
		return "", errors.New(fmt.Sprintf("No adjacent chunk found for %v", currentChunk))
	}

	// TODO use weights from allChunks to weigh result
	idx := rand.Intn(len(candidateChunks))

	return candidateChunks[idx][2:], nil
}

func getChunksFromCounts(chunkCounts ChunkCounts) []string {
	chunks := make([]string, len(chunkCounts))
	i := 0
	for k := range chunkCounts {
		chunks[i] = k
		i++
	}
	return chunks
}

func getWord(chunkCounts ChunkCounts, length int) (string, error) {
	chunks := getChunksFromCounts(chunkCounts)
	chunkMap := getAdjacentChunkMap(chunks)

	currentChunk := chunks[rand.Intn(len(chunks))]
	word := currentChunk
	failedChunks := make([]string, length)
	for i := 3; i <= length; i++ {
		nextChar, err := getNextChar(currentChunk, chunkCounts, chunkMap)
		if err != nil {
			failedChunks = append(failedChunks, currentChunk)
			if len(failedChunks) > length*2 {
				// TODO maybe return failedChunks for diagnosis
				return "", errors.New(fmt.Sprintf("Could not find a matching chunks %v times, maybe your input text is too small", len(failedChunks)))
			}
			// reset the loop, start fresh
			i = 3
			currentChunk = chunks[rand.Intn(len(chunks))]
			word = currentChunk
			continue
		}
		word += nextChar
		currentChunk = word[len(word)-3:]
	}
	return word, nil
}

const DefaultInputSmall = `But there is something that I must say to my people, who stand on the warm threshold which leads into the palace of justice: In the process of gaining our rightful place, we must not be guilty of wrongful deeds. Let us not seek to satisfy our thirst for freedom by drinking from the cup of bitterness and hatred. We must forever conduct our struggle on the high plane of dignity and discipline. We must not allow our creative protest to degenerate into physical violence. Again and again, we must rise to the majestic heights of meeting physical force with soul force. ... I am not unmindful that some of you have come here out of great trials and tribulations. Some of you have come fresh from narrow jail cells. And some of you have come from areas where your quest -- quest for freedom left you battered by the storms of persecution and staggered by the winds of police brutality. You have been the veterans of creative suffering. Continue to work with the faith that unearned suffering is redemptive. Go back to Mississippi, go back to Alabama, go back to South Carolina, go back to Georgia, go back to Louisiana, go back to the slums and ghettos of our northern cities, knowing that somehow this situation can and will be changed. Let us not wallow in the valley of despair, I say to you today, my friends. And so even though we face the difficulties of today and tomorrow, I still have a dream. It is a dream deeply rooted in the American dream. I have a dream that one day this nation will rise up and live out the true meaning of its creed: "We hold these truths to be self-evident, that all men are created equal." I have a dream that one day on the red hills of Georgia, the sons of former slaves and the sons of former slave owners will be able to sit down together at the table of brotherhood. I have a dream that one day even the state of Mississippi, a state sweltering with the heat of injustice, sweltering with the heat of oppression, will be transformed into an oasis of freedom and justice. I have a dream that my four little children will one day live in a nation where they will not be judged by the color of their skin but by the content of their character.`

func printDebugInfo(chunkCounts ChunkCounts) {
	fmt.Printf("%v", chunkCounts)
	fmt.Println("")
	fmt.Printf("%v", getAdjacentChunkMap(getChunksFromCounts(chunkCounts)))
	fmt.Println("")
}

func main() {
	chunkCounts := getChunksFromText(DefaultInputSmall)
	// printDebugInfo(chunkCounts)
	for i := 0; i < 50; i++ {
		word, err := getWord(chunkCounts, 7)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Printf("Word: %v\n", word)
		}
	}
	fmt.Println("")
}
