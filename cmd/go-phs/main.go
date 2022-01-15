package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Hand struct {
	cards []Card
	rank  int
}

type Card struct {
	value int
	suit  string
}

func main() {
	// STDIN Reader
	scanner := bufio.NewScanner(os.Stdin)

	// Loop over lines from STDIN
	for scanner.Scan() {
		line := scanner.Text()

		handOne, handTwo, err := parseHands(line)
		// There shouldn't be an error when input is as expected, so if there is, log it and quit as we can't continue
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Hand one: {%v}\n", handOne)
		fmt.Printf("Hand two: {%v}\n", handTwo)
	}
}

// Parse a line of input text to return the hands of the two players
func parseHands(line string) (Hand, Hand, error) {
	// Split the line on spaces into two character card strings
	cardStrings := strings.Split(line, " ")

	// Convert the character pairs to card objects
	handOne := Hand{cards: []Card{}}
	handTwo := Hand{cards: []Card{}}

	for i, pair := range cardStrings {
		// Split the 2 character long card string into two characters
		characters := strings.Split(pair, "")

		// Convert the first character (value) to an int for easy comparison
		value, err := charToValue(characters[0])
		// Errors shouldn't occur on the input file, but handle them just in case
		if err != nil {
			return Hand{}, Hand{}, err
		}

		// Use the second character (suit) as is
		suit := characters[1]

		// The first five cards go into player one's hand, the next five into player two's
		if i < 5 {
			handOne.cards = append(handOne.cards, Card{value, suit})
		} else {
			handTwo.cards = append(handTwo.cards, Card{value, suit})
		}
	}

	return handOne, handTwo, nil
}

// Convert the first character from the two character card strings to an integer representation for
// easier comparison
func charToValue(char string) (int, error) {
	// First try converting character directly to int which handles all digits (2-9)
	value, err := strconv.Atoi(char)

	// If converting to digit fails, try matching character with Ten, Jack, Queen, King, Ace
	if err != nil {
		switch char {
		case "T":
			value = 10
		case "J":
			value = 11
		case "Q":
			value = 12
		case "K":
			value = 13
		case "A":
			value = 14

		// If digit conversion and character matching both fail, something was wrong with the input
		default:
			return -1, fmt.Errorf("Unable to convert {%s} to a value", char)
		}
	}

	return value, nil
}
