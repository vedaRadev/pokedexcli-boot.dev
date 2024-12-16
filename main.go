package main

import (
    "strings"
    "bufio"
    "os"
    "fmt"
    "net/http"
    "encoding/json"
)

type LocationAreas struct {
    Count    int        `json:"count"`
    Next     *string    `json:"next"`
    Previous *string    `json:"previous"`
    Results  []struct {
        Name string     `json:"name"`
        URL  string     `json:"url"`
    } `json:"results"`
}

type cliCommand struct {
    name string
    description string
    callback func(*cliCommandConfig) error
}

type cliCommandConfig struct {
    nextLocationsPageUrl string
    prevLocationsPageUrl string
}

var COMMANDS map[string]cliCommand 

func commandExit(config *cliCommandConfig) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    // We'll never actually reach this return
    return nil
}

func commandHelp(config *cliCommandConfig) error {
    for _, command := range COMMANDS {
        fmt.Printf("%s: %s\n", command.name, command.description)
    }

    return nil
}

func getLocationAreas(config *cliCommandConfig, pageUrl string) error {
    res, err := http.Get(pageUrl)
    if err != nil { return err }
    defer res.Body.Close()

    if res.StatusCode < 200 || res.StatusCode > 299 {
        return fmt.Errorf("request failure - response code %s", res.Status)
    }

    var data LocationAreas
    if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
        return err
    }

    if data.Next != nil {
        config.nextLocationsPageUrl = *data.Next;
    } else {
        config.nextLocationsPageUrl = "";
    }

    if data.Previous != nil {
        config.prevLocationsPageUrl = *data.Previous;
    } else {
        config.prevLocationsPageUrl = "";
    }

    for _, location := range data.Results {
        fmt.Println(location.Name)
    }

    return nil
}

func commandMap(config *cliCommandConfig) error {
    if config.nextLocationsPageUrl == "" {
        fmt.Println("You are already on the last page")
        return nil
    }

    return getLocationAreas(config, config.nextLocationsPageUrl)
}

func commandMapB(config *cliCommandConfig) error {
    if config.prevLocationsPageUrl == "" {
        fmt.Println("You are already on the first page")
        return nil
    }

    return getLocationAreas(config, config.prevLocationsPageUrl)
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
        "map": {
            name: "map",
            description: "Display the next 20 locations in the Pokemon world",
            callback: commandMap,
        },
        "mapb": {
            name: "mapb",
            description: "Display the previous 20 locations in the Pokemon world",
            callback: commandMapB,
        },
    }

    fmt.Println("Welcome to the Pokedex!")
    scanner := bufio.NewScanner(os.Stdin)
    commandConfig := cliCommandConfig {
        nextLocationsPageUrl: "https://pokeapi.co/api/v2/location-area",
        prevLocationsPageUrl: "",
    }
    var lastCommand cliCommand
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())

        if len(input) == 0 {
            if lastCommand.callback != nil {
                lastCommand.callback(&commandConfig)
            }

            continue
        }

        if command, exists := COMMANDS[input[0]]; exists {
            command.callback(&commandConfig)
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
