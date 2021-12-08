package main

import (
	"fmt"
	"net/http"
	"puzzle/game"
)

func main() {
	// http.HandleFunc("/puzzle", puzzle)
	// http.HandleFunc("/", hello)
	// http.ListenAndServe(":8080", nil)
	game.CreateBoard()
}

func puzzle(w http.ResponseWriter, req *http.Request) {
	// game.CreateBoard()
	fmt.Fprintf(w, "puzzle\n")
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
