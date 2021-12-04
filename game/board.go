package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var emptySquare int

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
	for {
		playGame(emptySquare, gameBoard)
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

func playGame(emptySquare int, gameBoard map[int]string) {
	moveRequested := makeMove(emptySquare)
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, gameBoard)
	if !valid {
		fmt.Printf("\n%s\n\n", "Invalid move! Here is your board:")
		printBoard(gameBoard)
	} else {
		fmt.Printf("\n%s\n\n", "Good Move! Here is your changed game board!")
		flipSpaces(gameBoard, emptySquare, moveRequested, moveRequestedLocation)
		printBoard(gameBoard)
	}
}

func validMove(emptySquare int, moveRequested string, gameBoard map[int]string) (bool, int) {
	var moveRequestedLocation int
	for x, val := range gameBoard {
		if val == moveRequested {
			moveRequestedLocation = x
			break
		}
	}
	if emptySquare == moveRequestedLocation-1 || emptySquare == moveRequestedLocation+1 {
		return true, moveRequestedLocation
	}
	return false, -1
}
