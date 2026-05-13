package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pokedexcli/main/internal/pokecache"
	"strconv"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type locationId struct {
	id    int
	cache *pokecache.Cache
}

type config struct {
	next     *string
	previous *string
	cache    *pokecache.Cache
	//pointer to string provides type safety in case of null being returned in the url
	areaIDs map[string]int
}

type locationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type exploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

const baseUrl = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokdex",
		callback:    commandExit,
	},

	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},

	"map": {
		name:        "map",
		description: "Displays next 20 location areas",
		callback:    commandMap,
	},

	"mapb": {
		name:        "mapb",
		description: "Displays previous 20 location areas",
		callback:    commandMapb,
	},

	"explore": {
		name:        "explore <location_area>",
		description: "Displays all the pokemon encounters in the specific location area",
		callback:    commandExplore,
	},
}

func commandExit(cfg *config, args []string) error {
	fmt.Print("Closing the Pokedex...Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Print(`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: Displays next 20 location areas
mapb: Displays previous 20 location areas
explore <location_area>: Displays all the pokemon encounters in the specific location area
`)

	return nil
}

func commandMap(cfg *config, args []string) error {

	url := baseUrl

	if cfg.next != nil {
		url = *cfg.next
	}

	var body []byte

	cachedData, found := cfg.cache.Get(url)

	if found {
		body = cachedData
	} else {

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cfg.cache.Add(url, body)
	}

	var locations locationAreaResponse

	err := json.Unmarshal(body, &locations)
	if err != nil {
		return err
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	for _, area := range locations.Results {
		parts := strings.Split(area.URL, "/")

		id, err := strconv.Atoi(parts[len(parts)-2])
		if err != nil {
			return err
		}

		cfg.areaIDs[area.Name] = id
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(cfg *config, args []string) error {
	url := baseUrl

	if cfg.previous != nil {
		url = *cfg.previous
	}

	var body []byte

	cachedData, found := cfg.cache.Get(url)

	if found {
		body = cachedData
	} else {

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		cfg.cache.Add(url, body)
	}

	var locations locationAreaResponse

	err := json.Unmarshal(body, &locations)
	if err != nil {
		return err
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	for _, area := range locations.Results {
		parts := strings.Split(area.URL, "/")

		id, err := strconv.Atoi(parts[len(parts)-2])
		if err != nil {
			return err
		}

		cfg.areaIDs[area.Name] = id
		fmt.Println(area.Name)
	}

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a location area")
		return nil
	}
	areaName := args[0]
	id, exists := cfg.areaIDs[areaName]

	if !exists {
		fmt.Println("Area not found")
		return nil
	}

	url := fmt.Sprintf(
		"https://pokeapi.co/api/v2/location-area/%d/",
		id,
	)

	var body []byte

	cachedData, found := cfg.cache.Get(url)

	if found {
		body = cachedData
	} else {

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, body)
	}

	var encounters exploreResponse
	err := json.Unmarshal(body, &encounters)
	if err != nil {
		return err
	}
	for _, encounter := range encounters.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}
