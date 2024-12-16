package main

import (
    "strings"
    "bufio"
    "os"
    "fmt"
)

type cliCommand struct {
    name string
    description string
    callback func() error
}

var COMMANDS map[string]cliCommand 


func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    // We'll never actually reach this return
    return nil
}


func commandHelp() error {
    for _, command := range COMMANDS {
        fmt.Printf("%s: %s\n", command.name, command.description)
    }

    return nil
}


func main() {
    // TODO Is there a way to initialize this outside of main without causing an initialization
    // cycle due to the help command callback referencing the global variable?
    COMMANDS = map[string]cliCommand {
        "exit": {
            name: "exit",
            description: "Exit the Pokedex",
            callback: commandExit,
        },
        "help": {
            name: "help",
            description: "Displays commands and their usage information",
            callback: commandHelp,
        },
    }

    fmt.Println("Welcome to the Pokedex!")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())
        if command, exists := COMMANDS[input[0]]; !exists {
            fmt.Printf("unrecognized command: %s\n", input[0])
        } else {
            command.callback()
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
