package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var emptySquare int
var gameSquares []gameSquare

type gameSquare struct {
	square int
	value  int
}

func CreateBoard() {
	numTracker := make(map[int]int)
	gameBoard := make(map[int]string)
	for x := 1; x <= 15; x++ {
		num := uniqueNum(numTracker, gameBoard)
		gameBoard[x] = strconv.Itoa(num)
	}

	emptySquare = randomNumber()
	shiftBoard(gameBoard, emptySquare, "*")

	fmt.Printf("%s\n\n", "Welcome! Here is your board:")
	printBoard(gameBoard)

	gameSquares = sortBoardByValue(gameBoard)

	for {
		playGame(gameBoard)
	}
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

func printBoard(gameBoard map[int]string) {
	for x := 1; x <= 16; x++ {
		if x != 0 && x%4 == 0 {
			fmt.Printf("%s,\t\n", gameBoard[x])
		} else {
			fmt.Printf("%s,\t", gameBoard[x])
		}
	}
	fmt.Println()
}

func shiftBoard(gameBoard map[int]string, startShiftSpace int, shiftChar string) {
	gameBoardCP := make(map[int]string)
	for k, v := range gameBoard {
		gameBoardCP[k] = v
	}
	for x := startShiftSpace; x <= 16; x++ {
		if x == startShiftSpace {
			gameBoard[x] = shiftChar
		} else {
			gameBoard[x] = gameBoardCP[x-1]
		}
	}
}

func flipSpaces(gameBoard map[int]string, startShiftSpace int, shiftChar string, moveRequestedLocation int) {
	gameBoardCP := make(map[int]string)
	for k, v := range gameBoard {
		gameBoardCP[k] = v
	}
	gameBoard[startShiftSpace] = shiftChar
	gameBoard[moveRequestedLocation] = "*"
	emptySquare = moveRequestedLocation

	gameSquares = sortBoardByValue(gameBoard)
}

func makeMove(emptySquare int) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What would you like to do?")
	moveRequested, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("got error " + err.Error())
	}
	return strings.TrimSuffix(moveRequested, "\n") //ReadString returns value plus the delimiter.  So need to trim the delimiter off
}

func playGame(gameBoard map[int]string) {
	moveRequested := makeMove(emptySquare)
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, gameSquares)
	if !valid {
		fmt.Printf("\n%s\n\n", "Invalid move! Here is your board:")
		printBoard(gameBoard)
	} else {
		fmt.Printf("\n%s\n\n", "Good Move! Here is your changed game board!")
		flipSpaces(gameBoard, emptySquare, moveRequested, moveRequestedLocation)
		printBoard(gameBoard)
	}
}

func validMove(emptySquare int, moveRequested string, gameSquares []gameSquare) (bool, int) {
	var moveRequestedLocation int

	moveRequestedInt, _ := strconv.Atoi(moveRequested)
	moveRequestedLocation = findSquare(gameSquares, moveRequestedInt)

	if emptySquare == gameSquares[moveRequestedLocation].square-4 || emptySquare == gameSquares[moveRequestedLocation].square+4 {
		return true, gameSquares[moveRequestedLocation].square
	}

	if emptySquare == gameSquares[moveRequestedLocation].square-1 || emptySquare == gameSquares[moveRequestedLocation].square+1 {
		return true, gameSquares[moveRequestedLocation].square
	}

	return false, -1
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

func findSquare(sortedBoard []gameSquare, moveRequestedLocation int) int {
	//use binay search to find the square
	i := sort.Search(len(sortedBoard), func(i int) bool {
		return moveRequestedLocation <= sortedBoard[i].value
	})
	return i
}

// func genTestBoard() map[int]string {
// 	gameBoard := make(map[int]string)
// 	gameBoard[1] = "3"
// 	gameBoard[2] = "7"
// 	gameBoard[3] = "4"
// 	gameBoard[4] = "12"
// 	gameBoard[5] = "8"
// 	gameBoard[6] = "1"
// 	gameBoard[7] = "9"
// 	gameBoard[8] = "11"
// 	gameBoard[9] = "14"
// 	gameBoard[10] = "*"
// 	gameBoard[11] = "10"
// 	gameBoard[12] = "15"
// 	gameBoard[13] = "13"
// 	gameBoard[14] = "6"
// 	gameBoard[15] = "5"
// 	gameBoard[16] = "2"
// 	return gameBoard
// }
