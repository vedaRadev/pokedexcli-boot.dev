package main

import (
    "strings"
    "bufio"
    "os"
    "fmt"

    "github.com/vedaRadev/pokedexcli-boot.dev/commands"
)

func main() {
    fmt.Println("Welcome to the Pokedex!")
    scanner := bufio.NewScanner(os.Stdin)
    commandConfig := commands.InitCommandConfig()
    var lastCommand commands.CliCommand
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())

        if len(input) == 0 {
            if lastCommand.Callback != nil {
                lastCommand.Callback(&commandConfig)
            }

            continue
        }

        if command, exists := commands.Commands[input[0]]; exists {
            command.Callback(&commandConfig)
            lastCommand = command
        } else {
            fmt.Printf("unrecognized command: %s\n", input[0])
        }
    }
}


func cleanInput(text string) []string {
    var cleaned []string
    for _, word := range strings.Fields(text) {
        cleaned = append(cleaned, strings.ToLower(word))
    }
    return cleaned
}
