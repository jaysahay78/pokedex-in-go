package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"pokedexcli/main/internal/pokecache"
	"strconv"
	"strings"
	"time"
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
	areaIDs  map[string]int
	pokemons map[string]Pokemon
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

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
	Stats []struct {
		BaseState int `json:"base_stat"`
		Stat      struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
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

	"catch": {
		name:        "catch <pokemon_name>",
		description: "Throws a pokeball at a specific pokemon nearby to catch it",
		callback:    commandCatch,
	},

	"pokedex": {
		name:        "pokedex",
		description: "shows all the pokemon caught",
		callback:    commandPokdex,
	},

	"inspect": {
		name:        "inspect",
		description: "inspect pokemons that you have caught",
		callback:    commandInspect,
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

help: 					 Displays a help message
exit: 					 Exit the Pokedex
map:  					 Displays next 20 location areas
mapb: 					 Displays previous 20 location areas
explore <location_area>: Displays all the pokemon encounters in the specific location area
catch <pokemon_name>: 	 Throws a pokeball at a specific pokemon nearby to catch it
inspect: 	 			 Displays details about pokemons you have caught 
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

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide pokemon name")
		return nil
	}
	name := args[0]

	url := fmt.Sprintf(
		"https://pokeapi.co/api/v2/pokemon/%s/",
		name,
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

	var pokemon Pokemon
	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}
	_, exists := cfg.pokemons[pokemon.Name]

	if exists {
		fmt.Println("You already caught this pokemon!")
		return nil
	}
	baseExperience := pokemon.BaseExperience

	fmt.Println("Throwing a Pokeball at " + pokemon.Name + "...")
	minExp := 36
	maxExp := 608

	normalized := float64(baseExperience-minExp) / float64(maxExp-minExp)

	// harder pokemon => lower chance
	catchChance := 85 - int(normalized*70)

	// minimum catch chance
	if catchChance < 15 {
		catchChance = 15
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	roll := r.Intn(100)

	if roll < catchChance {
		fmt.Println(pokemon.Name + " was caught!")
		cfg.pokemons[pokemon.Name] = pokemon
	} else {
		fmt.Println(pokemon.Name + " escaped!")
	}

	return nil
}

func commandPokdex(cfg *config, args []string) error {

	if len(cfg.pokemons) == 0 {
		fmt.Println("No pokemons caught yet!")
	}

	fmt.Println("Your pokedex:")
	for _, pokemon := range cfg.pokemons {
		fmt.Println("- " + pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide pokemon name")
		return nil
	}
	name := args[0]
	var pokemon Pokemon

	if _, exists := cfg.pokemons[name]; exists {
		pokemon = cfg.pokemons[name]
	} else {
		fmt.Println("Unknown pokemon")
	}

	fmt.Println("Name: " + pokemon.Name)
	fmt.Println("Height: ", pokemon.Height)
	fmt.Println("Weight: ", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Println(" -", stat.Stat.Name, ": ", stat.BaseState)
	}
	fmt.Println("Types:")
	for _, types := range pokemon.Types {
		fmt.Println(" -" + types.Type.Name)
	}

	return nil
}
