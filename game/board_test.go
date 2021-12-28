package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
generated gameboard looks like below
5 3 7 4
12 8 1 9
13 11 14 2
10 15 * 6
*/
func genGameBoard() map[int]string {
	gameBoard := make(map[int]string)
	gameBoard[1] = "5"
	gameBoard[2] = "3"
	gameBoard[3] = "7"
	gameBoard[4] = "4"
	gameBoard[5] = "12"
	gameBoard[6] = "8"
	gameBoard[7] = "1"
	gameBoard[8] = "9"
	gameBoard[9] = "13"
	gameBoard[10] = "11"
	gameBoard[11] = "14"
	gameBoard[12] = "2"
	gameBoard[13] = "10"
	gameBoard[14] = "15"
	gameBoard[15] = "*"
	gameBoard[16] = "6"
	return gameBoard
}

func TestValidMove(t *testing.T) {
	emptySquare := 15 //*
	moveRequested := "6"
	gameBoard := genGameBoard()
	gameSquares := sortBoardByValue(gameBoard)
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, gameSquares)
	assert.True(t, valid)
	assert.Equal(t, 16, moveRequestedLocation)
}

func TestNotValidMove(t *testing.T) {
	emptySquare := 15 //*
	moveRequested := "1"
	gameBoard := genGameBoard()
	gameSquares := sortBoardByValue(gameBoard)
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, gameSquares)
	assert.False(t, valid)
	assert.Equal(t, -1, moveRequestedLocation)
}

func TestValidMove1(t *testing.T) {
	emptySquare := 15 //*
	moveRequested := "14"
	gameBoard := genGameBoard()
	gameSquares := sortBoardByValue(gameBoard)
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, gameSquares)
	assert.True(t, valid)
	assert.Equal(t, 11, moveRequestedLocation)
}

func TestValidMove2(t *testing.T) {
	emptySquare := 15 //*
	moveRequested := "15"
	gameBoard := genGameBoard()
	gameSquares := sortBoardByValue(gameBoard)
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, gameSquares)
	assert.True(t, valid)
	assert.Equal(t, 14, moveRequestedLocation)
}
