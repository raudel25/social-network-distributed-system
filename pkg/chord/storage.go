package chord

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
)

type Data struct {
	Value   string `json:"value"`
	Version int64  `json:"version"`
}

type Storage interface {
	Get(key string) Data
	GetAll() map[string]Data
	Set(key string, value Data)
	SetAll(dict map[string]Data)
	Remove(key string)
}

type RamStorage struct {
	store map[string]Data
}

func NewRamStorage() *RamStorage {
	return &RamStorage{store: make(map[string]Data)}
}

func (s *RamStorage) Get(key string) Data {
	return s.store[key]
}

func (s *RamStorage) Set(key string, value Data) {
	s.store[key] = value
}

func (s *RamStorage) Remove(key string) {
	delete(s.store, key)
}

func (s *RamStorage) GetAll() map[string]Data {
	return s.store
}

func (s *RamStorage) SetAll(dict map[string]Data) {
	for key, value := range dict {
		s.store[key] = value
	}
}

// DictStorage struct
type DictStorage struct {
	store    map[string]Data
	filename string
}

// NewDictStorage creates a new DictStorage
func NewDictStorage(filename string) *DictStorage {
	dict := &DictStorage{store: make(map[string]Data), filename: filename}
	dict.loadFromFile()
	return dict
}

// Get retrieves a value by key
func (ds *DictStorage) Get(key string) Data {
	return ds.store[key]
}

// Set sets a value by key
func (ds *DictStorage) Set(key string, value Data) {
	ds.store[key] = value
}

// Remove deletes a value by key
func (ds *DictStorage) Remove(key string) {
	delete(ds.store, key)
	ds.saveToFile()
}

func (ds *DictStorage) GetAll() map[string]Data {
	return ds.store
}

func (ds *DictStorage) SetAll(dict map[string]Data) {
	for key, value := range dict {
		ds.store[key] = value
	}
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
