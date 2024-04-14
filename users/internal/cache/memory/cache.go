package memory

import (
	"context"
	errs "github.com/JMURv/e-commerce/users/internal/cache"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"sync"
	"time"
)

type Cache struct {
	sync.Mutex
	data map[string]*model.User
}

func New() *Cache {
	return &Cache{data: make(map[string]*model.User)}
}

func (c *Cache) Get(_ context.Context, key string) (*model.User, error) {
	c.Lock()
	defer c.Unlock()
	if v, ok := c.data[key]; !ok {
		return nil, errs.ErrNotFoundInCache
	} else {
		return v, nil
	}
}

func (c *Cache) Set(_ context.Context, t time.Duration, key string, r *model.User) error {
	c.Lock()
	defer c.Unlock()
	c.data[key] = r
	return nil
}

func (c *Cache) Delete(_ context.Context, key string) error {
	c.Lock()
	defer c.Unlock()
	delete(c.data, key)
	return nil
}
