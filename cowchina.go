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
sa: Spades (black) a

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

type Player struct {
	id      int
	name    string
	cSBList map[int]CSymbol
	team    Team
}

type Move struct {
	card   Card
	player int
}

var moves map[int]Move

type Card struct {
	cardSymbol CSymbol
	cardType   CSuit
}

type CSymbol int64

const (
	CLUBS    = iota
	DIAMONDS = iota
	HEARTS   = iota
	SPADES   = iota
	JOKER    = iota
	NULL     = 99
)

type CSuit int64

const (
	ACE   = iota
	JACK  = iota
	QUEEN = iota
	KING  = iota
	RED   = iota
	BLACK = iota
	N6    = iota
	N7    = iota
	N8    = iota
	N9    = iota
	N10   = iota
)

type Team int64

const (
	TEAM_A = iota
	TEAM_B = iota
)

var player1 Player
var player2 Player
var player3 Player
var player4 Player

type Bid struct {
	wins       int
	cardSymbol CSymbol
	team       Team
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
	player1 = Player{1, p1Name, make(map[int]CSymbol), TEAM_A}
	player2 = Player{2, p2Name, make(map[int]CSymbol), TEAM_B}
	player3 = Player{3, p3Name, make(map[int]CSymbol), TEAM_A}
	player4 = Player{4, p4Name, make(map[int]CSymbol), TEAM_B}
	infoFmt.Println(player1.name, "and", player3.name, "is Team A.")
	infoFmt.Println(player2.name, "and", player4.name, "is Team B.")
	// Get the bid
	currentBid := Bid{}
	for currentBid.wins == 0 {
		var highestBid, biddingTeam, biddingSymbol string
		var bidTeam Team
		// Get input
		inputFmt.Print("Enter highest bid: ")
		fmt.Scanln(&highestBid)
		inputFmt.Print("Enter highest bidding team: ")
		fmt.Scanln(&biddingTeam)
		inputFmt.Print("Enter bidding suit chosen: ")
		fmt.Scanln(&biddingSymbol)

		// Convert the bid input to string
		bidInt, err := strconv.Atoi(highestBid)
		if err == nil {
			// Get the input of the team
			switch biddingTeam {
			case "a":
				bidTeam = TEAM_A
			case "b":
				bidTeam = TEAM_B
			default:
				bidTeam = NULL
			}
			if bidTeam != NULL {
				// Get the symbol
				symbol, err := getSymbolFromText(biddingSymbol)
				if err == nil {
					// No errors, lets set the currentBid!
					currentBid = Bid{bidInt, symbol, bidTeam}
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
	var deckFirstSymbol Card
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
		var cardSymbol CSymbol = NULL
		var cardType CSuit = NULL
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

		} else if strings.HasPrefix(moveIn, "c") || strings.HasPrefix(moveIn, "d") || strings.HasPrefix(moveIn, "h") || strings.HasPrefix(moveIn, "s") {
			cardSymbolIn := fmt.Sprintf("%c", moveIn[0])
			cardSymbolGet, err := getSymbolFromText(cardSymbolIn)
			if err != nil {
				errorFmt.Print(err)
				fmt.Println()
			} else {
				cardSymbol = cardSymbolGet
			}

			if len(moveIn) == 2 {
				cardTypeIn := fmt.Sprintf("%c", moveIn[1])
				cardTypeGet, err := getSuitFromText(cardTypeIn)
				if err != nil {
					errorFmt.Print(err)
					fmt.Println()
				} else {
					cardType = cardTypeGet
				}
			} else if len(moveIn) == 3 {
				cardTypeIn := fmt.Sprintf("%c", moveIn[1]) + fmt.Sprintf("%c", moveIn[2])

				cardTypeGet, err := getSuitFromText(cardTypeIn)
				if err != nil {
					errorFmt.Print(err)
					fmt.Println()
				} else {
					cardType = cardTypeGet
				}
			}

		} else if moveIn == "z" {
			cardSymbol = JOKER
			cardType = RED
		} else if moveIn == "x" {
			cardSymbol = JOKER
			cardType = BLACK
		} else {

		}
		if cardSymbol != NULL && cardType != NULL && !pass {
			if playerTurn < 4 {
				playerTurn++
			} else {
				playerTurn = 1
			}
			cardPlaced := Card{cardSymbol, cardType}
			thisMove := Move{cardPlaced, currentPlayer.id}

			// Check with previous moves if this move changes symbol
			if !newFace {
				if deckFirstSymbol.cardSymbol == cardPlaced.cardSymbol {
					// This is compatible
				} else {
					// The player placed a card not the same symbol as the first. We will add it to blacklist.
					blacklist := currentPlayer.cSBList
					blacklist[len(blacklist)] = cardSymbol
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
				deckFirstSymbol = Card{}
				newFace = true
				infoFmt.Print("===== New deck! =====")
				fmt.Println()
			}
			moves[moveCount] = thisMove
			moveCount++
		} else if pass {
			infoFmt.Println(currentPlayer.name + " Pass!")
		} else {
			errorFmt.Print("Invalid input, enter again!")
			fmt.Println()
		}
	}
	fmt.Println("End of game! (or max cycles reached)")
}

func cardToText(card Card) string {
	var symbol, suit string
	if card.cardSymbol == JOKER {
		switch card.cardType {
		case BLACK:
			return "Small (black) Joker"
		case RED:
			return "Big (red) Joker"
		}
	}
	switch card.cardSymbol {
	case CLUBS:
		symbol = "Clubs"
	case DIAMONDS:
		symbol = "Diamonds"
	case HEARTS:
		symbol = "Hearts"
	case SPADES:
		symbol = "Spades"
	}

	switch card.cardType {
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
	}
	return suit + " of " + symbol

}

func getSymbolFromText(symbol string) (CSymbol, error) {
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
		return NULL, errors.New("Unknown Symbol")
	}
}

func getSuitFromText(suit string) (CSuit, error) {
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
			return NULL, errors.New("Unknown Suit")
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
			return NULL, errors.New("Unknown Suit")
		}
	}

}
