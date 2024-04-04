package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// KeyValueEntry represents a key/value entry with metadata.
type KeyValueEntry struct {
    Key       string    `json:"key"`
    Value     string    `json:"value"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// KeyValueStore represents the in-memory key/value store.
type KeyValueStore struct {
    data map[string]KeyValueEntry
    mu   sync.RWMutex // Mutex for concurrent access
}

// NewKeyValueStore initializes a new KeyValueStore.
func NewKeyValueStore() *KeyValueStore {
    return &KeyValueStore{
        data: make(map[string]KeyValueEntry),
    }
}

// Set sets a key-value pair in the store.
func (kv *KeyValueStore) Set(key, value string) {
    kv.mu.Lock()
    defer kv.mu.Unlock()
    entry := KeyValueEntry{
        Key:       key,
        Value:     value,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    kv.data[key] = entry
}

// Get retrieves a value from the store based on the given key.
func (kv *KeyValueStore) Get(key string) (KeyValueEntry, bool) {
    kv.mu.RLock()
    defer kv.mu.RUnlock()
    entry, ok := kv.data[key]
    return entry, ok
}

func main() {
    // Create a new KeyValueStore instance
    store := NewKeyValueStore()

    // Define HTTP handlers
    http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Query().Get("key")
        value := r.URL.Query().Get("value")
        if key == "" || value == "" {
            http.Error(w, "Key or value not provided", http.StatusBadRequest)
            return
        }
        store.Set(key, value)
        fmt.Fprintf(w, "Key %s set to value %s\n", key, value)
    })

    http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Query().Get("key")
        if key == "" {
            http.Error(w, "Key not provided", http.StatusBadRequest)
            return
        }
        entry, ok := store.Get(key)
        if !ok {
            http.Error(w, "Key not found", http.StatusNotFound)
            return
        }
        jsonEntry, err := json.Marshal(entry)
        if err != nil {
            http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonEntry)
    })

    // Start HTTP server
    log.Fatal(http.ListenAndServe(":8080", nil))
}