package client

import (
	"fmt"
	"log"
	"time"
)

type mockLoader struct {
	data map[string]interface{}
}

func NewMockLoader(data map[string]interface{}) Loader {
	return &mockLoader{data}
}
func (m *mockLoader) Import(data []byte) error {
	return nil
}
func (m *mockLoader) Initialize() error {
	return nil
}
func (m *mockLoader) Get(key string) ([]byte, error) {
	if ret, ok := m.data[key]; ok {
		if result, ok := ret.([]byte); ok {
			return result, nil
		}
	}
	return nil, fmt.Errorf("Key (%s) not set in mock.", key)
}

func (m *mockLoader) MustGetString(key string) string {
	if ret, ok := m.data[key]; ok {
		if result, ok := ret.(string); ok {
			return result
		}
	}
	log.Fatalf("Key (%s) not set in mock", key)
	return ""
}

func (m *mockLoader) MustGetBool(key string) bool {
	if ret, ok := m.data[key]; ok {
		if result, ok := ret.(bool); ok {
			return result
		}
	}
	log.Fatalf("Key (%s) not set in mock", key)
	return false
}

func (m *mockLoader) MustGetInt(key string) int {
	if ret, ok := m.data[key]; ok {
		if result, ok := ret.(int); ok {
			return result
		}
	}
	log.Fatalf("Key (%s) not set in mock", key)
	return 0
}

func (m *mockLoader) MustGetDuration(key string) time.Duration {
	if ret, ok := m.data[key]; ok {
		if result, ok := ret.(time.Duration); ok {
			return result
		}
	}
	log.Fatalf("Key (%s) not set in mock", key)
	return 0
}
