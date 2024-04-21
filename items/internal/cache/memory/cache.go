package memory

import (
	"context"
	errs "github.com/JMURv/e-commerce/items/internal/cache"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/opentracing/opentracing-go"
	"sync"
	"time"
)

type Cache struct {
	sync.Mutex
	data map[string]*model.Item
}

func New() *Cache {
	return &Cache{data: make(map[string]*model.Item)}
}

func (c *Cache) Get(ctx context.Context, key string) (*model.Item, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.GetFromCache")
	defer span.Finish()

	c.Lock()
	defer c.Unlock()
	if v, ok := c.data[key]; !ok {
		return nil, errs.ErrNotFoundInCache
	} else {
		return v, nil
	}
}

func (c *Cache) Set(ctx context.Context, t time.Duration, key string, r *model.Item) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.SetToCache")
	defer span.Finish()

	c.Lock()
	defer c.Unlock()
	c.data[key] = r
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.DeleteFromCache")
	defer span.Finish()

	c.Lock()
	defer c.Unlock()
	delete(c.data, key)
	return nil
}
