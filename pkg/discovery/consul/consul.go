package consul

import (
	"context"
	"fmt"
	"github.com/JMURv/e-commerce/pkg/discovery"
	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

func (r *Registry) Register(_ context.Context, id string, name string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of <host>:<port>, example: localhost:8081")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(
		&consul.AgentServiceRegistration{
			Address: parts[0],
			ID:      id,
			Name:    name,
			Port:    port,
			Check: &consul.AgentServiceCheck{
				CheckID: id,
				TTL:     "5s",
			},
		},
	)
}

func (r *Registry) Deregister(_ context.Context, id string, _ string) error {
	return r.client.Agent().ServiceDeregister(id)
}

// ServiceAddresses returns the list of addresses of
// active instances of the given service.
func (r *Registry) ServiceAddresses(_ context.Context, name string) ([]string, error) {
	entries, _, err := r.client.Health().Service(name, "", true, nil)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}

	return res, nil
}

// ReportHealthyState is a push mechanism for
// reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(id string, _ string) error {
	return r.client.Agent().PassTTL(id, "")
}
