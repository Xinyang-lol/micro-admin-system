package registry

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/api"
)

type Config struct {
	Address string
}

type Client struct {
	api     *api.Client
	indexes sync.Map
}

func New(address string) (*Client, error) {
	cfg := api.DefaultConfig()
	cfg.Address = address
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{api: client}, nil
}

func (c *Client) Register(ctx context.Context, name string, id string, host string, port int, tags []string) error {
	reg := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Address: host,
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			CheckID:                        "service:" + id,
			TTL:                            "15s",
			DeregisterCriticalServiceAfter: "60s",
		},
	}
	if err := c.api.Agent().ServiceRegister(reg); err != nil {
		return err
	}
	if err := c.api.Agent().UpdateTTL("service:"+id, "alive", api.HealthPassing); err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				_ = c.api.Agent().ServiceDeregister(id)
				return
			case <-ticker.C:
				if err := c.api.Agent().UpdateTTL("service:"+id, "alive", api.HealthPassing); err != nil {
					log.Printf("consul heartbeat failed: %v", err)
				}
			}
		}
	}()
	return nil
}

func (c *Client) NextAddress(service string) (string, error) {
	entries, _, err := c.api.Health().Service(service, "", true, nil)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("service %s has no healthy instance", service)
	}
	var counter *uint64
	value, _ := c.indexes.LoadOrStore(service, new(uint64))
	counter = value.(*uint64)
	idx := atomic.AddUint64(counter, 1)
	entry := entries[int(idx)%len(entries)]
	return fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port), nil
}

func StartRuntimeReporter(ctx context.Context, service string) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				log.Printf("[%s] health reporter tick at %s", service, t.Format(time.RFC3339))
			}
		}
	}()
}
