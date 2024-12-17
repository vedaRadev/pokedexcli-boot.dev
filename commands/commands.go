package commands

import (
    "net/http"
    "encoding/json"
    "os"
    "fmt"
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

type CliCommand struct {
    Name string
    Description string
    Execute func(*CliCommandConfig) error
}

// TODO if no file outside of this one needs to access CliCommandConfig fields, stop exporting them
type CliCommandConfig struct {
    NextLocationsPageUrl string
    PrevLocationsPageUrl string
}

var Commands map[string]CliCommand

func init() {
    Commands = map[string]CliCommand {
        "exit": {
            Name: "exit",
            Description: "Exit the Pokedex",
            Execute: commandExit,
        },
        "help": {
            Name: "help",
            Description: "Displays commands and their usage information",
            Execute: commandHelp,
        },
        "map": {
            Name: "map",
            Description: "Display the next 20 locations in the Pokemon world",
            Execute: commandMap,
        },
        "mapb": {
            Name: "mapb",
            Description: "Display the previous 20 locations in the Pokemon world",
            Execute: commandMapB,
        },
    }
}

func InitCommandConfig() CliCommandConfig {
    return CliCommandConfig {
        NextLocationsPageUrl: "https://pokeapi.co/api/v2/location-area",
        PrevLocationsPageUrl: "",
    }
}

func GetCommand(name string) (CliCommand, bool) {
    command, exists := Commands[name]
    return command, exists
}

func commandExit(config *CliCommandConfig) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    // We'll never actually reach this return
    return nil
}

func commandHelp(config *CliCommandConfig) error {
    for _, command := range Commands {
        fmt.Printf("%s: %s\n", command.Name, command.Description)
    }

    return nil
}

func getLocationAreas(config *CliCommandConfig, pageUrl string) error {
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
        config.NextLocationsPageUrl = *data.Next;
    } else {
        config.NextLocationsPageUrl = "";
    }

    if data.Previous != nil {
        config.PrevLocationsPageUrl = *data.Previous;
    } else {
        config.PrevLocationsPageUrl = "";
    }

    for _, location := range data.Results {
        fmt.Println(location.Name)
    }

    return nil
}

func commandMap(config *CliCommandConfig) error {
    if config.NextLocationsPageUrl == "" {
        fmt.Println("You are already on the last page")
        return nil
    }

    return getLocationAreas(config, config.NextLocationsPageUrl)
}

func commandMapB(config *CliCommandConfig) error {
    if config.PrevLocationsPageUrl == "" {
        fmt.Println("You are already on the first page")
        return nil
    }

    return getLocationAreas(config, config.PrevLocationsPageUrl)
}
