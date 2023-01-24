package main

import (
	"testing"
)

func TestGetCharacterChunks(t *testing.T) {
	chunks := getCharacterChunks("explain")
	expected := []string{"exp", "xpl", "pla", "lai", "ain"}

	if len(expected) != len(chunks) {
		t.Fatalf(`Want: %v Got: %v (%v)`, len(expected), len(chunks), chunks)
	}
	for i := range expected {
		if chunks[i] != expected[i] {
			t.Fatalf(`Want: %v Got: %v at index %v`, expected[i], chunks[i], i)
		}
	}
}

func TestGetCharacterChunksInThreeLetterWord(t *testing.T) {
	chunks := getCharacterChunks("hat")
	expected := []string{"hat"}
	for i := range expected {
		if chunks[i] != expected[i] {
			t.Fatalf(`Want: %v Got: %v at index %v`, expected[i], chunks[i], i)
		}
	}
}

func TestGetCharacterChunksInSmallWords(t *testing.T) {
	empty := getCharacterChunks("")
	oneLetter := getCharacterChunks("i")
	twoLetters := getCharacterChunks("is")
	results := [][]string{empty, oneLetter, twoLetters}
	for i := range results {
		if len(results[i]) > 0 {
			t.Fatalf(`Wanted empty result, got length %v: %v`, len(results[i]), results[i])

		}
	}
}

func TestSplitTextToLowercaseWords(t *testing.T) {
	type Case struct {
		name     string
		input    string
		expected []string
	}
	cases := []Case{
		{
			"simple case",
			"Explain this to me",
			[]string{"explain", "this", "to", "me"},
		},
		{
			"punctuation at the end",
			"Explain THIS to me!",
			[]string{"explain", "this", "to", "me"},
		},

		{
			"no empty elements on spaces following punctuation",
			"EXTERMINATE! EXTERMINATE!",
			[]string{"exterminate", "exterminate"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			words := splitTextToLowercaseWords(c.input)
			if len(c.expected) != len(words) {
				t.Fatalf(`Want: %v Got: %v (%v)`, len(c.expected), len(words), words)
			}
			for i := range c.expected {
				if words[i] != c.expected[i] {
					t.Fatalf(`Want: %v Got: %v at index %v`, c.expected[i], words[i], i)
				}
			}
		})

	}
}

func TestGetChunksFromText(t *testing.T) {
	type Case struct {
		name     string
		input    string
		expected map[string]int
	}
	cases := []Case{
		{
			"simple case",
			"Explain this to me",
			map[string]int{"exp": 1, "xpl": 1, "pla": 1, "lai": 1, "ain": 1, "thi": 1, "his": 1},
		},
		{
			"overlapping chunks",
			"This is his hat, that is her hair, whistfully",
			map[string]int{"thi": 1, "his": 3, "hat": 2, "tha": 1, "her": 1, "hai": 1, "air": 1, "whi": 1, "ist": 1, "stf": 1, "tfu": 1, "ful": 1, "ull": 1, "lly": 1},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			chunks := getChunksFromText(c.input)
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

func TestGetAdjacentChunkMap(t *testing.T) {
	chunkMap := getAdjacentChunkMap([]string{"the", "thi", "her", "hep", "hel", "ill"})
	expected := map[string][]string{"th": {"the", "thi"}, "he": {"her", "hep", "hel"}, "il": {"ill"}}

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
