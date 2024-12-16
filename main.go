package main

import (
    "strings"
    "bufio"
    "os"
    "fmt"
)

func cleanInput(text string) []string {
    var cleaned []string
    for _, word := range strings.Fields(text) {
        cleaned = append(cleaned, strings.ToLower(word))
    }

    return cleaned
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())
        fmt.Printf("Your command was: %s\n", input[0])
    }
}
