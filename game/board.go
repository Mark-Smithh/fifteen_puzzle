package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var emptySquare int
var sortedGameSquares []gameSquare

var allGamesMap map[string]GameContext = make(map[string]GameContext)

type gameSquare struct {
	square int
	value  int
}

type Board struct {
	GameBoard []string `json:"game_board"`
}

type RequestedMove struct {
	Move string `json:"move"`
}

type MoveResult struct {
	Msg       string `json:"msg"`
	GameBoard Board  `json:"board"`
}

type GameContext struct {
	GameBoardMap      map[int]string
	EmptySquare       int
	SortedGameSquares []gameSquare
	GameBoard         []string
	Board             Board
}

type MissingQSParam struct {
	Msg string `json:"msg"`
}

type MissingGameContext struct {
	Msg string `json:"msg"`
}

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

func StartGame(w io.Writer, gameBoard map[int]string) {
	fmt.Fprintf(w, "%s\n\n", "Welcome! Here is your board:")
	PrintBoard(w, gameBoard)
	for {
		playGame(w, gameBoard)
	}
}

func PuzzleHttp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gameBoard, emptySquare, sortedGameSquares := CreateBoard()
	gameContext := httpBoard(req, w, gameBoard, emptySquare, sortedGameSquares)
	if gameContext.GameBoard == nil {
		return
	}
	json.NewEncoder(w).Encode(gameContext.Board)
}

func httpBoard(req *http.Request, w http.ResponseWriter, gameBoard map[int]string, emptySquare int, sortedGameSquares []gameSquare) GameContext {
	board := Board{}

	found, user := paramCheck(req, w, "user", "Username must be passed as a querystring paramater.  Example localhost:8080/puzzle?user=mark")
	if !found {
		return GameContext{}
	}

	gContext := GameContext{}

	found, existingGC := gameContext(user)
	if found {
		board.GameBoard = existingGC.GameBoard
		return existingGC
	}

	allSquares := []string{}
	for x := 1; x <= 16; x++ {
		allSquares = append(allSquares, gameBoard[x])
	}
	board.GameBoard = allSquares

	gContext.GameBoard = allSquares
	gContext.GameBoardMap = gameBoard
	gContext.EmptySquare = emptySquare
	gContext.SortedGameSquares = sortedGameSquares
	gContext.Board = board
	allGamesMap[user] = gContext

	return gContext
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

// func PrintBoardJson(w io.Writer, gameBoard map[int]string) {
// 	for x := 1; x <= 16; x++ {
// 		if x != 0 && x%4 == 0 {
// 			fmt.Fprintf(w, "%s,\t\n", gameBoard[x])
// 		} else {
// 			fmt.Fprintf(w, "%s,\t", gameBoard[x])
// 		}
// 	}
// 	fmt.Fprint(w, "\n")
// }

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

func flipSpacesHttp(gc GameContext, moveRequested string, moveRequestedLocation int, user string) GameContext {
	gameBoard := gc.GameBoardMap
	startShiftSpace := gc.EmptySquare

	gameBoard[startShiftSpace] = moveRequested
	gameBoard[moveRequestedLocation] = "*"

	board := Board{}
	allSquares := []string{}
	for x := 1; x <= 16; x++ {
		allSquares = append(allSquares, gameBoard[x])
	}
	board.GameBoard = allSquares

	gc.EmptySquare = moveRequestedLocation
	gc.SortedGameSquares = sortBoardByValue(gameBoard)
	gc.GameBoard = allSquares
	gc.GameBoardMap = gameBoard
	gc.Board = board
	allGamesMap[user] = gc

	return gc
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

func MakeMoveHttp(w http.ResponseWriter, req *http.Request) {
	found, user := paramCheck(req, w, "user", "Username must be passed as a querystring paramater.  Example localhost:8080/puzzle?user=mark")
	if !found {
		return
	}

	found, moveRequested := paramCheck(req, w, "square", "Square must be passed as a querystring paramater.  Example localhost:8080/puzzle?user=mark&square=4") // NEED TO FIX & IS NOT COMING ACROSS AS EXPECTED
	if !found {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	moveRequest := RequestedMove{}
	result := MoveResult{Msg: "valid move"}
	moveRequest.Move = moveRequested

	found, existingGC := gameContext(user)
	if !found {
		m := MissingGameContext{
			Msg: "Cannot find game context for user " + user,
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	valid, moveRequestedLocation := validMoveHttp(moveRequest, user)
	if !valid {
		result.Msg = "invalid move"
		result.GameBoard = existingGC.Board
		json.NewEncoder(w).Encode(result)
		return
	}

	updatedGC := flipSpacesHttp(existingGC, moveRequested, moveRequestedLocation, user)

	result.GameBoard = updatedGC.Board
	json.NewEncoder(w).Encode(result)
}

func validMoveHttp(moveRequested RequestedMove, user string) (bool, int) {
	gContext := allGamesMap[user]
	sortedGameSquares := gContext.SortedGameSquares
	emptySquare := gContext.EmptySquare

	var moveRequestedLocation int
	moveRequestedInt, _ := strconv.Atoi(moveRequested.Move)
	moveRequestedLocation = findSquare(sortedGameSquares, moveRequestedInt)

	if emptySquare == sortedGameSquares[moveRequestedLocation].square-4 || emptySquare == sortedGameSquares[moveRequestedLocation].square+4 {
		return true, sortedGameSquares[moveRequestedLocation].square
	}

	if emptySquare == sortedGameSquares[moveRequestedLocation].square-1 || emptySquare == sortedGameSquares[moveRequestedLocation].square+1 {
		return true, sortedGameSquares[moveRequestedLocation].square
	}

	return false, -1
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

func validMove(emptySquare int, moveRequested string, sortedGameSquares []gameSquare) (bool, int) {
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

func paramCheck(req *http.Request, w http.ResponseWriter, param string, msg string) (bool, string) {
	v, ok := req.URL.Query()[param]
	if !ok {
		mParam := MissingQSParam{
			Msg: msg,
		}
		json.NewEncoder(w).Encode(mParam)
		return false, ""
	}
	return true, v[0]
}

func gameContext(user string) (bool, GameContext) {
	existingGC, found := allGamesMap[user]
	if found {
		return true, existingGC
	}
	return false, GameContext{}
}
