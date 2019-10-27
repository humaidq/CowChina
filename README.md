# CowChina

## 1. Description
CowChina is a logger for a variant of the Spades playing card game. It logs the moves, invalid cards (cheating) and winners.  

This could only be used with our variant of Spades, called Hokm.  

## 2. Requirements

The following packages must be installed on your system.

- Go
- Git

## 3. Copying and contributing

This program is written by Humaid AlQassimi, and is distributed
under the [MIT](https://humaidq.ae/license/mit) license.  


## 4. Download and install

```sh
$ go get -u git.sr.ht/~humaid/cowchina
$ go install git.sr.ht/~humaid/cowchina
```

## 5. Usage

To run the program:
```sh
$ cowchina
```

To begin, the names of the four players, the
highest bid, bidding suit chosen, and highest bidding team is
entered.  

Then the each card is inputed. Below is a cheat sheet of possible
inputs.

### 5.1. Input cheat sheet

| Input | Description |
| ---- | ----- |
| `ca` | Ace of Clubs |
| `da` | Ace of Diamonds |
| `ha` | Ace of Hearts |
| `sa` | Ace of Spades |
| `c{6-9}` | Clubs 6 to 9 |
| `s{6-9}` | Diamonds 6 to 9 |
| `h{6-9}` | Hearts 6 to 9 |
| `s{6-9}` | Spades 6 to 9 |
| `c{j,q,k}` | Clubs {jack, queen, king} |
| `d{j,q,k}` | Diamonds {jack, queen, king} |
| `h{j,q,k}` | Hearts {jack, queen, king} |
| `s{j,q,k}` | Spades {jack, queen, king} |
| `z` | Black Joker |
| `x` | Red Joker |
| `{enter}` | PASS |
| `{tab}` | Go back **(TODO)** |
| `:<id>` | Jump to player **(TODO)** |

