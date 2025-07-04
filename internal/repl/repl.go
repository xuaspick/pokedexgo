package repl

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/xuaspick/pokedexgo/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Client, ...string) error
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	supportedCommands := getCommands()

	pokeClient := pokeapi.NewClient(5 * time.Second)

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

		if err := supportedCommands[promptCommand].callback(pokeClient, prompt[1:]...); err != nil {
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

func commandExit(cli *pokeapi.Client, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cli *pokeapi.Client, args ...string) error {
	commands := getCommands()
	names := make([]string, 0, len(commands))

	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, name := range names {
		fmt.Printf("%s: %s\n", name, commands[name].description)
	}
	return nil
}

func commandMap(cli *pokeapi.Client, args ...string) error {
	locAreas, err := cli.GetLocationAreas("forward")
	if err != nil {
		return err
	}
	for _, l := range locAreas {
		fmt.Println(l.Name)
	}
	return nil
}

func commandMapb(cli *pokeapi.Client, args ...string) error {
	locAreas, err := cli.GetLocationAreas("back")
	if err != nil {
		return err
	}

	for _, l := range locAreas {
		fmt.Println(l.Name)
	}
	return nil
}

func commandExplore(cli *pokeapi.Client, args ...string) error {
	if len(args) == 0 {
		fmt.Println("Location name must be provided after command 'explore'")
		return nil
	}
	pokemonFound, err := cli.GetPokemonInArea(args[0])
	if err != nil {
		return err
	}
	for _, pokemonName := range pokemonFound {
		fmt.Println(pokemonName)
	}
	return nil
}

func commandCatch(cli *pokeapi.Client, args ...string) error {
	if len(args) == 0 {
		fmt.Println("A Pokemon name must be provided to attempt catching")
		return nil
	}
	_, err := cli.CatchPokemon(args[0])
	if err != nil {
		return err
	}
	return nil
}

func commandInspect(cli *pokeapi.Client, args ...string) error {
	if len(args) == 0 {
		fmt.Println("A Pokemon name must be provided to display data")
		return nil
	}

	cli.InspectPokemon(args[0])
	return nil
}

func commandPokedex(cli *pokeapi.Client, args ...string) error {
	err := cli.ListCaughtPokemon()
	return err
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attepmts to catch a pokemon",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore <location_name>",
			description: "Shows the pokemon in the indicated location",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "Shows the data of a captured pokemon",
			callback:    commandInspect,
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
		"pokedex": {
			name:        "pokedex",
			description: "shows a list of all caught pokemon",
			callback:    commandPokedex,
		},
	}
}
