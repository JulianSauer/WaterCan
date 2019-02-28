package hue

import (
    "os"
    "encoding/json"
)

type cache struct {
    BridgeIp string
    User     string
}

const CACHE = "cache.json"
var cacheContent *cache = nil

func loadCache() (*cache, error) {
    if cacheContent != nil {
        return cacheContent, nil
    }
    file, e := os.Open(CACHE)
    if e != nil {
        return nil, e
    }

    defer file.Close()
    decoder := json.NewDecoder(file)
    cache := cache{}
    e = decoder.Decode(&cache)
    cacheContent = &cache
    return &cache, e
}

func saveCache(cache *cache) error {
    file, e := os.Create(CACHE)
    if e != nil {
        return e
    }

    encoder := json.NewEncoder(file)
    return encoder.Encode(*cache)
}
