package chord

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
)

type Storage interface {
	Get(key string) string
	Set(key string, value string)
	Remove(key string)
}

type RamStorage struct {
	store map[string]string
}

func NewRamStorage() *RamStorage {
	return &RamStorage{store: make(map[string]string)}
}

func (s *RamStorage) Get(key string) string {
	return s.store[key]
}

func (s *RamStorage) Set(key string, value string) {
	s.store[key] = value
}

func (s *RamStorage) Remove(key string) {
	delete(s.store, key)
}

// DictStorage struct
type DictStorage struct {
	store    map[string]string
	filename string
}

// NewDictStorage creates a new DictStorage
func NewDictStorage(filename string) *DictStorage {
	dict := &DictStorage{store: make(map[string]string), filename: filename}
	dict.loadFromFile()
	return dict
}

// Get retrieves a value by key
func (ds *DictStorage) Get(key string) string {
	return ds.store[key]
}

// Set sets a value by key
func (ds *DictStorage) Set(key string, value string) {
	ds.store[key] = value
}

// Remove deletes a value by key
func (ds *DictStorage) Remove(key string) {
	delete(ds.store, key)
	ds.saveToFile()
}

// SaveToFile saves the storage to a JSON file
func (ds *DictStorage) saveToFile() error {
	data, err := json.Marshal(ds.store)
	if err != nil {
		return err
	}
	return os.WriteFile(ds.filename, data, fs.FileMode(0644))
}

// LoadFromFile loads the storage from a JSON file
func (ds *DictStorage) loadFromFile() error {
	file, err := os.Open(ds.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(io.Reader(file))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &ds.store)
}
