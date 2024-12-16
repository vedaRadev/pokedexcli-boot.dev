package main

import (
    "strings"
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
    fmt.Println("Hello, World!")
}
