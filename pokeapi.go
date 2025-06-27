package main

import (
	"encoding/json"
	"net/http"
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

func GetLocationAreas(direction string, c *config) ([]locationArea, error) {
	resp := reslocationArea{}
	var url string

	if direction == "forward" {
		url = c.nextUrl
	} else {
		url = c.previousUrl
	}

	res, err := http.Get(url)
	if err != nil {
		return resp.Results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return resp.Results, err
	}
	c.nextUrl = resp.Next
	if resp.Previous != nil {
		c.previousUrl = *resp.Previous
	} else {
		c.previousUrl = ""
	}
	return resp.Results, nil

}
