package balancer

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

type mockBalancer struct {
	services map[string]*ServiceLocation
}

func NewMockDNSBalancer(services map[string]*ServiceLocation) DNS {
	return &mockBalancer{services: services}
}

func (m *mockBalancer) FindService(serviceName string) (*ServiceLocation, error) {
	if s, ok := m.services[serviceName]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("Could not find %s", serviceName)
}

func (r *mockBalancer) GetHttpUrl(serviceName string, useTLS bool) (url.URL, error) {
	result := url.URL{}
	loc, err := r.FindService(serviceName)
	if err != nil {
		return result, err
	}
	result.Host = loc.URL
	if loc.Port != 0 {
		result.Host = net.JoinHostPort(loc.URL, strconv.Itoa(loc.Port))
	}
	if useTLS {
		result.Scheme = "https"
	} else {
		result.Scheme = "http"
	}
	return result, nil
}
