package wordgenerator_test

import (
	"github.com/gbirke/random-words/wordgenerator"
	"strings"
	"testing"
)

func TestNewChunkListFromWord(t *testing.T) {

	type Case struct {
		name     string
		input    string
		expected wordgenerator.ChunkList
	}
	cases := []Case{
		{
			"one word",
			"explain",
			wordgenerator.ChunkList{"exp", "xpl", "pla", "lai", "ain"},
		},
		{
			"three letter word",
			"hat",
			wordgenerator.ChunkList{"hat"},
		},

		{
			"two letter word",
			"is",
			wordgenerator.ChunkList{},
		},
		{
			"one letter word",
			"i",
			wordgenerator.ChunkList{},
		},
		{
			"empty word",
			"",
			wordgenerator.ChunkList{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			chunks := wordgenerator.NewChunkListFromWord(c.input)

			if len(c.expected) != len(chunks) {
				t.Fatalf(`Want: %v Got: %v (%v)`, len(c.expected), len(chunks), chunks)
			}
			for i := range c.expected {
				if chunks[i] != c.expected[i] {
					t.Fatalf(`Want: %v Got: %v at index %v`, c.expected[i], chunks[i], i)
				}
			}
		})
	}
}

func TestNewChunkCountFromReader(t *testing.T) {
	type Case struct {
		name     string
		input    string
		expected wordgenerator.ChunkCounts
	}
	cases := []Case{
		{
			"simple case",
			"Explain this to me",
			wordgenerator.ChunkCounts{"exp": 1, "xpl": 1, "pla": 1, "lai": 1, "ain": 1, "thi": 1, "his": 1},
		},
		{
			"overlapping chunks",
			"This is his hat, that is her hair, whistfully",
			wordgenerator.ChunkCounts{"thi": 1, "his": 3, "hat": 2, "tha": 1, "her": 1, "hai": 1, "air": 1, "whi": 1, "ist": 1, "stf": 1, "tfu": 1, "ful": 1, "ull": 1, "lly": 1},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			chunks, err := wordgenerator.NewChunkCountsFromReader(strings.NewReader(c.input))
			if err != nil {
				t.Fatalf("Got error %v", err)
			}
			if len(c.expected) != len(chunks) {
				t.Fatalf(`Want: %v (%v) Got: %v (%v)`, len(c.expected), c.expected, len(chunks), chunks)
			}
			for chunk, chunkCount := range c.expected {
				actualCount, ok := chunks[chunk]
				if !ok {
					t.Fatalf(`Chunk "%v" not found in %v`, chunk, chunks)
				}
				if actualCount != chunkCount {
					t.Fatalf(`Want: %v Got: %v for chunk %v`, chunkCount, actualCount, chunk)
				}
			}
		})
	}
}

func TestNewAdjacentChunkMap(t *testing.T) {
	chunkMap := wordgenerator.NewAdjacentChunkMap(wordgenerator.ChunkList{"the", "thi", "her", "hep", "hel", "ill"})
	expected := wordgenerator.AdjacentChunkMap{"th": {"the", "thi"}, "he": {"her", "hep", "hel"}, "il": {"ill"}}

	if len(expected) != len(chunkMap) {
		t.Fatalf(`Want: %v Got: %v (%v)`, len(expected), len(chunkMap), chunkMap)
	}
	for startChars, chunks := range expected {
		actualChunks, ok := chunkMap[startChars]
		if !ok {
			t.Fatalf(`Starting chars "%v" not found in %v`, startChars, chunkMap)
		}
		if len(actualChunks) != len(chunks) {
			t.Fatalf(`Want: %v (%v) Got: %v (%v)`, len(actualChunks), actualChunks, len(chunks), chunks)
		}
	}
}
