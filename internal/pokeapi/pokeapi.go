package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/xuaspick/pokedexgo/internal/pokecache"
)

type reslocationArea struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous *string        `json:"previous"`
	Results  []locationArea `json:"results"`
}

type locationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Config struct {
	NextUrl     string
	PreviousUrl string
}

type Client struct {
	cache  *pokecache.Cache
	Config Config
}

func NewClient(cacheInterval time.Duration) *Client {
	return &Client{
		cache: pokecache.NewCache(cacheInterval),
		Config: Config{
			NextUrl: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		},
	}
}

func (cli *Client) GetLocationAreas(direction string) ([]locationArea, error) {
	resp := reslocationArea{}
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

	if cached, ok := cli.cache.Get(url); ok {
		err := json.Unmarshal(cached, &resp)
		if err != nil {
			return resp.Results, err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return resp.Results, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return resp.Results, err
		}
		cli.cache.Add(url, body)

		err = json.Unmarshal(body, &resp)
		if err != nil {
			return resp.Results, err
		}
	}

	cli.Config.NextUrl = resp.Next
	if resp.Previous != nil {
		cli.Config.PreviousUrl = *resp.Previous
	} else {
		cli.Config.PreviousUrl = ""
	}

	return resp.Results, nil

}
