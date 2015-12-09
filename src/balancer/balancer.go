package balancer

// DNS balancer finds services through dns and balances load across them
type DNS interface {
	FindService(serviceName string) (ServiceLocation, error)
}

// ServiceLocation is a represensation of where a service lives
type ServiceLocation struct {
	URL  string
	Port int
}
