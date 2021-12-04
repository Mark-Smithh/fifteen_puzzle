package hello

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Hello() {
	input := takeInput()
	fmt.Println(uCaseOutput(input))
}

func takeInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter text: ")
	textVal, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("got err %s", err.Error())
	}
	return textVal
}

func uCaseOutput(text string) string {
	return strings.ToUpper(text)
}
