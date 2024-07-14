package chord

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type Data struct {
	Value   string `json:"value"`
	Version int64  `json:"version"`
}

type Storage interface {
	Get(key string) (Data, error)
	GetAll() (map[string]Data, error)
	Set(key string, value Data) error
	SetAll(dict map[string]Data) error
	Remove(key string) error
}

type RamStorage struct {
	store map[string]Data
}

func NewRamStorage() *RamStorage {
	return &RamStorage{store: make(map[string]Data)}
}

func (s *RamStorage) Get(key string) (Data, error) {
	return s.store[key], nil
}

func (s *RamStorage) Set(key string, value Data) error {
	s.store[key] = value
	return nil
}

func (s *RamStorage) Remove(key string) error {
	delete(s.store, key)
	return nil
}

func (s *RamStorage) GetAll() (map[string]Data, error) {
	return s.store, nil
}

func (s *RamStorage) SetAll(dict map[string]Data) error {
	for key, value := range dict {
		s.store[key] = value
	}
	return nil
}

type DiskStorage struct {
	basePath string
	mutex    sync.RWMutex
}

func NewDiskStorage(basePath string) (*DiskStorage, error) {
	// Asegurarse de que el directorio base exista
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &DiskStorage{
		basePath: basePath,
	}, nil
}

func (s *DiskStorage) Get(key string) (Data, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	fullPath := filepath.Join(s.basePath, key)
	data, err := os.ReadFile(fullPath)
	if os.IsNotExist(err) {
		return Data{}, nil // Retorna nil si el archivo no existe
	}
	if err != nil {
		return Data{}, err
	}

	var value Data
	err = json.Unmarshal(data, &value)
	return value, err
}

func (s *DiskStorage) Set(key string, value Data) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fullPath := filepath.Join(s.basePath, key)

	// Asegurarse de que el directorio padre exista
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}

func (s *DiskStorage) Remove(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fullPath := filepath.Join(s.basePath, key)
	return os.Remove(fullPath)
}

func (s *DiskStorage) GetAll() (map[string]Data, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make(map[string]Data)
	err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(s.basePath, path)
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var value Data
			if err := json.Unmarshal(data, &value); err != nil {
				return err
			}
			result[relPath] = value
		}
		return nil
	})

	return result, err
}

func (s *DiskStorage) SetAll(dict map[string]Data) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, value := range dict {
		if err := s.Set(key, value); err != nil {
			return err
		}
	}
	return nil
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
func (ds *DictStorage) Get(key string) (Data, error) {
	return ds.store[key], nil
}

// Set sets a value by key
func (ds *DictStorage) Set(key string, value Data) error {
	ds.store[key] = value
	return nil
}

// Remove deletes a value by key
func (ds *DictStorage) Remove(key string) error {
	delete(ds.store, key)
	return ds.saveToFile()
}

func (ds *DictStorage) GetAll() (map[string]Data, error) {
	return ds.store, nil
}

func (ds *DictStorage) SetAll(dict map[string]Data) error {
	for key, value := range dict {
		ds.store[key] = value
	}
	return nil
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
