package main

import (
	"fmt"
	"github.com/gbirke/random-words/wordgenerator"
	"os"
	"strconv"
)

const DefaultInputSmall = `But there is something that I must say to my people, who stand on the warm threshold which leads into the palace of justice: In the process of gaining our rightful place, we must not be guilty of wrongful deeds. Let us not seek to satisfy our thirst for freedom by drinking from the cup of bitterness and hatred. We must forever conduct our struggle on the high plane of dignity and discipline. We must not allow our creative protest to degenerate into physical violence. Again and again, we must rise to the majestic heights of meeting physical force with soul force. ... I am not unmindful that some of you have come here out of great trials and tribulations. Some of you have come fresh from narrow jail cells. And some of you have come from areas where your quest -- quest for freedom left you battered by the storms of persecution and staggered by the winds of police brutality. You have been the veterans of creative suffering. Continue to work with the faith that unearned suffering is redemptive. Go back to Mississippi, go back to Alabama, go back to South Carolina, go back to Georgia, go back to Louisiana, go back to the slums and ghettos of our northern cities, knowing that somehow this situation can and will be changed. Let us not wallow in the valley of despair, I say to you today, my friends. And so even though we face the difficulties of today and tomorrow, I still have a dream. It is a dream deeply rooted in the American dream. I have a dream that one day this nation will rise up and live out the true meaning of its creed: "We hold these truths to be self-evident, that all men are created equal." I have a dream that one day on the red hills of Georgia, the sons of former slaves and the sons of former slave owners will be able to sit down together at the table of brotherhood. I have a dream that one day even the state of Mississippi, a state sweltering with the heat of injustice, sweltering with the heat of oppression, will be transformed into an oasis of freedom and justice. I have a dream that my four little children will one day live in a nation where they will not be judged by the color of their skin but by the content of their character.`

func main() {
	args := os.Args[1:]
	corpusFile := "corpus.txt"
	wordLength := 6
	numWords := 150
	var err error
	// TODO use proper options instead of positional arguments
	if len(args) > 0 {
		corpusFile = args[0]
	}
	if len(args) > 1 {
		wordLength, err = strconv.Atoi(args[1])
		if err != nil || wordLength > 20 {
			fmt.Println("Invalid word length")
			wordLength = 6
		}
	}
	if len(args) > 2 {
		numWords, err = strconv.Atoi(args[2])
		if err != nil || numWords > 1000 {
			fmt.Println("Maximum for number of words is 1000")
			numWords = 1000
		}
	}
	generator, err := wordgenerator.NewWordGeneratorFromFile(corpusFile)
	if err != nil {
		fmt.Printf("Could not read file %v\nUsing default corpus.\n", corpusFile)
		generator, err = wordgenerator.NewWordGeneratorFromString(DefaultInputSmall)
		if err != nil {
			panic("Could not generate corpus from defaults")
		}
	}
	// wordgenerator.PrintDebugInfo(chunkCounts)
	for i := 0; i < numWords; i++ {
		word, err := generator.GetWord(wordLength)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Printf("%v  ", word)
			// TODO add CLI flag for output format
			if (i+1)%5 == 0 {
				fmt.Println("")
			}
		}
	}
	fmt.Println("")
}
