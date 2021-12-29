package main

import (
	"fmt"
	"net/http"
	"os"
	"puzzle/game"
	"strings"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "local" {
		localGame()
		return
	} else if len(os.Args) >= 2 && strings.ToLower(os.Args[1]) == "help" {
		fmt.Println(game.Help())
		return
	}
	apiGame()
}

func localGame() {
	gameBoard, _, _ := game.CreateBoard()
	game.StartGame(os.Stdout, gameBoard)
}

func apiGame() {
	http.HandleFunc("/puzzle", game.StartGameHTTP) //localhost:8080/puzzle
	http.HandleFunc("/move", game.MakeMoveHTTP)    //localhost:8080/move?square=1
	http.HandleFunc("/help", game.HelpTxt)         //localhost:8080/help
	http.ListenAndServe(":8080", nil)
}

// func hello(w http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintf(w, "hello\n")
// }
