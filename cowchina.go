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
	cardType   CType
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

type CType int64

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
	currentBid := Bid{}
	for currentBid.wins == 0 {
		var highestBid, biddingTeam, biddingSymbol string
		inputFmt.Print("Enter highest bid: ")
		fmt.Scanln(&highestBid)
		inputFmt.Print("Enter highest bidding team: ")
		fmt.Scanln(&biddingTeam)
		inputFmt.Print("Enter bidding suit chosen: ")
		fmt.Scanln(&biddingSymbol)
		bidInt, err := strconv.Atoi(highestBid)

		var bidTeam Team
		if err == nil {
			switch biddingTeam {
			case "a":
				bidTeam = TEAM_A
			case "b":
				bidTeam = TEAM_B
			default:
				bidTeam = NULL
			}
			if bidTeam != NULL {
				symbol, err := getSymbolFromText(biddingSymbol)
				if err == nil {
					currentBid = Bid{bidInt, symbol, bidTeam}
				}

			}
		}
		if currentBid.wins == 0 {
			errorFmt.Print("Invalid bid!")
			fmt.Println()
		}
	}
	fmt.Println(currentBid)

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
		var cardType CType = NULL
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
				cardTypeGet, err := getSuitFromText(cardTypeIn, "")
				if err != nil {
					errorFmt.Print(err)
					fmt.Println()
				} else {
					cardType = cardTypeGet
				}
			} else if len(moveIn) == 3 {
				cardTypeIn := fmt.Sprintf("%c", moveIn[1])
				cardTypeIn2 := fmt.Sprintf("%c", moveIn[2])
				cardTypeGet, err := getSuitFromText(cardTypeIn, cardTypeIn2)
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
			// Invalid
			return
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

func getSuitFromText(suitA, suitB string) (CType, error) {
	fmt.Println(suitA, suitB)
	if suitB != "" {
		switch suitA {
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
	} else {
		suitBint, err := strconv.Atoi(suitB)
		if err != nil {
			return NULL, errors.New("Unknown Suit")
		}
		switch suitBint {
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
	}

}
