package main

import (
	"regexp"
	"strings"
	"fmt"
)
func CleanInput(text string) []string{
	text = strings.ToLower(text)

	re := regexp.MustCompile(`[^a-z0-9\s]+`)
	text = re.ReplaceAllString(text, "")

	words := strings.Fields(text)

	return words
}

func main() {
	result := CleanInput("Hello, World! Go is GREAT.")
	fmt.Println(result)
}