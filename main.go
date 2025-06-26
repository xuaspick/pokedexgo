package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var prompt string
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		prompt = scanner.Text()
		fmt.Printf("Your command was: %s \n", cleanInput(prompt)[0])
		fmt.Print("Pokedex > ")
	}

}

func cleanInput(text string) []string {
	var slicedString []string
	var splitted = strings.Split(strings.Trim(strings.ToLower(text), " "), " ")
	slicedString = append(slicedString, splitted...)
	return slicedString
}

func commandError() {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}
