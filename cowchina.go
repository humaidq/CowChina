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
	card      Card
	player    int
	blackMove bool
}

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

const (
	brand = "CowChina"
	ver   = "0.1.0"
)

var players map[int]Player
var moves map[int]Move

// Bid is a struct which can save the bid information of a game.
type Bid struct {
	wins     int
	cardSuit CSuit
	team     Team
}

func main() {
	errorFmt := color.New(color.BgRed).Add(color.Underline)
	anticheatFmt := color.New(color.BgHiMagenta)
	inputFmt := color.New(color.FgCyan)
	infoFmt := color.New(color.FgMagenta)
	//debugFmt := color.New(color.FgGreen).Add(color.Bold) // Not used in prod
	fmt.Println("Running " + brand + " v" + ver)
	// Get player names
	players = make(map[int]Player)
	for i := 1; i < 5; i++ {
		var name string
		inputFmt.Print("Enter name of the " + playerNumberToSimpleText(i) + " player: ")
		fmt.Scanln(&name)
		// TODO assign the team of each player
		players[i] = Player{i, name, make(map[int]CSuit), TeamA}
	}
	// Get the bid
	currentBid := Bid{}
	for currentBid == (Bid{}) {
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
		if err == nil && bidInt > 5 && bidInt < 11 {
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
		if currentBid == (Bid{}) {
			// Uh, oh! There was an error in the previous input and the Bid was not set
			errorFmt.Print("Invalid bid!")
			fmt.Println()
		}
	}
	// We will initialize the moves map and declare these variables to be used
	moves = make(map[int]Move)
	playerTurn := 1
	moveCount := 1
	var deckFirstCard Card
	deckMoves := make(map[int]Move)
	deckCount := 1
	for i := 1; i < 99; i++ {
		infoFmt.Println("It's " + players[playerTurn].name + "'s turn!")
		inputFmt.Print("Enter move's card (e.g. c9): ")
		var moveIn string

		pass := false
		fmt.Scanln(&moveIn)
		var cardSuit CSuit
		var cardSymbol CSymbol
		var failed bool
		if moveIn == "" {
			pass = true
			if playerTurn < 4 {
				playerTurn++
			} else {
				playerTurn = 1
			}
		} else if (strings.HasPrefix(moveIn, "c") || strings.HasPrefix(moveIn, "d") || strings.HasPrefix(moveIn, "h") || strings.HasPrefix(moveIn, "s")) && len(moveIn) > 1 && len(moveIn) < 4 {
			// TODO: Make the suit and symbol interchangeable
			cardSuitIn := fmt.Sprintf("%c", moveIn[0])
			cardSuitGet, err := getSuitFromText(cardSuitIn)
			if err != nil {
				errorFmt.Print(err)
				fmt.Println()
				failed = true
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
				failed = true
			} else {
				cardSymbol = cardSymbolGet
			}
		} else if moveIn == "z" {
			cardSuit = JOKER
			cardSymbol = RED
		} else if moveIn == "x" {
			cardSuit = JOKER
			cardSymbol = BLACK
		} else {
			errorFmt.Print("Invalid input, enter again!")
			fmt.Println()
			failed = true
		}

		if pass {
			infoFmt.Println(players[playerTurn].name + " Pass!")
		} else if !failed {

			cardPlaced := Card{cardSuit, cardSymbol}

			blackMove := false
			for _, blSuit := range players[playerTurn].cSBList {
				if cardPlaced.cardSuit == blSuit {
					blackMove = true
					anticheatFmt.Print(players[playerTurn].name + " cheated, " + players[playerTurn].name + " should have placed " + cardSuitToText(cardPlaced.cardSuit) + " in a previous deck.")
					fmt.Println()
					break
				}
			}
			thisMove := Move{cardPlaced, players[playerTurn].id, blackMove}
			// Check with previous moves if this move changes symbol
			if deckFirstCard != (Card{}) {
				if deckFirstCard.cardSuit == cardPlaced.cardSuit || cardPlaced.cardSuit == JOKER {
					// This is compatible
				} else {
					// The player placed a card not the same symbol as the first. We will add it to blacklist.
					players[playerTurn].cSBList[len(players[playerTurn].cSBList)] = deckFirstCard.cardSuit
					//anticheatFmt.Print(players[playerTurn].name + " placed " + cardSuitToText(cardPlaced.cardSuit) + " and we do not expect the player to have " + cardSuitToText(deckFirstCard.cardSuit) + ".")
					//fmt.Println()
				}
			} else if cardPlaced.cardSuit != JOKER {
				deckFirstCard = cardPlaced
			}
			infoFmt.Println(cardToText(cardPlaced))
			// Continue for the next move
			moves[moveCount] = thisMove
			deckMoves[len(deckMoves)] = thisMove
			moveCount++
			if deckCount < 4 {
				deckCount++
				if playerTurn < 4 {
					playerTurn++
				} else {
					playerTurn = 1
				}
			} else {
				// Get the winner
				winnerMove := getWinnerFromDeck(deckMoves)
				infoFmt.Println(players[winnerMove.player].name + " eats this deck!")
				playerTurn = winnerMove.player
				deckMoves = make(map[int]Move)
				deckCount = 1
				deckFirstCard = Card{}
				infoFmt.Println("========= New deck! ========")
			}
		} else {
			errorFmt.Print("Invalid input, enter again!")
			fmt.Println()
		}
	}
	fmt.Println("End of game! (or max cycles reached)")
}

func getWinnerFromDeck(deck map[int]Move) Move {
	// Check if there are jokers in the deck
	for _, move := range deck {
		if move.card.cardSuit == JOKER {
			if move.card.cardSymbol == RED {
				return move
			} else if move.card.cardSymbol == BLACK {
				return move
			}
		}
	}
	return deck[0]
}

func cardSuitToText(suit CSuit) string {
	switch suit {
	case CLUBS:
		return "Clubs"
	case DIAMONDS:
		return "Diamonds"
	case HEARTS:
		return "Hearts"
	case SPADES:
		return "Spades"
	case JOKER:
		return "Joker"
	}
	return "Unknown"
}

func cardSymbolToText(symbol CSymbol) string {
	switch symbol {
	case ACE:
		return "Ace"
	case JACK:
		return "Jack"
	case QUEEN:
		return "Queen"
	case KING:
		return "King"
	case N6:
		return "6"
	case N7:
		return "7"
	case N8:
		return "8"
	case N9:
		return "9"
	case N10:
		return "10"
	case BLACK:
		return "(black)"
	case RED:
		return "(red)"
	}
	return "Unknown"
}

func cardToText(card Card) string {
	if card.cardSuit == JOKER {
		switch card.cardSymbol {
		case BLACK:
			return "Small (black) Joker"
		case RED:
			return "Big (red) Joker"
		}
	}
	return cardSuitToText(card.cardSuit) + " of " + cardSymbolToText(card.cardSymbol)
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

func playerNumberToSimpleText(number int) string {
	switch number {
	case 1:
		return "first"
	case 2:
		return "second"
	case 3:
		return "third"
	case 4:
		return "fourth"
	default:
		return "unknown"
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
