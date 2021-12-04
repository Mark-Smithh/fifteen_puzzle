package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, gameBoard)
	assert.True(t, valid)
	assert.Equal(t, 16, moveRequestedLocation)
}

func TestNotValidMove(t *testing.T) {
	emptySquare := 15 //*
	moveRequested := "1"
	gameBoard := genGameBoard()
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, gameBoard)
	assert.False(t, valid)
	assert.Equal(t, -1, moveRequestedLocation)
}