package main

import (
	"regexp"
	"strings"
	"fmt"
)
func cleanInput(text string) []string{
	text = strings.ToLower(text)

	re := regexp.MustCompile(`[^a-z0-9\s]+`)
	text = re.ReplaceAllString(text, "")
	words := strings.Fields(text)
	return words
}

func main() {
	result := cleanInput("Hello, World! Go is GREAT.")
	formatted := "[\"" + strings.Join(result, "\", \"") + "\"]"
	fmt.Println(formatted)
}