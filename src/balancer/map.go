package balancer

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

type mapBalancer struct {
	services map[string]string
}

// NewMapBalancer will return the URL from the configured map name -> URL:port
func NewMapBalancer(services map[string]string) DNS {
	return &mapBalancer{services: services}
}

func (m *mapBalancer) FindService(serviceName string) (*ServiceLocation, error) {
	s, ok := m.services[serviceName]
	if !ok {
		return nil, fmt.Errorf("Could not find %s", serviceName)
	}

	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return nil, err
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	return &ServiceLocation{URL: host, Port: p}, nil
}

func (m *mapBalancer) GetHttpUrl(serviceName string, useTLS bool) (url.URL, error) {
	result := url.URL{}
	loc, err := m.FindService(serviceName)
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
