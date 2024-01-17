package memory

import (
	"context"
	"github.com/JMURv/market/pkg/discovery"
	"sync"
	"time"
)

type serviceName string
type instanceID string

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

func (r *Registry) Register(_ context.Context, id string, name string, port string) error {
	r.Lock()
	defer r.Unlock()

	serviceNameKey := serviceName(name)
	instanceIDKey := instanceID(id)

	if _, ok := r.serviceAddrs[serviceNameKey]; !ok {
		r.serviceAddrs[serviceNameKey] = map[instanceID]*serviceInstance{}
	}

	r.serviceAddrs[serviceNameKey][instanceIDKey] = &serviceInstance{hostPort: port, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(_ context.Context, id string, name string) error {
	r.Lock()
	defer r.Unlock()

	serviceNameKey := serviceName(name)
	instanceIDKey := instanceID(id)

	if _, ok := r.serviceAddrs[serviceNameKey]; !ok {
		return nil
	}

	delete(r.serviceAddrs[serviceNameKey], instanceIDKey)
	return nil
}

func (r *Registry) ReportHealthyState(id string, name string) error {
	r.Lock()
	defer r.Unlock()

	serviceNameKey := serviceName(name)
	instanceIDKey := instanceID(id)

	if _, ok := r.serviceAddrs[serviceNameKey]; !ok {
		return discovery.ServiceIsNotRegistered
	}

	if _, ok := r.serviceAddrs[serviceNameKey][instanceIDKey]; !ok {
		return discovery.ServiceInstanceIsNotRegistered
	}
	r.serviceAddrs[serviceNameKey][instanceIDKey].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of
// active instances of the given service.
func (r *Registry) ServiceAddresses(_ context.Context, name string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	serviceNameKey := serviceName(name)

	if len(r.serviceAddrs[serviceNameKey]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string
	for _, i := range r.serviceAddrs[serviceNameKey] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
