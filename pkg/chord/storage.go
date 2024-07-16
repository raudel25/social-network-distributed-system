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
	Active  bool   `json:"active"`
}

type Storage interface {
	Get(key string) (Data, error)
	GetAll() (map[string]Data, error)
	GetRemoveAll() (map[string]Data, error)
	Set(key string, value Data) error
	SetAll(dict map[string]Data) error
	Remove(key string, time int64) error
	RemoveAll(dict map[string]int64) error
}

type RamStorage struct {
	store map[string]Data
}

func NewRamStorage() *RamStorage {
	return &RamStorage{store: make(map[string]Data)}
}

func (s *RamStorage) Get(key string) (Data, error) {
	data, ok := s.store[key]
	if !ok || !data.Active {
		return Data{}, nil
	}
	return data, nil
}

func (s *RamStorage) Set(key string, value Data) error {
	value.Active = true
	s.store[key] = value
	return nil
}

func (s *RamStorage) Remove(key string) error {
	value := s.store[key]
	value.Active = false
	s.store[key] = value
	return nil
}

func (s *RamStorage) GetAll() (map[string]Data, error) {
	newDict := make(map[string]Data)

	for key, value := range s.store {
		if value.Active {
			newDict[key] = value
		}
	}

	return newDict, nil
}

// GetRemoveAll(map[string]Data, error) error
func (s *RamStorage) GetRemoveAll() (map[string]Data, error) {
	newDict := make(map[string]Data)

	for key, value := range s.store {
		if !value.Active {
			newDict[key] = value
		}
	}

	return newDict, nil
}

func (s *RamStorage) SetAll(dict map[string]Data) error {
	for key, value := range dict {
		value.Active = true
		s.store[key] = value
	}
	return nil
}

func (s *RamStorage) RemoveAll(dict []string) error {
	for _, key := range dict {
		value := s.store[key]
		value.Active = false
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
	if !value.Active {
		return Data{}, nil
	}
	return value, err
}

func (s *DiskStorage) Set(key string, value Data) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	value.Active = true

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

	value, err := s.Get(key)
	if err != nil {
		return err
	}
	value.Active = false

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
			if value.Active {
				result[relPath] = value
			}
		}
		return nil
	})

	return result, err
}

func (s *DiskStorage) GetRemoveAll() (map[string]Data, error) {
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
			if value.Active {
				result[relPath] = value
			}
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

func (s *DiskStorage) RemoveAll(dict []string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, key := range dict {
		if err := s.Remove(key); err != nil {
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
	data, ok := ds.store[key]
	if !ok || !data.Active {
		return Data{}, nil
	}
	return data, nil
}

// Set sets a value by key
func (ds *DictStorage) Set(key string, value Data) error {
	value.Active = true
	ds.store[key] = value
	return ds.saveToFile()
}

// Remove deletes a value by key
func (ds *DictStorage) Remove(key string, time int64) error {
	value := ds.store[key]
	value.Active = false
	ds.store[key] = value
	return ds.saveToFile()
}

func (ds *DictStorage) GetAll() (map[string]Data, error) {
	newDict := make(map[string]Data)

	for key, value := range ds.store {
		if value.Active {
			newDict[key] = value
		}
	}

	return newDict, nil
}

// GetRemoveAll(map[string]Data, error) error
func (ds *DictStorage) GetRemoveAll() (map[string]Data, error) {
	newDict := make(map[string]Data)

	for key, value := range ds.store {
		if !value.Active {
			newDict[key] = value
		}
	}

	return newDict, nil
}

func (ds *DictStorage) SetAll(dict map[string]Data) error {
	for key, value := range dict {
		value.Active = true
		ds.store[key] = value
	}
	return ds.saveToFile()
}

func (ds *DictStorage) RemoveAll(dict map[string]int64) error {
	for key, time := range dict {
		value := ds.store[key]
		value.Active = false
		value.Version = time
		ds.store[key] = value
	}
	return ds.saveToFile()
}

func (ds *DictStorage) saveToFile() error {
	data, err := json.Marshal(ds.store)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(ds.filename, os.O_WRONLY|os.O_CREATE, fs.FileMode(0644))
	if err != nil {
		return err
	}
	defer file.Close()

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
