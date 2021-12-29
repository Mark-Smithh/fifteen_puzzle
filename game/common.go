package game

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// CreateBoard creates the gameboard
func CreateBoard() (map[int]string, int, []Square) {
	gameBoard := make(map[int]string)
	shuffled := shufNums([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	for x := 1; x <= 15; x++ {
		num := shuffled[x-1]
		gameBoard[x] = strconv.Itoa(num)
	}
	emptySquare = randomNumber()
	shiftBoard(gameBoard, emptySquare, "*")
	sortedGameSquares = sortBoardByValue(gameBoard)
	return gameBoard, emptySquare, sortedGameSquares
}

func shufNums(nums []int) []int {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	return nums
}

func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 15
	return rand.Intn(max-min+1) + min
}

func sortBoardByValue(gameBoard map[int]string) []Square {
	keys := make([]int, 0, len(gameBoard))
	for k := range gameBoard {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	sortedBoard := []Square{}
	for _, v := range keys {
		valI, _ := strconv.Atoi(gameBoard[v])
		sortedBoard = append(sortedBoard, Square{square: v, value: valI})
	}

	sort.SliceStable(sortedBoard, func(i, j int) bool { return sortedBoard[i].value < sortedBoard[j].value })
	return sortedBoard
}

//Help return help text
func Help() string {
	return "Move the * around the gameboard.  It can only be moved to an adjacent square.  Use /puzzle endpoint to start a new game.  Use /move endpoint to make a move."
}
