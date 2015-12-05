package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/divideandconquer/go-merge/merge"
	"github.com/hashicorp/consul/api"
)

const divider = "/"

// Loader is a object that can import, initialize, and Get config values
type Loader interface {
	Import(data []byte) error
	Initialize() error
	Get(key string) ([]byte, error)

	// Must functions will panic if they can't do what is requested.
	// They are maingly meant for use with configs that are required for an app to start up
	MustGetString(key string) string
	MustGetBool(key string) bool
	MustGetInt(key string, base int, bitsize int) int64
	MustGetDuration(key string) time.Duration

	//TODO add array support?
}

type cachedLoader struct {
	namespace string
	cacheLock sync.RWMutex
	cache     map[string][]byte
}

// NewCachedLoader creates a Loader that will cache the provided namespace on initialization
// and return data from that cache on Get
func NewCachedLoader(namespace string) Loader {
	return &cachedLoader{namespace: namespace}
}

// Import takes a json byte array and inserts the key value pairs into consul prefixed by the namespace
func (c *cachedLoader) Import(data []byte) error {
	conf := make(map[string]interface{})
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return fmt.Errorf("Unable to parse json data: %v", err)
	}
	kvMap, err := c.compileKeyValues(conf, c.namespace)
	if err != nil {
		return fmt.Errorf("Unable to complie KVs: %v", err)
	}

	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return fmt.Errorf("Could not create consul client: %v", err)
	}
	for k, v := range kvMap {
		p := &api.KVPair{Key: k, Value: v}
		_, err = consul.KV().Put(p, nil)
		if err != nil {
			return fmt.Errorf("Could not write key to consul (%s | %s) %v", k, v, err)
		}
	}
	return nil
}

func (c *cachedLoader) compileKeyValues(data map[string]interface{}, prefix string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for k, v := range data {
		if subMap, ok := v.(map[string]interface{}); ok {
			//recurse and merge results
			compiled, err := c.compileKeyValues(subMap, prefix+divider+k)
			if err != nil {
				return nil, err
			}
			merged := merge.Merge(result, compiled)
			if mm, ok := merged.(map[string][]byte); ok {
				result = mm
			}
		} else {
			//for other types json marshal will turn then into string byte slice for storage
			j, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			result[k] = j
		}
	}
	return result, nil
}

// Initialize loads the consul KV's from the namespace into cache for later retrieval
func (c *cachedLoader) Initialize() error {
	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return fmt.Errorf("Could not create consul client: %v", err)
	}

	pairs, _, err := consul.KV().List(c.namespace, nil)
	if err != nil {
		return fmt.Errorf("Could not pull config from consul: %v", err)
	}

	//write lock the cache incase init is called more than once
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	c.cache = make(map[string][]byte)
	for _, kv := range pairs {
		c.cache[kv.Key] = kv.Value
	}
	return nil
}

// Get fetches the raw config from cache
func (c *cachedLoader) Get(key string) ([]byte, error) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if ret, ok := c.cache[key]; ok {
		return ret, nil
	}
	return nil, fmt.Errorf("Could not find value for key: %s", key)
}

// MustGetString fetches the config and parses it into a string.  Panics on failure.
func (c *cachedLoader) MustGetString(key string) string {
	b, err := c.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}
	return string(b)
}

// MustGetBool fetches the config and parses it into a bool.  Panics on failure.
func (c *cachedLoader) MustGetBool(key string) bool {
	b, err := c.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}
	ret, err := strconv.ParseBool(string(b))
	if err != nil {
		panic(fmt.Sprintf("Could not parse config (%s) into a bool: %v", key, err))
	}
	return ret
}

// MustGetInt fetches the config and parses it into an int.  Panics on failure.
// base and bitsize work the same way as strconv.ParseInt()
func (c *cachedLoader) MustGetInt(key string, base int, bitsize int) int64 {
	b, err := c.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}
	ret, err := strconv.ParseInt(string(b), base, bitsize)
	if err != nil {
		panic(fmt.Sprintf("Could not parse config (%s) into an int: %v", key, err))
	}
	return ret
}

// MustGetDuration fetches the config and parses it into a duration.  Panics on failure.
func (c *cachedLoader) MustGetDuration(key string) time.Duration {
	b, err := c.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Could not fetch config (%s) %v", key, err))
	}
	ret, err := time.ParseDuration(string(b))
	if err != nil {
		panic(fmt.Sprintf("Could not parse config (%s) into a duration: %v", key, err))
	}
	return ret
}
