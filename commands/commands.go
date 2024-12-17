package commands

import (
    "encoding/json"
    "os"
    "fmt"
    "errors"
    "math/rand/v2"
    "strings"

    "github.com/vedaRadev/pokedexcli-boot.dev/pokeapi"
)



type CliCommand struct {
    Name string
    Description string
    Execute func(*CliCommandConfig, ...string) error
}

// TODO if no file outside of this one needs to access CliCommandConfig fields, stop exporting them
type CliCommandConfig struct {
    NextLocationsPageUrl string
    PrevLocationsPageUrl string
    currentAreaName string
    // TODO Just store the name then fetch from network when inspecting?
    caughtPokemon map[string]*pokeapi.PokemonDetails
}

// TODO stop exporting if nothing outside of this package needs to touch this directly
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
        "explore": {
            // TODO need to add a Usage field?
            Name: "explore <area>",
            Description: "Enter and explore an area to find pokemon",
            Execute: commandExplore,
        },
        "catch": {
            Name: "catch <pokemon name>",
            Description: "Attempt to catch a pokemon and add it to your pokedex",
            Execute: commandCatch,
        },
        "area": {
            Name: "area",
            Description: "Print your current location",
            Execute: commandArea,
        },
        "inspect": {
            Name: "explore <pokemon name>",
            Description: "inspect a pokemon that you've caught",
            Execute: commandInspect,
        },
    }
}

func InitCommandConfig() CliCommandConfig {
    return CliCommandConfig {
        NextLocationsPageUrl: "https://pokeapi.co/api/v2/location-area",
        PrevLocationsPageUrl: "",
        currentAreaName: "",
        caughtPokemon: map[string]*pokeapi.PokemonDetails {},
    }
}

func GetCommand(name string) (CliCommand, bool) {
    command, exists := Commands[name]
    return command, exists
}

func commandExit(config *CliCommandConfig, params ...string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    // We'll never actually reach this return
    return nil
}

func commandHelp(config *CliCommandConfig, params ...string) error {
    for _, command := range Commands {
        fmt.Printf("%s: %s\n", command.Name, command.Description)
    }

    return nil
}

func printLocationAreas(config *CliCommandConfig, pageUrl string) error {
    jsonData, _, err := pokeapi.Get(pageUrl)
    if err != nil { return err }

    var data pokeapi.LocationAreasPaged
    if err := json.Unmarshal(jsonData, &data); err != nil { return err }

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

func commandMap(config *CliCommandConfig, params ...string) error {
    if config.NextLocationsPageUrl == "" {
        fmt.Println("You are already on the last page")
        return nil
    }

    return printLocationAreas(config, config.NextLocationsPageUrl)
}

func commandMapB(config *CliCommandConfig, params ...string) error {
    if config.PrevLocationsPageUrl == "" {
        fmt.Println("You are already on the first page")
        return nil
    }

    return printLocationAreas(config, config.PrevLocationsPageUrl)
}

func commandExplore(config *CliCommandConfig, params ...string) error {
    // TODO if no areaName provided, just explore the player's current area
    if len(params) == 0 {
        return errors.New("Command takes 1 argument: area_name")
    }

    areaName := params[0]
    locationDetails, err := pokeapi.GetLocationAreaDetails(areaName)
    if err != nil { return err }

    config.currentAreaName = areaName
    fmt.Printf("You enter %s and explore to find ", areaName)
    if len(locationDetails.PokemonEncounters) > 0 {
        fmt.Println()
        for _, encounter := range locationDetails.PokemonEncounters {
            fmt.Printf(" - %s\n", encounter.Pokemon.Name)
        }
    } else {
        fmt.Println("nothing.")
    }

    return nil
}

func commandCatch(config *CliCommandConfig, params ...string) error {
    if len(params) == 0 {
        return errors.New("Command takes 1 argument: pokemon_name")
    }

    locationDetails, err := pokeapi.GetLocationAreaDetails(config.currentAreaName)
    if err != nil { return err }

    pokemonName := strings.ToLower(params[0])
    pokemonIsInArea := false
    for _, encounter := range locationDetails.PokemonEncounters {
        if strings.ToLower(encounter.Pokemon.Name) == pokemonName {
            pokemonIsInArea = true
            break
        }
    }
    if !pokemonIsInArea { return fmt.Errorf("%s is not in your current area!", pokemonName) }

    pokemon, err := pokeapi.GetPokemon(pokemonName)
    if err != nil { return err }

    fmt.Printf("Throwing a pokeball at %s...\n", pokemonName)

    // TODO need WAY better catch determination
    didCatch := rand.IntN(pokemon.BaseExperience) < 50
    if didCatch {
        fmt.Printf("You caught %s!\n", pokemonName)
        config.caughtPokemon[pokemonName] = pokemon
    } else {
        fmt.Printf("%s got away!\n", pokemonName)
    }

    return nil
}

func commandArea(config *CliCommandConfig, params ...string) error {
    if config.currentAreaName == "" {
        fmt.Println("You have not entered an area yet")
    } else {
        fmt.Printf("You are in %s\n", config.currentAreaName)
    }
    
    return nil
}

func commandInspect(config *CliCommandConfig, params ...string) error {
    if len(params) == 0 {
        return errors.New("Command takes 1 argument: pokemon_name")
    }

    name := strings.ToLower(params[0])
    details, exists := config.caughtPokemon[name];
    if !exists { return fmt.Errorf("You have not caught a %s", name) }

    fmt.Printf("Name: %v\n", details.Name)
    fmt.Printf("Height: %v\n", details.Height)
    fmt.Printf("Weight: %v\n", details.Weight)
    fmt.Println("Stats:")
    for _, Stat := range details.Stats {
        fmt.Printf(" -%v: %v\n", Stat.Stat.Name, Stat.BaseStat)
    }
    fmt.Println("Types:")
    for _, Type := range details.Types {
        fmt.Printf(" -%v\n", Type.Type.Name)
    }

    return nil
}
