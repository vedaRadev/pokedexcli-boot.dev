package main

import (
    "strings"
    "bufio"
    "os"
    "fmt"

    "github.com/vedaRadev/pokedexcli-boot.dev/commands"
)

func main() {
    fmt.Println("Welcome to the Pokedex! Type \"help\" to see commands. Press enter with no command to rerun the previous command.")
    scanner := bufio.NewScanner(os.Stdin)
    commandConfig := commands.InitCommandConfig()
    var lastCommand commands.CliCommand
    lastParams := []string {}
    for {
        fmt.Print("\nPokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())

        if len(input) == 0 {
            if lastCommand.Execute != nil {
                lastCommand.Execute(&commandConfig, lastParams...)
            }
        } else if command, exists := commands.GetCommand(input[0]); exists {
            params := input[1:]
            err := command.Execute(&commandConfig, params...)
            if err != nil { fmt.Println(err) }
            lastCommand = command
            lastParams = params
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
