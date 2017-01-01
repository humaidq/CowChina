# CowChina [![Travis](https://img.shields.io/travis/hmksq/CowChina.svg)]() [![Codecov](https://img.shields.io/codecov/c/github/hmksq/CowChina.svg)]()
CowChina is a logger for a variant of the Spades playing card game. It logs the moves, invalid cards (cheating) and winners. It is licensed under the [MIT License].  
Right now, this could only be used with our variant of Spades, called Hokm.

## Hokm
There are four players, each player is the partner with the player in front of them. The game and deals go anti-clockwise.  
A deck of 52 cards is used, listed from highest to lowest: Big (red) Joker, Small (black) Joker, A, K, Q, J, 10, 9, 8, 7, 6. Other cards are not included. Each player will have 9 cards.  
Bidding usually has a minimum of 6, the highest bidder will choose the Hokm (suit).  
The rest of the game is like Spades, but the suit will be chosen by the highest bidder than being default to spades.  

## Contributing
We appreciate your contributions! Make sure your code is run through `golint` before commiting.
