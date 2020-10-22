package config

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type mappedLoader struct {
	data     map[string]json.RawMessage
	dataLock sync.RWMutex
}

// NewMappedLoader creates a Loader that will build config from JSON and store it in a a basic map
// This is meant to be a stubbed / example implementation of a loader, its not really actually meant to be used in a production environment.
func NewMappedLoader(data []byte) (Loader, error) {
	ret := &mappedLoader{}
	err := ret.Import(data)
	return ret, err
}

// Import takes a json byte array and inserts the key value pairs into consul prefixed by the namespace
func (m *mappedLoader) Import(data []byte) error {
	conf := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return fmt.Errorf("Unable to parse json data: %v", err)
	}
	m.data = conf
	return nil
}

// Initialize loads the consul KV's from the namespace into cache for later retrieval
func (m *mappedLoader) Initialize() error {
	//noop
	return nil
}

// Get fetches the raw config from cache
func (m *mappedLoader) Get(key string) ([]byte, error) {
	m.dataLock.RLock()
	defer m.dataLock.RUnlock()

	if ret, ok := m.data[key]; ok {
		return ret, nil
	}
	return nil, fmt.Errorf("Could not find value for key: %s", key)
}

// MustGetString fetches the config and parses it into a string.  Panics on failure.
func (m *mappedLoader) MustGetString(key string) string {
	b, err := m.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}

	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal config (%s) %v", key, err))
	}

	return s
}

// MustGetBool fetches the config and parses it into a bool.  Panics on failure.
func (m *mappedLoader) MustGetBool(key string) bool {
	b, err := m.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}
	var ret bool
	err = json.Unmarshal(b, &ret)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal config (%s) %v", key, err))
	}
	return ret
}

// MustGetInt fetches the config and parses it into an int.  Panics on failure.
func (m *mappedLoader) MustGetInt(key string) int {
	b, err := m.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}

	var ret int
	err = json.Unmarshal(b, &ret)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal config (%s) %v", key, err))
	}
	return ret
}

// MustGetDuration fetches the config and parses it into a duration.  Panics on failure.
func (m *mappedLoader) MustGetDuration(key string) time.Duration {
	s := m.MustGetString(key)
	ret, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("Could not parse config (%s) into a duration: %v", key, err))
	}
	return ret
}

func (m *mappedLoader) Put(key string, value []byte) error {
	m.dataLock.Lock()
	defer m.dataLock.Unlock()
	m.data[key] = value
	return nil
}
