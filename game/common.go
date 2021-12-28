package game

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

func CreateBoard() (map[int]string, int, []gameSquare) {
	numTracker := make(map[int]int)
	gameBoard := make(map[int]string)
	for x := 1; x <= 15; x++ {
		num := uniqueNum(numTracker, gameBoard)
		gameBoard[x] = strconv.Itoa(num)
	}
	emptySquare = randomNumber()
	shiftBoard(gameBoard, emptySquare, "*")
	sortedGameSquares = sortBoardByValue(gameBoard)
	return gameBoard, emptySquare, sortedGameSquares
}

func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 15
	return rand.Intn(max-min+1) + min
}

func duplicateNum(numTracker map[int]int, num int) bool {
	return numTracker[num] > 1
}

func uniqueNum(numTracker map[int]int, gameBoard map[int]string) int {
	for {
		num := randomNumber()
		numTracker[num]++
		if !duplicateNum(numTracker, num) {
			return num
		}
	}
}

func sortBoardByValue(gameBoard map[int]string) []gameSquare {
	keys := make([]int, 0, len(gameBoard))
	for k := range gameBoard {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	sortedBoard := []gameSquare{}
	for _, v := range keys {
		valI, _ := strconv.Atoi(gameBoard[v])
		sortedBoard = append(sortedBoard, gameSquare{square: v, value: valI})
	}

	sort.SliceStable(sortedBoard, func(i, j int) bool { return sortedBoard[i].value < sortedBoard[j].value })
	return sortedBoard
}
