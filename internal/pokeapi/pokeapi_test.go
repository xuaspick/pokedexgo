package pokeapi

import (
	"slices"
	"testing"
	"time"
)

func TestNewClientLocs(t *testing.T) {
	client := NewClient(5 * time.Second)

	if len(client.Config.NextUrl) == 0 {
		t.Errorf("NextUrl expected to be non empty")
	}

	locs, _ := client.GetLocationAreas("back")

	if loclen := len(locs); loclen > 0 {
		t.Errorf("len of locations should be 0, %v found", loclen)
	}

	locs, _ = client.GetLocationAreas("forward")

	if loclen := len(locs); loclen != 20 {
		t.Errorf("len of locations should be 20, %v found", loclen)
	}
}

func TestPokemonInLocation(t *testing.T) {
	client := NewClient(5 * time.Second)

	pokemonFound, err := client.GetPokemonInArea("canalave-city-area")
	if err != nil {
		t.Errorf("slice of pokemon expected")
	}

	pokemonToFind := "tentacool"
	if !slices.Contains(pokemonFound, pokemonToFind) {
		t.Errorf("expected to find a %s", pokemonToFind)
	}
}

func TestCatchPokemon(t *testing.T) {
	cli := NewClient(5 * time.Second)

	captured, err := cli.CatchPokemon("caterpie", 1)
	if err != nil {
		t.Errorf("error found catching pokemon %v", err)
	}

	if captured {
		_, ok := cli.Pokedex["caterpie"]
		if !ok {
			t.Errorf("caterpie expected to be in Pokedex after capture")
		}
	}

}
