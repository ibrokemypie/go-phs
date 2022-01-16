package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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

	handOne.rank = rankHand(handOne)
	handTwo.rank = rankHand(handTwo)

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

func rankHand(hand Hand) int {
	switch {
	case isRoyalFlush(hand.cards):
		return 10
	case isStraightFlush(hand.cards):
		return 9
	case isFourOfAKind(hand.cards):
		return 8
	case isFullHouse(hand.cards):
		return 7
	case isFlush(hand.cards):
		return 6
	case isStraight(hand.cards):
		return 5
	case isThreeOfAKind(hand.cards):
		return 4
	case isTwoPair(hand.cards):
		return 3
	case isPair(hand.cards):
		return 2

	// "High Card", the cards in the hand do not make up any poker hand
	default:
		return 1
	}
}

func isRoyalFlush(cards []Card) bool {
	suit := cards[0].suit

	for _, card := range cards {
		// If any of the cards have a different suit the hand isn't a flush
		if card.suit != suit {
			return false
		}

		// If any of the cards aren't a Ten, Jack, Queen, King or Ace this hand isn't a royal flush
		switch card.value {
		case 10, 11, 12, 13, 14:
			continue

		default:
			return false
		}
	}

	// If the cards in the hand are all of the same suit and are the cards from Ten to Ace, then the hand is
	// a royal flush
	return true
}

func isStraightFlush(cards []Card) bool {
	suit := cards[0].suit

	// First sort the cards by value to make checking for consecutiveness a single pass over the cards
	sort.Slice(cards, func(i, j int) bool { return cards[i].value < cards[j].value })

	for i, card := range cards {
		// If any of the cards have a different suit the hand isn't a flush
		if card.suit != suit {
			return false
		}

		// The cards are sorted, so if any card but the last isn't followed by one of consecutive value the hand
		// can't be a straight
		if i != len(cards)-1 {
			if cards[i+1].value != card.value+1 {
				return false
			}
		}
	}

	return true
}

func isFourOfAKind(cards []Card) bool {
	// Create a map for storing card values and their number of occurrences
	cardValueCounts := make(map[int]int)

	// Iterate over all the cards and count the number of occurrences of each value
	for _, card := range cards {
		if _, valueCounted := cardValueCounts[card.value]; valueCounted {
			cardValueCounts[card.value]++
		} else {
			cardValueCounts[card.value] = 1
		}
	}

	// If any value occurs 4 times, the hand is a four of a kind
	for _, count := range cardValueCounts {
		if count == 4 {
			return true
		}
	}

	// Otherwise, it is not
	return false
}

func isFullHouse(cards []Card) bool {
	// Create a map for storing card values and their number of occurrences
	cardValueCounts := make(map[int]int)

	// Iterate over all the cards and count the number of occurrences of each value
	for _, card := range cards {
		if _, valueCounted := cardValueCounts[card.value]; valueCounted {
			cardValueCounts[card.value]++
		} else {
			cardValueCounts[card.value] = 1
		}
	}

	containsPair := false
	containsThreeOfAKind := false

	// Check if a value occurs twice and another value occurs thrice
	for _, count := range cardValueCounts {
		if count == 2 {
			containsPair = true
		}

		if count == 3 {
			containsThreeOfAKind = true
		}
	}

	// If the hand contains both a pair and a three of a kind it is a full house, otherwise it is not
	return containsPair && containsThreeOfAKind
}

func isFlush(cards []Card) bool {
	suit := cards[0].suit

	for _, card := range cards {
		// If any of the cards have a different suit the hand isn't a flush
		if card.suit != suit {
			return false
		}
	}

	// If all cards have the same suit, the hand is flush
	return true
}

func isStraight(cards []Card) bool {
	// First sort the cards by value to make checking for consecutiveness a single pass over the cards
	sort.Slice(cards, func(i, j int) bool { return cards[i].value < cards[j].value })

	for i, card := range cards {
		// The cards are sorted, so if any card but the last isn't followed by one of consecutive value the hand
		// can't be a straight
		if i != len(cards)-1 {
			if cards[i+1].value != card.value+1 {
				return false
			}
		}
	}

	// If all the cards are consecutively valued, the hand is a straight
	return true
}

func isThreeOfAKind(cards []Card) bool {
	// Create a map for storing card values and their number of occurrences
	cardValueCounts := make(map[int]int)

	// Iterate over all the cards and count the number of occurrences of each value
	for _, card := range cards {
		if _, valueCounted := cardValueCounts[card.value]; valueCounted {
			cardValueCounts[card.value]++
		} else {
			cardValueCounts[card.value] = 1
		}
	}

	// Check if any card value occurs 3 times
	for _, count := range cardValueCounts {
		if count == 3 {
			return true
		}
	}

	// Otherwise the hand does not contain a three of a kind
	return false
}

func isTwoPair(cards []Card) bool {
	// Create a map for storing card values and their number of occurrences
	cardValueCounts := make(map[int]int)

	// Iterate over all the cards and count the number of occurrences of each value
	for _, card := range cards {
		if _, valueCounted := cardValueCounts[card.value]; valueCounted {
			cardValueCounts[card.value]++
		} else {
			cardValueCounts[card.value] = 1
		}
	}

	numPairs := 0
	// Check how many values have two occurrences
	for _, count := range cardValueCounts {
		if count == 2 {
			numPairs++
		}
	}

	// If there are two values that occur twice in the hand, it is a two pair
	return numPairs == 2
}

func isPair(cards []Card) bool {
	// Create a map for storing card values and their number of occurrences
	cardValueCounts := make(map[int]int)

	// Iterate over all the cards and count the number of occurrences of each value
	for _, card := range cards {
		if _, valueCounted := cardValueCounts[card.value]; valueCounted {
			cardValueCounts[card.value]++
		} else {
			cardValueCounts[card.value] = 1
		}
	}

	// If any value occurs twice the hand contains a pair
	for _, count := range cardValueCounts {
		if count == 2 {
			return true
		}
	}

	// Otherwise it does not
	return false
}
