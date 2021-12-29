package game

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var allGamesMap map[string]gameContext = make(map[string]gameContext)

//Square used when sorting the squares
type Square struct {
	square int
	value  int
}

type board struct {
	GameBoard []string `json:"gameboard"`
}

type requestedMove struct {
	Move string `json:"move"`
}

type moveResult struct {
	Msg       string `json:"msg"`
	GameBoard board  `json:"board"`
}

type gameContext struct {
	GameBoardMap      map[int]string `json:"gameboard_map"`
	EmptySquare       int            `json:"empty_square"`
	SortedGameSquares []Square       `json:"game_square"`
	GameBoard         []string       `json:"game_board"`
	Board             board          `json:"board"`
}

type missingQSParam struct {
	Msg string `json:"msg"`
}

type missingGameContext struct {
	Msg string `json:"msg"`
}

type help struct {
	Msg string `json:"msg"`
}

// StartGameHTTP start an api game
func StartGameHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gameBoard, emptySquare, sortedGameSquares := CreateBoard()
	gC := createGameContext(req, w, gameBoard, emptySquare, sortedGameSquares)
	if gC.GameBoard == nil {
		return
	}
	json.NewEncoder(w).Encode(gC.Board)
}

func createGameContext(req *http.Request, w http.ResponseWriter, gameBoard map[int]string, emptySquare int, sortedGameSquares []Square) gameContext {
	board := board{}

	found, user := paramCheck(req, w, "user", "Username must be passed as a querystring paramater.  Example localhost:8080/puzzle?user=mark")
	if !found {
		return gameContext{}
	}

	gContext := gameContext{}

	found, existingGC := existingGC(user)
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

func flipSpacesHTTP(gc gameContext, moveRequested string, moveRequestedLocation int, user string) gameContext {
	gameBoard := gc.GameBoardMap
	startShiftSpace := gc.EmptySquare

	gameBoard[startShiftSpace] = moveRequested
	gameBoard[moveRequestedLocation] = "*"

	board := board{}
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

//MakeMoveHTTP make an api game move
func MakeMoveHTTP(w http.ResponseWriter, req *http.Request) {
	found, user := paramCheck(req, w, "user", "User must be passed as a querystring paramater.  Example localhost:8080/puzzle?user=mark")
	if !found {
		return
	}

	found, existingGC := existingGC(user)
	if !found {
		m := missingGameContext{
			Msg: "Cannot find game context for user " + user + ".  A new game must be started before a move can be made.",
		}
		json.NewEncoder(w).Encode(m)
		return
	}

	found, moveRequested := paramCheck(req, w, "square", "Square must be passed as a querystring paramater.  Example localhost:8080/move?user=mark&square=4")
	if !found {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	moveRequest := requestedMove{}
	result := moveResult{Msg: "valid move"}
	moveRequest.Move = moveRequested

	valid, moveRequestedLocation := validMoveHTTP(moveRequest, user)
	if !valid {
		result.Msg = "invalid move"
		result.GameBoard = existingGC.Board
		json.NewEncoder(w).Encode(result)
		return
	}

	updatedGC := flipSpacesHTTP(existingGC, moveRequested, moveRequestedLocation, user)

	result.GameBoard = updatedGC.Board
	json.NewEncoder(w).Encode(result)
}

func validMoveHTTP(moveRequested requestedMove, user string) (bool, int) {
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
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // needed to not escape problematic HTML characters should be inside JSON quoted strings.  added to not change & to \u0026
	v, ok := req.URL.Query()[param]
	if !ok {
		mParam := missingQSParam{
			Msg: msg,
		}
		enc.Encode(mParam)
		return false, ""
	}
	return true, v[0]
}

func existingGC(user string) (bool, gameContext) {
	existingGC, found := allGamesMap[user]
	if found {
		return true, existingGC
	}
	return false, gameContext{}
}

//HelpTxt return help txt
func HelpTxt(w http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // needed to not escape problematic HTML characters should be inside JSON quoted strings.
	h := help{
		Msg: Help(),
	}
	enc.Encode(h)
}
