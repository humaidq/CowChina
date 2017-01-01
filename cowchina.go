/*
The MIT License (MIT)

Copyright (c) 2016-2017 Sheikh Humaid AlQassimi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

/*
Possible inputs:

ca: Club (black) A
da: Diamond (red) A
ha: Heart (red) A
sa: Spades (black) A

c{6-9}: Club (black) {2-10}
d{6-9}: Diamond (red) {2-10}
h{6-9}: Heart (red) {2-10}
s{6-9}: Spades (black) {2-10}

c{j, q, k}
d{j, q, k}
h{j, q, k}
s{j, q, k}

z: Black JOKER
x: Red JOKER

{empty} = PASS
{tab} = Go back
:<player id> = Jump to player

Cheat Detection:
1. If the symbol is changed, blacklist the type before it.
  If the player uses the type blacklisted later, he is a cheater.
2. Burned Joker detection

*/

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Player is a struct that contains all the information of a player.
type Player struct {
	id      int
	name    string
	cSBList map[int]CSuit
	team    Team
}

// Move is a struct that contains a Card and the player (which placed it).
type Move struct {
	card   Card
	player int
}

var moves map[int]Move

// Card is a struct that contains all the attributes of a playing card.
type Card struct {
	cardSuit   CSuit
	cardSymbol CSymbol
}

// CSuit is a type that holds the symbol constant.
type CSuit int64

const (
	// CLUBS is the clubs symbol on a card.
	CLUBS = iota
	// DIAMONDS is the diamonds symbol on a card.
	DIAMONDS = iota
	// HEARTS is the hearts symbol on a card.
	HEARTS = iota
	// SPADES is the spades symbol on a card.
	SPADES = iota
	// JOKER is the joker card.
	JOKER = iota
	// NULL is a constant used to specify a null in an int64 type.
	NULL = 99
)

// CSymbol is a type that holds the suit const.
type CSymbol int64

const (
	// ACE is the "A" ace suit on a card.
	ACE = iota
	// JACK is the "J" jack suit on a card.
	JACK = iota
	// QUEEN is the "Q" queen suit on a card.
	QUEEN = iota
	// KING is the "K" king suit on a card.
	KING = iota
	// RED is a colour attribute of the joker card.
	RED = iota
	// BLACK is a colour attribute of the joker card.
	BLACK = iota
	// N6 is the number 6 on a card.
	N6 = iota
	// N7 is the number 7 on a card.
	N7 = iota
	// N8 is the number 8 on a card.
	N8 = iota
	// N9 is the number 9 on a card.
	N9 = iota
	// N10 is the number 10 on a card.
	N10 = iota
)

// Team is a type that can hold the team const.
type Team int64

const (
	// TeamA is a team of two players in a game.
	TeamA = iota
	// TeamB is a team of two players in a game.
	TeamB = iota
)

var player1 Player
var player2 Player
var player3 Player
var player4 Player

// Bid is a struct which can save the bid information of a game.
type Bid struct {
	wins     int
	cardSuit CSuit
	team     Team
}

func main() {
	errorFmt := color.New(color.BgRed).Add(color.Underline)
	inputFmt := color.New(color.FgCyan)
	infoFmt := color.New(color.FgMagenta)
	fmt.Println("Running CowChina v0.1.0 ")
	// Get player names
	var p1Name, p2Name, p3Name, p4Name string
	inputFmt.Print("Enter name of the first player: ")
	fmt.Scanln(&p1Name)
	inputFmt.Print("Enter name of the second player: ")
	fmt.Scanln(&p2Name)
	inputFmt.Print("Enter name of the third player: ")
	fmt.Scanln(&p3Name)
	inputFmt.Print("Enter name of the fourth player: ")
	fmt.Scanln(&p4Name)
	player1 = Player{1, p1Name, make(map[int]CSuit), TeamA}
	player2 = Player{2, p2Name, make(map[int]CSuit), TeamB}
	player3 = Player{3, p3Name, make(map[int]CSuit), TeamA}
	player4 = Player{4, p4Name, make(map[int]CSuit), TeamB}
	infoFmt.Println(player1.name, "and", player3.name, "is Team A.")
	infoFmt.Println(player2.name, "and", player4.name, "is Team B.")
	// Get the bid
	currentBid := Bid{}
	for currentBid.wins == 0 {
		var highestBid, biddingTeam, biddingSuit string
		var bidTeam Team
		// Get input
		inputFmt.Print("Enter highest bid: ")
		fmt.Scanln(&highestBid)
		inputFmt.Print("Enter bidding suit chosen: ")
		fmt.Scanln(&biddingSuit)
		inputFmt.Print("Enter highest bidding team: ")
		fmt.Scanln(&biddingTeam)

		// Convert the bid input to string
		bidInt, err := strconv.Atoi(highestBid)
		if err == nil {
			// Get the input of the team
			switch biddingTeam {
			case "a":
				bidTeam = TeamA
			case "b":
				bidTeam = TeamB
			default:
				bidTeam = NULL
			}
			if bidTeam != NULL {
				// Get the symbol
				suit, err := getSuitFromText(biddingSuit)
				if err == nil {
					// No errors, lets set the currentBid!
					currentBid = Bid{bidInt, suit, bidTeam}
				}
			}
		}
		if currentBid.wins == 0 {
			// Uh, oh! There was an error in the previous input and the Bid was not set
			errorFmt.Print("Invalid bid!")
			fmt.Println()
		}
	}
	fmt.Println(currentBid)
	// We will initialize the moves map and declare these variables to be used
	moves = make(map[int]Move)
	playerTurn := 1
	moveCount := 1
	deckCount := 1
	for i := 1; i < 99; i++ {
		var currentPlayer *Player
		switch playerTurn {
		case 1:
			currentPlayer = &player1
		case 2:
			currentPlayer = &player2
		case 3:
			currentPlayer = &player3
		case 4:
			currentPlayer = &player4
		default:
			errorFmt.Print("Error in finding the player for playerTurn! playerTurn invalid!")
			fmt.Println()
		}

		infoFmt.Print("It's " + currentPlayer.name + "'s turn!")
		fmt.Println()
		inputFmt.Print("Enter move's card (e.g. c9): ")
		var moveIn string
		newFace := false

		pass := false
		fmt.Scanln(&moveIn)
		var cardSuit CSuit
		var cardSymbol CSymbol
		var failed bool
		if moveIn == "" {
			if len(moves) == 0 {
				pass = true
				if playerTurn < 4 {
					playerTurn++
				} else {
					playerTurn = 1
				}
			} else {
				errorFmt.Print("Cannot pass mid-game!")
				fmt.Println()
			}
		} else if len(moveIn) < 2 {
			errorFmt.Print("Invalid input, enter again!")
			fmt.Println()
		} else if strings.HasPrefix(moveIn, "c") || strings.HasPrefix(moveIn, "d") || strings.HasPrefix(moveIn, "h") || strings.HasPrefix(moveIn, "s") {
			cardSuitIn := fmt.Sprintf("%c", moveIn[0])
			cardSuitGet, err := getSuitFromText(cardSuitIn)
			if err != nil {
				errorFmt.Print(err)
				fmt.Println()
			} else {
				cardSuit = cardSuitGet
			}

			cardSymbolIn := fmt.Sprintf("%c", moveIn[1])
			if len(moveIn) == 3 {
				cardSymbolIn += fmt.Sprintf("%c", moveIn[2])
			}
			cardSymbolGet, err := getSymbolFromText(cardSymbolIn)
			if err != nil {
				errorFmt.Print(err)
				fmt.Println()
			} else {
				cardSymbol = cardSymbolGet
			}
		} else if moveIn == "z" {
			cardSuit = JOKER
			cardSymbol = RED
		} else if moveIn == "x" {
			cardSuit = JOKER
			cardSymbol = BLACK
		}

		if pass {
			infoFmt.Println(currentPlayer.name + " Pass!")
		} else if !failed {
			if playerTurn < 4 {
				playerTurn++
			} else {
				playerTurn = 1
			}
			cardPlaced := Card{cardSuit, cardSymbol}
			thisMove := Move{cardPlaced, currentPlayer.id}

			// Check with previous moves if this move changes symbol
			if !newFace {
				if currentBid.cardSuit == cardPlaced.cardSuit {
					// This is compatible
				} else {
					// The player placed a card not the same symbol as the bid. We will add it to blacklist.
					blacklist := currentPlayer.cSBList
					blacklist[len(blacklist)] = cardSuit
					currentPlayer.cSBList = blacklist
				}
			}
			infoFmt.Println(cardToText(cardPlaced))
			// Continue for the next mov
			if deckCount < 4 {
				deckCount++
				newFace = false
			} else {
				deckCount = 1
				newFace = true
				infoFmt.Print("===== New deck! =====")
				fmt.Println()
			}
			moves[moveCount] = thisMove
			moveCount++
		} else {
			errorFmt.Print("Invalid input, enter again!")
			fmt.Println()
		}
	}
	fmt.Println("End of game! (or max cycles reached)")
}

func cardToText(card Card) string {
	var symbol, suit string
	if card.cardSuit == JOKER {
		switch card.cardSuit {
		case BLACK:
			return "Small (black) Joker"
		case RED:
			return "Big (red) Joker"
		}
	}
	switch card.cardSuit {
	case CLUBS:
		symbol = "Clubs"
	case DIAMONDS:
		symbol = "Diamonds"
	case HEARTS:
		symbol = "Hearts"
	case SPADES:
		symbol = "Spades"
	}

	switch card.cardSymbol {
	case ACE:
		suit = "Ace"
	case JACK:
		suit = "Jack"
	case QUEEN:
		suit = "Queen"
	case KING:
		suit = "King"
	case N6:
		suit = "6"
	case N7:
		suit = "7"
	case N8:
		suit = "8"
	case N9:
		suit = "9"
	case N10:
		suit = "10"
	}
	return suit + " of " + symbol
}

func getSuitFromText(symbol string) (CSuit, error) {
	switch symbol {
	case "c":
		return CLUBS, nil
	case "d":
		return DIAMONDS, nil
	case "h":
		return HEARTS, nil
	case "s":
		return SPADES, nil
	default:
		return NULL, errors.New("Unknown Suit")
	}
}

func getSymbolFromText(suit string) (CSymbol, error) {
	suitint, err := strconv.Atoi(suit)
	if err == nil {
		switch suitint {
		case 6:
			return N6, nil
		case 7:
			return N7, nil
		case 8:
			return N8, nil
		case 9:
			return N9, nil
		case 10:
			return N10, nil
		default:
			return NULL, errors.New("Unknown Symbol")
		}
	} else {
		switch suit {
		case "a":
			return ACE, nil
		case "j":
			return JACK, nil
		case "q":
			return QUEEN, nil
		case "k":
			return KING, nil
		default:
			return NULL, errors.New("Unknown Symbol")
		}
	}
}
