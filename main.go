package main

import (
	"net/http"
	"os"
	"puzzle/game"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "local" {
		localGame()
	} else {
		apiGame()
	}
}

func localGame() {
	gameBoard, _, _ := game.CreateBoard()
	game.StartGame(os.Stdout, gameBoard)
}

func apiGame() {
	http.HandleFunc("/puzzle", game.PuzzleHttp) //localhost:8080/puzzle
	http.HandleFunc("/move", game.MakeMoveHttp) //localhost:8080/move?square=1
	http.ListenAndServe(":8080", nil)
}

// func hello(w http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintf(w, "hello\n")
// }
