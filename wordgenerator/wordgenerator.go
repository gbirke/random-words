package wordgenerator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

// A Chunk is a 3-letter string that can return parts of itself
type Chunk string

// Return last 2 chars
func (c Chunk) lastChars() string {
	return string(c)[1:]
}

// Return first 2 chars
func (c Chunk) firstChars() string {
	return string(c)[0:2]
}

// return last char (as a "starting place" for the next word)
func (c Chunk) nextChar() string {
	return string(c)[2:]
}

type ChunkList []Chunk

// split word into 3-character chunks
func NewChunkListFromWord(word string) ChunkList {
	const ChunkLength = 3
	length := len(word)
	if length < ChunkLength {
		return ChunkList{}
	}

	var result = make(ChunkList, length-2)
	for i := 0; i < length-ChunkLength+1; i++ {
		result[i] = Chunk(word[i : i+ChunkLength])
	}
	return result
}

// A collection of unique chunks, with counts of how often they have appeared in a source text
type ChunkCounts map[Chunk]int

func (cc ChunkCounts) GetChunkList() ChunkList {
	chunks := make(ChunkList, len(cc))
	i := 0
	for k := range cc {
		chunks[i] = k
		i++
	}
	return chunks
}

func NewChunkCountsFromReader(r io.Reader) (ChunkCounts, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	nonWordChars := regexp.MustCompile(`[\W_]+`)
	chunkCounts := make(ChunkCounts)
	for scanner.Scan() {
		word := nonWordChars.ReplaceAllString(strings.ToLower(scanner.Text()), "")
		// TODO investigate removing stop words like "the", "not", etc to avoid skewing the corpus
		for _, chunk := range NewChunkListFromWord(word) {
			if chunk != "" {
				chunkCounts[chunk]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// TODO we should postprocess the result to avoid having "dead end" chunks that have ending chars with no matching start chars
	//      then we can get rid of the error in getNextChar

	return chunkCounts, nil
}

// A collection of chunks, sorted by their first 2 starting characters
type AdjacentChunkMap map[string]ChunkList

func NewAdjacentChunkMap(chunks ChunkList) AdjacentChunkMap {
	result := make(AdjacentChunkMap)
	for _, chunk := range chunks {
		startChars := chunk.firstChars()
		result[startChars] = append(result[startChars], chunk)
	}
	return result
}

type WordGenerator struct {
	chunkCounts ChunkCounts
}

func NewWordGeneratorFromString(text string) (*WordGenerator, error) {

	chunkCounts, err := NewChunkCountsFromReader(strings.NewReader(text))

	if err != nil {
		return nil, err
	}

	return &WordGenerator{chunkCounts: chunkCounts}, nil
}

func NewWordGeneratorFromFile(filename string) (*WordGenerator, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	chunkCounts, err := NewChunkCountsFromReader(f)

	if err != nil {
		return nil, err
	}

	return &WordGenerator{chunkCounts: chunkCounts}, nil
}

func (w WordGenerator) getNextChar(currentChunk Chunk, adjacentChunks AdjacentChunkMap) (string, error) {
	lastChars := currentChunk.lastChars()
	candidateChunks, ok := adjacentChunks[lastChars]
	if !ok {
		return "", errors.New(fmt.Sprintf("No adjacent chunk found for %v", currentChunk))
	}

	// No random selection for small adjacencies
	if len(candidateChunks) == 1 {
		return candidateChunks[0].nextChar(), nil
	}

	// Weighted selection, see https://stackoverflow.com/a/11872928/130121
	weightedChunks := make(ChunkCounts, len(candidateChunks))
	sumWeights := 0
	for _, chunk := range candidateChunks {
		weightedChunks[chunk] = w.chunkCounts[chunk]
		sumWeights += weightedChunks[chunk]
	}
	targetWeight := rand.Intn(sumWeights)
	for chunk, weight := range weightedChunks {
		targetWeight -= weight
		if targetWeight <= 0 {
			return chunk.nextChar(), nil
		}
	}

	// This should never happen
	// TODO use different error classes to distinguish algorithm error and corpus error
	return "", errors.New("did not get weighted result, this should never happen")
}

func (w WordGenerator) GetWord(length int) (string, error) {
	chunks := w.chunkCounts.GetChunkList()
	chunkMap := NewAdjacentChunkMap(chunks)

	currentChunk := chunks[rand.Intn(len(chunks))]
	word := string(currentChunk)
	failedChunks := make([]Chunk, length)
	for i := 3; i <= length; i++ {
		nextChar, err := w.getNextChar(currentChunk, chunkMap)
		if err != nil {
			failedChunks = append(failedChunks, currentChunk)
			if len(failedChunks) > length*2 {
				// TODO maybe return failedChunks for diagnosis
				return "", errors.New(fmt.Sprintf("Could not find a matching chunks %v times, maybe your input text is too small", len(failedChunks)))
			}
			// reset the loop, start fresh
			i = 3
			currentChunk = chunks[rand.Intn(len(chunks))]
			word = string(currentChunk)
			continue
		}
		word += nextChar
		// TODO create static method for this
		currentChunk = Chunk(word[len(word)-3:])
	}
	return word, nil
}

func (w WordGenerator) PrintDebugInfo() {
	fmt.Printf("%v", w.chunkCounts)
	fmt.Println("")
	fmt.Printf("%v", NewAdjacentChunkMap(w.chunkCounts.GetChunkList()))
	fmt.Println("")
}
