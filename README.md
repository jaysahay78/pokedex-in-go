# Pokedex CLI

A small interactive command-line Pokedex written in Go. It uses the [PokeAPI](https://pokeapi.co/) to browse Pokemon location areas, explore Pokemon encounters, catch Pokemon, and inspect Pokemon you have already caught.

## Features

- Browse location areas in pages of 20
- Move forward and backward through the location list
- Explore Pokemon encounters in a selected location area
- Attempt to catch Pokemon with a probability based on base experience
- View and inspect Pokemon caught during the current session
- Cache API responses in memory

## Requirements

- Go 1.26.2 or newer
- Internet access for PokeAPI requests

## Getting Started

Run the CLI from the project root:

```powershell
go run ./main
```

Or build and run it:

```powershell
go build -o pokedex.exe ./main
.\pokedex.exe
```

## Commands

Once the prompt appears, use these commands:

```text
help
exit
map
mapb
explore <location_area>
catch <pokemon_name>
pokedex
inspect <pokemon_name>
```

### Example Session

```text
Pokedex > map
canalave-city-area
eterna-city-area
...

Pokedex > explore canalave-city-area
tentacool
tentacruel
...

Pokedex > catch tentacool
Throwing a Pokeball at tentacool...
tentacool was caught!

Pokedex > pokedex
Your pokedex:
- tentacool

Pokedex > inspect tentacool
Name: tentacool
Height: 9
Weight: 455
Stats:
 - hp : 40
...
```

## Testing

Run the test suite with:

```powershell
go test ./...
```

## Project Structure

```text
.
+-- go.mod
`-- main
    +-- main.go
    +-- commands.go
    +-- cache_test.go
    `-- internal
        `-- pokecache
            `-- cache.go
```



## Notes

- Caught Pokemon are stored only in memory and reset when the CLI exits.
- Location areas must be loaded with `map` or `mapb` before they can be explored.
- API responses are cached for the current process to reduce repeated network requests.

### Blog

link to a blog I wrote about it on substack - https://voluptatibusasper467509.substack.com/p/building-a-pokedex-in-go-what-the
