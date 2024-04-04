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
	if !ok {
		return entry, false
	}
	var value interface{}
	err := json.Unmarshal([]byte(entry.Value), &value)
	if err != nil {
		// If value is not a valid JSON, return the entry as is
		return entry, true
	}

	// If parsing succeeds, update the entry with the parsed JSON value
	// Convert the parsed JSON value back to a string
	jsonValue, err := json.Marshal(value)
	if err != nil {
		// If there's an error marshaling the value, return the entry as is
		return entry, true
	}
	entry.Value = string(jsonValue)
	return entry, true
}

// CleanupRoutine periodically removes entries older than 24 hours.
func (kv *KeyValueStore) CleanupRoutine() {
	for {
		time.Sleep(24 * time.Hour)
		kv.mu.Lock()
		for key, entry := range kv.data {
			if time.Since(entry.CreatedAt) > 24*time.Hour {
				delete(kv.data, key)
			}
		}
		kv.mu.Unlock()
	}
}

func main() {
	// Create a new KeyValueStore instance
	store := NewKeyValueStore()

	// Start cleanup routine
	go store.CleanupRoutine()

	// Define HTTP handlers
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Handle GET requests (query parameters)
			key := r.URL.Query().Get("key")
			value := r.URL.Query().Get("value")
			if key == "" || value == "" {
				http.Error(w, "Key or value not provided", http.StatusBadRequest)
				return
			}
			store.Set(key, value)
			fmt.Fprintf(w, "Key %s set to value %s\n", key, value)
		} else if r.Method == http.MethodPost {
			// Handle POST requests (JSON body)
			decoder := json.NewDecoder(r.Body)
			var data map[string]interface{}
			if err := decoder.Decode(&data); err != nil {
				http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
				return
			}
			key := data["key"].(string)
			value, ok := data["value"].(string)
			if !ok {
				// If value is not a string, assume it's a JSON object and marshal it
				jsonValue, err := json.Marshal(data["value"])
				if err != nil {
					http.Error(w, "Failed to marshal JSON value", http.StatusInternalServerError)
					return
				}
				value = string(jsonValue)
			}
			if key == "" || value == "" {
				http.Error(w, "Key or value not provided", http.StatusBadRequest)
				return
			}
			store.Set(key, value)
			fmt.Fprintf(w, "Key %s set to value %s\n", key, value)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
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

		var jsonResponse interface{}
		var err error
		// Attempt to parse the entry value as JSON
		err = json.Unmarshal([]byte(entry.Value), &jsonResponse)
		if err != nil {
			// If it's not a valid JSON object, treat it as a string
			jsonResponse = entry.Value
		}

		// Include all fields in the final JSON response
		finalResponse, err := json.Marshal(map[string]interface{}{
			"key":       entry.Key,
			"value":     jsonResponse,
			"createdAt": entry.CreatedAt,
			"updatedAt": entry.UpdatedAt,
		})
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(finalResponse)
	})

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
