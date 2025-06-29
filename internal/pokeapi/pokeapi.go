package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/xuaspick/pokedexgo/internal/pokecache"
)

type resLocationArea struct {
	Next     string         `json:"next"`
	Previous *string        `json:"previous"`
	Results  []locationArea `json:"results"`
}

type locationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type resPokemonPerArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name    string `json:"name"`
	BaseExp int    `json:"base_experience"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Stats   []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type Client struct {
	cache   *pokecache.Cache
	Config  Config
	Pokedex map[string]Pokemon
}

type Config struct {
	NextUrl     string
	PreviousUrl string
}

const (
	thresholdToCatch = 30
)

func NewClient(cacheInterval time.Duration) *Client {
	return &Client{
		cache: pokecache.NewCache(cacheInterval),
		Config: Config{
			NextUrl: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		},
		Pokedex: map[string]Pokemon{},
	}
}

func (cli *Client) processRequest(url string) (httpBody []byte, err error) {
	var httpReadBody []byte
	if cached, ok := cli.cache.Get(url); ok {
		return cached, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return httpReadBody, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return httpReadBody, err
	}
	cli.cache.Add(url, body)
	return body, nil

}

func (cli *Client) GetLocationAreas(direction string) ([]locationArea, error) {
	resp := resLocationArea{}
	var url string

	if direction == "forward" {
		url = cli.Config.NextUrl
	} else {
		url = cli.Config.PreviousUrl
	}

	if url == "" {
		fmt.Println("you're on the first page")
		return resp.Results, nil
	}

	body, err := cli.processRequest(url)
	if err != nil {
		return resp.Results, err
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return resp.Results, err
	}

	cli.Config.NextUrl = resp.Next
	if resp.Previous != nil {
		cli.Config.PreviousUrl = *resp.Previous
	} else {
		cli.Config.PreviousUrl = ""
	}

	return resp.Results, nil

}

func (cli *Client) GetPokemonInArea(cityName string) (pokemonNames []string, err error) {
	var pokemonFound []string
	pokeFound := resPokemonPerArea{}
	url := "https://pokeapi.co/api/v2/location-area/" + cityName
	body, err := cli.processRequest(url)

	if err != nil {
		return pokemonFound, err
	}
	err = json.Unmarshal(body, &pokeFound)
	if err != nil {
		return pokemonFound, err
	}
	for _, encounter := range pokeFound.PokemonEncounters {
		pokemonFound = append(pokemonFound, encounter.Pokemon.Name)
	}

	return pokemonFound, nil
}

func (cli *Client) CatchPokemon(pokemonName string, testMode ...int) (captured bool, err error) {
	var pokemonData Pokemon
	if pokemonName == "" {
		return false, nil
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	body, err := cli.processRequest(url)
	if err != nil {
		return false, nil
	}
	err = json.Unmarshal(body, &pokemonData)
	if err != nil {
		return false, nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	leTime := time.Now().UnixNano()
	r := rand.New(rand.NewSource(leTime))

	var capturePoints int
	if len(testMode) > 0 {
		capturePoints = pokemonData.BaseExp
	} else {
		capturePoints = r.Intn(pokemonData.BaseExp)
	}
	basePoints := pokemonData.BaseExp - thresholdToCatch

	// fmt.Printf("attemptin capture, BasePoints: %v <= CapturePoints: %v\n", basePoints, capturePoints)

	if basePoints > capturePoints {
		fmt.Printf("%s escaped!\n", pokemonName)
		return false, nil
	}

	cli.Pokedex[pokemonName] = pokemonData
	fmt.Printf("%s was caught!\n", pokemonName)

	return true, nil
}

func (cli *Client) InspectPokemon(pokemonName string) error {
	pokemon, ok := cli.Pokedex[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %v\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}
	return nil
}

func (cli *Client) ListCaughtPokemon() error {
	if len(cli.Pokedex) == 0 {
		fmt.Println("No pokemon caught, go catch 'em all!")
		return nil
	}

	fmt.Println("Your pokedex:")
	for _, pokemon := range cli.Pokedex {
		fmt.Printf(" - %s \n", pokemon.Name)
	}
	return nil
}
