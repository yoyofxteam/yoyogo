package servicediscovery

import "errors"

type ISelector interface {
	Select(serviceName string) (ServiceInstance, error)
}

type Selector struct {
	DiscoveryCache Cache    //service discovery cache
	Strategy       Strategy //load balancing strategy
}

func NewSelector(discoveryCache Cache, strategy Strategy) *Selector {
	return &Selector{discoveryCache, strategy}
}

// will set strategy and cache options
// Selector( strategy ,  cache ).Select(serviceName).(ServiceInstance)
func (s *Selector) Select(serviceName string) (ServiceInstance, error) {
	service, err := s.DiscoveryCache.GetService(serviceName)
	if err != nil {
		return nil, err
	}
	if len(service.Nodes) == 0 {
		return nil, errors.New("this service don't have any instance")
	}
	return s.Strategy.Next(service.Nodes)
}
