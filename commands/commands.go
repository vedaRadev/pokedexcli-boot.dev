package commands

import (
    "net/http"
    "encoding/json"
    "os"
    "fmt"
    "time"
    "io"
    "errors"

    "github.com/vedaRadev/pokedexcli-boot.dev/pokecache"
)

// TODO put pokeapi stuff into its own package?
type LocationAreas struct {
    Count    int        `json:"count"`
    Next     *string    `json:"next"`
    Previous *string    `json:"previous"`
    Results  []struct {
        Name string     `json:"name"`
        URL  string     `json:"url"`
    } `json:"results"`
}

type LocationAreaDetails struct {
    EncounterMethodRates []struct {
        EncounterMethod struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"encounter_method"`
        VersionDetails []struct {
            Rate    int `json:"rate"`
            Version struct {
                Name string `json:"name"`
                URL  string `json:"url"`
            } `json:"version"`
        } `json:"version_details"`
    } `json:"encounter_method_rates"`
    GameIndex int `json:"game_index"`
    ID        int `json:"id"`
    Location  struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"location"`
    Name  string `json:"name"`
    Names []struct {
        Language struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"language"`
        Name string `json:"name"`
    } `json:"names"`
    PokemonEncounters []struct {
        Pokemon struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"pokemon"`
        VersionDetails []struct {
            EncounterDetails []struct {
                Chance          int   `json:"chance"`
                ConditionValues []any `json:"condition_values"`
                MaxLevel        int   `json:"max_level"`
                Method          struct {
                    Name string `json:"name"`
                    URL  string `json:"url"`
                } `json:"method"`
                MinLevel int `json:"min_level"`
            } `json:"encounter_details"`
            MaxChance int `json:"max_chance"`
            Version   struct {
                Name string `json:"name"`
                URL  string `json:"url"`
            } `json:"version"`
        } `json:"version_details"`
    } `json:"pokemon_encounters"`
}

func pokeApiGet(url string) ([]byte, error) {
    var result []byte

    if value, exists := requestCache.Get(url); exists {
        result = value
    } else {
        res, err := http.Get(url)
        if err != nil { return nil, err }
        defer res.Body.Close()

        if res.StatusCode < 200 || res.StatusCode > 299 {
            return nil, fmt.Errorf("request failure - response code %s", res.Status)
        }

        result, err = io.ReadAll(res.Body)
        if err != nil {
            return nil, err
        }

        requestCache.Add(url, result)
    }

    return result, nil
}

type CliCommand struct {
    Name string
    Description string
    Execute func(*CliCommandConfig, ...string) error
}

// TODO if no file outside of this one needs to access CliCommandConfig fields, stop exporting them
type CliCommandConfig struct {
    NextLocationsPageUrl string
    PrevLocationsPageUrl string
}

// TODO stop exporting if nothing outside of this package needs to touch this directly
var Commands map[string]CliCommand
var requestCache *pokecache.Cache

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
            Description: "Explore an area to find pokemon",
            Execute: commandExplore,
        },
    }

    // TODO tune this interval
    requestCache = pokecache.NewCache(5 * time.Second)
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

func getLocationAreas(config *CliCommandConfig, pageUrl string) error {
    jsonData, err := pokeApiGet(pageUrl)
    if err != nil { return err }

    var data LocationAreas
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

    return getLocationAreas(config, config.NextLocationsPageUrl)
}

func commandMapB(config *CliCommandConfig, params ...string) error {
    if config.PrevLocationsPageUrl == "" {
        fmt.Println("You are already on the first page")
        return nil
    }

    return getLocationAreas(config, config.PrevLocationsPageUrl)
}

func commandExplore(config *CliCommandConfig, params ...string) error {
    if len(params) == 0 {
        return errors.New("Command takes 1 argument: area_name")
    }

    areaName := params[0]
    jsonData, err := pokeApiGet("https://pokeapi.co/api/v2/location-area/" + areaName)
    if err != nil { return err }

    var locationDetails LocationAreaDetails
    if err := json.Unmarshal(jsonData, &locationDetails); err != nil { return err }

    if len(locationDetails.PokemonEncounters) > 0 {
        fmt.Println("Found pokemon:")
        for _, encounter := range locationDetails.PokemonEncounters {
            fmt.Printf(" - %s\n", encounter.Pokemon.Name)
        }
    } else {
        fmt.Println("No pokemon found...")
    }

    return nil
}
