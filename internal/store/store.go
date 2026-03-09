package store

import (
	"encoding/json"
	"errors"
	"keep/internal/keep"
	"os"
	"sync"
)

const StoreFile = "store/keep.json"

var mu sync.Mutex

// loadKeeps reads the keeps from the JSON file
func loadKeeps() ([]keep.Keep, error) {
	file, err := os.Open(StoreFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []keep.Keep{}, nil // treat as empty list if file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	var keeps []keep.Keep
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&keeps); err != nil && err.Error() != "EOF" {
		return nil, err
	}
	return keeps, nil
}

// saveKeeps writes the keeps to the JSON file
func saveKeeps(keeps []keep.Keep) error {
	if err := os.MkdirAll("store", os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(StoreFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(keeps)
}

// List returns all keeps.
func List() ([]keep.Keep, error) {
	mu.Lock()
	defer mu.Unlock()
	return loadKeeps()
}

// Create adds a new keep.
func Create(newKeep keep.Keep) (*keep.Keep, error) {
	mu.Lock()
	defer mu.Unlock()

	keeps, err := loadKeeps()
	if err != nil {
		return nil, err
	}

	keeps = append(keeps, newKeep)

	if err := saveKeeps(keeps); err != nil {
		return nil, err
	}
	return &newKeep, nil
}

// Update updates the keep with the given UUID.
func Update(uuid keep.KeepUUID, updatedKeep string) error {
	mu.Lock()
	defer mu.Unlock()

	keeps, err := loadKeeps()
	if err != nil {
		return err
	}

	found := false
	for i, k := range keeps {
		if k.UUID == uuid {
			keeps[i].Content = updatedKeep
			found = true
			break
		}
	}
	if !found {
		return errors.New("keep not found")
	}
	return saveKeeps(keeps)
}

// Delete removes a keep with the given UUID.
func Delete(uuid keep.KeepUUID) error {
	mu.Lock()
	defer mu.Unlock()

	keeps, err := loadKeeps()
	if err != nil {
		return err
	}

	found := false
	newKeeps := make([]keep.Keep, 0, len(keeps))
	for _, k := range keeps {
		if k.UUID == uuid {
			found = true
			continue
		}
		newKeeps = append(newKeeps, k)
	}
	if !found {
		return errors.New("keep not found")
	}
	return saveKeeps(newKeeps)
}
