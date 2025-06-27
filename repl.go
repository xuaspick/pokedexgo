package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	previousUrl string
	nextUrl     string
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	supportedCommands := getCommands()

	cfg := config{
		nextUrl: "https://pokeapi.co/api/v2/location-area/",
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		prompt := cleanInput(scanner.Text())
		if len(prompt) == 0 {
			continue
		}
		promptCommand := prompt[0]

		if _, ok := supportedCommands[promptCommand]; !ok {
			fmt.Println("Unkwon command")
			continue
		}

		if err := supportedCommands[promptCommand].callback(&cfg); err != nil {
			fmt.Printf("Error executing command '%s': %v\n", promptCommand, err)
			continue
		}
	}
}

func cleanInput(text string) []string {
	var slicedString []string
	var splitted = strings.Split(strings.Trim(strings.ToLower(text), " "), " ")
	slicedString = append(slicedString, splitted...)
	return slicedString
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, data := range getCommands() {
		fmt.Printf("%s: %s\n", data.name, data.description)
	}
	return nil
}

func commandMap(c *config) error {
	locAreas, err := GetLocationAreas("forward", c)
	if err != nil {
		return err
	}
	for _, l := range locAreas {
		fmt.Println(l.Name)
	}
	return nil
}

func commandMapb(c *config) error {
	if c.previousUrl == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locAreas, err := GetLocationAreas("back", c)
	if err != nil {
		return err
	}
	for _, l := range locAreas {
		fmt.Println(l.Name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the names of the previous 20 location areas in the pokemon world",
			callback:    commandMapb,
		},
	}
}
