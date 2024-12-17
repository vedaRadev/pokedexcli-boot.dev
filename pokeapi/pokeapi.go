package pokeapi

import (
    "time"
    "net/http"
    "fmt"
    "io"
    "encoding/json"
    "strings"

    "github.com/vedaRadev/pokedexcli-boot.dev/pokecache"
)

var requestCache *pokecache.Cache = pokecache.NewCache(10 * time.Second)

func Get(url string) ([]byte, int, error) {
    var result []byte

    if value, exists := requestCache.Get(url); exists {
        result = value
    } else {
        res, err := http.Get(url)
        if err != nil { return nil, 0, err }
        defer res.Body.Close()

        if res.StatusCode < 200 || res.StatusCode > 299 {
            return nil, res.StatusCode, fmt.Errorf("request failure - response code %s", res.Status)
        }

        result, err = io.ReadAll(res.Body)
        if err != nil {
            return nil, 0, err
        }

        requestCache.Add(url, result)
    }

    return result, 0, nil
}

func GetLocationAreaDetails(areaName string) (*LocationAreaDetails, error) {
    jsonData, statusCode, err := Get("https://pokeapi.co/api/v2/location-area/" + areaName)
    if statusCode == 404 { return nil, fmt.Errorf("Area %s not found", areaName) }
    if err != nil { return nil, err }

    var locationDetails LocationAreaDetails
    if err := json.Unmarshal(jsonData, &locationDetails); err != nil { return nil, err }
    return &locationDetails, nil
}

func GetPokemon(pokemonName string) (*PokemonDetails, error) {
    pokemonName = strings.ToLower(pokemonName)
    jsonData, statusCode, err := Get("https://pokeapi.co/api/v2/pokemon/" + pokemonName)
    if statusCode == 404 { return nil, fmt.Errorf("%s is not a pokemon", pokemonName) }
    if err != nil { return nil, err }

    var pokemon PokemonDetails
    if err := json.Unmarshal(jsonData, &pokemon); err != nil { return nil, err }
    return &pokemon, nil
}
