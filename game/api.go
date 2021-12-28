package game

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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
