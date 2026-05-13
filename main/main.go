package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/main/internal/pokecache"
	"regexp"
	"strings"
	"time"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)

	re := regexp.MustCompile(`[^a-z0-9\s-]+`)
	text = re.ReplaceAllString(text, "")
	words := strings.Fields(text)
	return words
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		cache:    pokecache.NewCache(5 * time.Minute),
		areaIDs:  make(map[string]int),
		pokemons: make(map[string]Pokemon),
	}

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}
		commandName := words[0]
		args := words[1:] //we used 1: here to get []string{"canalave-city-area"} here as only 1 would give us a string instead of slice

		/*map lookup in go
		command -> retrieved value
		exists -> boolean indicating if key exists
		value,exists := map[key]*/

		command, exists := commands[commandName]

		if exists {
			err := command.callback(cfg, args)

			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}
