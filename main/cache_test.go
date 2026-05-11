package main

import (
	"pokedexcli/main/internal/pokecache"
	"testing"
	"time"
)

func TestCacheAddAndGet(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Minute)

	key := "https://pokeapi.co/api/v2/location-area/"
	value := []byte("test data")

	// add to cache
	cache.Add(key, value)

	// retrieve from cache
	cachedValue, found := cache.Get(key)

	if !found {
		t.Errorf("expected cache hit, got cache miss")
	}

	if string(cachedValue) != string(value) {
		t.Errorf("expected %s, got %s", value, cachedValue)
	}
}
