package main

import (
	"fmt"
	"net/http"
	"puzzle/game"
)

func main() {
	http.HandleFunc("/puzzle", game.PuzzleHttp) //localhost:8080/puzzle
	http.HandleFunc("/move", game.MakeMoveHttp) //localhost:8080/move?square=1
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}

// func main() {
// 	gameBoard, _, _ := game.CreateBoard()
// 	game.StartGame(os.Stdout, gameBoard)
// }

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
