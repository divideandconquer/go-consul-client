package balancer

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

type cachedServiceLocation struct {
	Services []*ServiceLocation
	CachedAt time.Time
}

type randomBalancer struct {
	environment   string
	consulCatalog *api.Health
	cache         map[string]cachedServiceLocation
	cacheLock     sync.RWMutex //TODO lock per serviceName
	ttl           time.Duration
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewRandomDNSBalancer(environment string, consulAddr string, cacheTTL time.Duration) (DNS, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr
	consul, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to consul: %v", err)
	}

	r := randomBalancer{}
	r.cache = make(map[string]cachedServiceLocation)
	r.environment = environment
	r.ttl = cacheTTL
	r.consulCatalog = consul.Health()
	return &r, nil
}

func (r *randomBalancer) FindService(serviceName string) (*ServiceLocation, error) {
	services, err := r.getServiceFromCache(serviceName)
	if err != nil || len(services) == 0 {
		services, err = r.writeServiceToCache(serviceName)
		if err != nil {
			return nil, err
		}
	}
	return r.pickService(services), nil
}

func (r *randomBalancer) pickService(services []*ServiceLocation) *ServiceLocation {
	return services[rand.Intn(len(services))]
}

func (r *randomBalancer) getServiceFromCache(serviceName string) ([]*ServiceLocation, error) {
	r.cacheLock.RLock()
	defer r.cacheLock.RUnlock()

	if result, ok := r.cache[serviceName]; ok {
		if time.Now().UTC().Before(result.CachedAt.Add(r.ttl)) {
			return result.Services, nil
		}
		return nil, fmt.Errorf("Cache for %s is expired", serviceName)
	}
	return nil, fmt.Errorf("Could not find %s in cache", serviceName)
}

// writeServiceToCache locks specifically to alleviate load on consul some additional lock time
// is preferable to extra consul calls
func (r *randomBalancer) writeServiceToCache(serviceName string) ([]*ServiceLocation, error) {
	//acquire a write lock
	r.cacheLock.Lock()
	defer r.cacheLock.Unlock()

	//check the cache again in case we've fetched since the last check
	//(our lock could have been waiting for another call to this function)
	if result, ok := r.cache[serviceName]; ok {
		if time.Now().UTC().Before(result.CachedAt.Add(r.ttl)) {
			return result.Services, nil
		}
	}

	//it still isn't in the cache, lets put it there
	consulServices, _, err := r.consulCatalog.Service(serviceName, r.environment, true, nil)
	if err != nil {
		return nil, fmt.Errorf("Error reaching consul for service lookup %v", err)
	}

	if len(consulServices) == 0 {
		return nil, fmt.Errorf("No services found for %s", serviceName)
	}

	//setup service locations
	var services []*ServiceLocation
	for _, v := range consulServices {
		s := &ServiceLocation{}
		s.URL = v.Service.Address
		s.Port = v.Service.Port
		services = append(services, s)
	}

	// cache
	c := cachedServiceLocation{Services: services, CachedAt: time.Now().UTC()}
	r.cache[serviceName] = c
	return services, nil
}
