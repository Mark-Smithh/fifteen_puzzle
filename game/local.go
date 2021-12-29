package game

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

var emptySquare int
var sortedGameSquares []Square

//StartGame start a local game
func StartGame(w io.Writer, gameBoard map[int]string) {
	fmt.Fprintf(w, "%s\n\n", "Welcome! Here is your board:")
	PrintBoard(w, gameBoard)
	for {
		playGame(w, gameBoard)
	}
}

//PrintBoard print the board
func PrintBoard(w io.Writer, gameBoard map[int]string) {
	for x := 1; x <= 16; x++ {
		if x != 0 && x%4 == 0 {
			fmt.Fprintf(w, "%s,\t\n", gameBoard[x])
		} else {
			fmt.Fprintf(w, "%s,\t", gameBoard[x])
		}
	}
	fmt.Fprint(w, "\n")
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
	gameBoard[startShiftSpace] = shiftChar
	gameBoard[moveRequestedLocation] = "*"
	emptySquare = moveRequestedLocation

	sortedGameSquares = sortBoardByValue(gameBoard)
}

func makeMove(w io.Writer, emptySquare int) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(w, "What would you like to do?\n")
	moveRequested, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("got error " + err.Error())
	}
	return strings.TrimSuffix(moveRequested, "\n") //ReadString returns value plus the delimiter.  So need to trim the delimiter off
}

func playGame(w io.Writer, gameBoard map[int]string) {
	moveRequested := makeMove(w, emptySquare)
	valid, moveRequestedLocation := validMove(emptySquare, moveRequested, sortedGameSquares)
	if !valid {
		fmt.Fprintf(w, "\n%s\n\n", "Invalid move! Here is your board:")
		PrintBoard(w, gameBoard)
	} else {
		fmt.Fprintf(w, "\n%s\n\n", "Good Move! Here is your changed game board!")
		flipSpaces(gameBoard, emptySquare, moveRequested, moveRequestedLocation)
		PrintBoard(w, gameBoard)
	}
}

func validMove(emptySquare int, moveRequested string, sortedGameSquares []Square) (bool, int) {
	var moveRequestedLocation int

	moveRequestedInt, _ := strconv.Atoi(moveRequested)
	moveRequestedLocation = findSquare(sortedGameSquares, moveRequestedInt)

	if emptySquare == sortedGameSquares[moveRequestedLocation].square-4 || emptySquare == sortedGameSquares[moveRequestedLocation].square+4 {
		return true, sortedGameSquares[moveRequestedLocation].square
	}

	if emptySquare == sortedGameSquares[moveRequestedLocation].square-1 || emptySquare == sortedGameSquares[moveRequestedLocation].square+1 {
		return true, sortedGameSquares[moveRequestedLocation].square
	}

	return false, -1
}

func findSquare(sortedBoard []Square, moveRequestedLocation int) int {
	//use binay search to find the square
	i := sort.Search(len(sortedBoard), func(i int) bool {
		return moveRequestedLocation <= sortedBoard[i].value
	})
	return i
}
