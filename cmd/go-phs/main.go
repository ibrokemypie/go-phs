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
	cards     []Card
	pokerHand PokerHand
}

type PokerHand struct {
	rank          int
	highCardValue int
}

type Card struct {
	value int
	suit  string
}

func main() {
	// STDIN Reader
	scanner := bufio.NewScanner(os.Stdin)

	playerOneScore := 0
	playerTwoScore := 0

	// Loop over lines from STDIN
	for scanner.Scan() {
		line := scanner.Text()

		handOne, handTwo, err := parseHands(line)
		// There shouldn't be an error when input is as expected, so if there is, log it and quit as we can't continue
		if err != nil {
			log.Fatal(err)
		}

		if handOne.compareHand(handTwo) == 1 {
			playerOneScore++
		} else if handOne.compareHand(handTwo) == -1 {
			playerTwoScore++
		} else {
			// A tie shouldn't happen in this data set, but handle it's possibility anyway
			playerOneScore++
			playerTwoScore++
		}
	}

	// Print the final scores of each player
	fmt.Printf("Player 1: %d\nPlayer 2: %d\n", playerOneScore, playerTwoScore)
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

	handOne.pokerHand = rankHand(handOne)
	handTwo.pokerHand = rankHand(handTwo)

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
			return -1, fmt.Errorf("Unable to convert %s to a value", char)
		}
	}

	return value, nil
}

func rankHand(hand Hand) PokerHand {
	pokerHand := PokerHand{
		rank:          0,
		highCardValue: 0,
	}

	// First sort the hand's cards by value
	sort.Slice(hand.cards, func(i, j int) bool { return hand.cards[i].value < hand.cards[j].value })

	// Next check the hand against each type of poker hand from highest rank to lowest to find the hand's
	// rank and the highest card that is part of the poker hand
	// (This isn't a switch statement so that we can get the hands high card as well)
	if isRoyalFlush, highCardValue := checkRoyalFlush(hand.cards); isRoyalFlush {
		pokerHand.rank = 10
		pokerHand.highCardValue = highCardValue
	} else if isStraightFlush, highCardValue := checkStraightFlush(hand.cards); isStraightFlush {
		pokerHand.rank = 9
		pokerHand.highCardValue = highCardValue
	} else if isFourOfAKind, highCardValue := checkFourOfAKind(hand.cards); isFourOfAKind {
		pokerHand.rank = 8
		pokerHand.highCardValue = highCardValue
	} else if isFullHouse, highCardValue := checkFullHouse(hand.cards); isFullHouse {
		pokerHand.rank = 7
		pokerHand.highCardValue = highCardValue
	} else if isFlush, highCardValue := checkFlush(hand.cards); isFlush {
		pokerHand.rank = 6
		pokerHand.highCardValue = highCardValue
	} else if isStraight, highCardValue := checkStraight(hand.cards); isStraight {
		pokerHand.rank = 5
		pokerHand.highCardValue = highCardValue
	} else if isThreeOfAKind, highCardValue := checkThreeOfAKind(hand.cards); isThreeOfAKind {
		pokerHand.rank = 4
		pokerHand.highCardValue = highCardValue
	} else if isTwoPair, highCardValue := checkTwoPair(hand.cards); isTwoPair {
		pokerHand.rank = 3
		pokerHand.highCardValue = highCardValue
	} else if isPair, highCardValue := checkPair(hand.cards); isPair {
		pokerHand.rank = 2
		pokerHand.highCardValue = highCardValue
	} else {
		// If the cards don't fit any poker hand, the rank is 1 (High Card) and the highest card is simply the
		// last one (as the cards are sorted by ascending value)
		pokerHand.rank = 1
		pokerHand.highCardValue = hand.cards[len(hand.cards)-1].value
	}

	return pokerHand
}

func checkRoyalFlush(cards []Card) (bool, int) {
	suit := cards[0].suit

	for _, card := range cards {
		// If any of the cards have a different suit the hand isn't a flush
		if card.suit != suit {
			return false, 0
		}

		// If any of the cards aren't a Ten, Jack, Queen, King or Ace this hand isn't a royal flush
		switch card.value {
		case 10, 11, 12, 13, 14:
			continue

		default:
			return false, 0
		}
	}

	// If the cards in the hand are all of the same suit and are the cards from Ten to Ace, then the hand is
	// a royal flush. The highest card in a royal flush is always the ace (14)
	return true, 14
}

func checkStraightFlush(cards []Card) (bool, int) {
	suit := cards[0].suit

	for i, card := range cards {
		// If any of the cards have a different suit the hand isn't a flush
		if card.suit != suit {
			return false, 0
		}

		// The cards are sorted, so if any card but the last isn't followed by one of consecutive value the hand
		// can't be a straight
		if i != len(cards)-1 {
			if cards[i+1].value != card.value+1 {
				return false, 0
			}
		}
	}

	// As a straight uses all five cards, the highest value is simply the last card in the hand
	return true, cards[len(cards)-1].value
}

func checkFourOfAKind(cards []Card) (bool, int) {
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

	// If any value occurs 4 times, the hand is a four of a kind, and return that value as the high card
	for value, count := range cardValueCounts {
		if count == 4 {
			return true, value
		}
	}

	// Otherwise, it is not
	return false, 0
}

func checkFullHouse(cards []Card) (bool, int) {
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
	threeOfAKindValue := 0

	// Check if a value occurs twice and another value occurs thrice
	for value, count := range cardValueCounts {
		if count == 2 {
			containsPair = true
		}

		if count == 3 {
			containsThreeOfAKind = true
			threeOfAKindValue = value
		}
	}

	// If the hand contains both a pair and a three of a kind it is a full house, otherwise it is not
	// Poker tiebreaker rules state that a tie of a full house is broken by the higher three of a kind value
	return containsPair && containsThreeOfAKind, threeOfAKindValue
}

func checkFlush(cards []Card) (bool, int) {
	suit := cards[0].suit

	for _, card := range cards {
		// If any of the cards have a different suit the hand isn't a flush
		if card.suit != suit {
			return false, 0
		}
	}

	// If all cards have the same suit, the hand is flush. Like a straight, a flush uses all five cards so the last
	// of the sorted cards in the hand is the high card for this hand
	return true, cards[len(cards)-1].value
}

func checkStraight(cards []Card) (bool, int) {
	for i, card := range cards {
		// The cards are sorted, so if any card but the last isn't followed by one of consecutive value the hand
		// can't be a straight
		if i != len(cards)-1 {
			if cards[i+1].value != card.value+1 {
				return false, 0
			}
		}
	}

	// If all the cards are consecutively valued, the hand is a straight. The last of the cards in the hand is
	// the high card of this hand.
	return true, cards[len(cards)-1].value
}

func checkThreeOfAKind(cards []Card) (bool, int) {
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

	// Check if any card value occurs 3 times, if so, its value is the high card for this hand.
	for value, count := range cardValueCounts {
		if count == 3 {
			return true, value
		}
	}

	// Otherwise the hand does not contain a three of a kind
	return false, 0
}

func checkTwoPair(cards []Card) (bool, int) {
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
	highestPair := 0
	// Check how many values have two occurrences
	for value, count := range cardValueCounts {
		if count == 2 {
			numPairs++
			if value > highestPair {
				highestPair = value
			}
		}
	}

	// If there are two values that occur twice in the hand, it is a two pair and return the pair's cards' value
	return numPairs == 2, highestPair
}

func checkPair(cards []Card) (bool, int) {
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

	// If any value occurs twice the hand contains a pair, return the cards' value
	for value, count := range cardValueCounts {
		if count == 2 {
			return true, value
		}
	}

	// Otherwise it does not
	return false, 0
}

// Compare two hands. -1 represents smaller, 0 represents equal, 1 represents greater
func (handOne Hand) compareHand(handTwo Hand) int {
	// First check if the rank is higher or lower
	if handOne.pokerHand.rank > handTwo.pokerHand.rank {
		return 1
	} else if handOne.pokerHand.rank < handTwo.pokerHand.rank {
		return -1
	}

	// If the rank is the same, use the poker hand's high card to break the tie
	if handOne.pokerHand.highCardValue > handTwo.pokerHand.highCardValue {
		return 1
	} else if handOne.pokerHand.highCardValue < handTwo.pokerHand.highCardValue {
		return -1
	}

	// If the rank and poker hand's high card both tie, go through each sorted hand backwards until one
	// card has a higher value than the other
	for i := 0; i < 5; i++ {
		if handOne.cards[4-i].value > handTwo.cards[4-i].value {
			return 1
		} else if handOne.cards[4-i].value < handTwo.cards[4-i].value {
			return -1
		}
	}

	// If the hands are exactly the same then it is a tie that cannot be broken
	return 0
}
