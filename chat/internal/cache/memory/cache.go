package memory

import (
	"context"
	errs "github.com/JMURv/e-commerce/chat/internal/cache"
	"github.com/JMURv/e-commerce/chat/pkg/model"
	"github.com/opentracing/opentracing-go"
	"sync"
	"time"
)

type Cache struct {
	sync.Mutex
	data map[string]*model.Room
}

func New() *Cache {
	return &Cache{data: make(map[string]*model.Room)}
}

func (c *Cache) Get(ctx context.Context, key string) (*model.Room, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "chat.GetFromCache")
	defer span.Finish()

	c.Lock()
	defer c.Unlock()
	if v, ok := c.data[key]; !ok {
		return nil, errs.ErrNotFoundInCache
	} else {
		return v, nil
	}
}

func (c *Cache) Set(ctx context.Context, t time.Duration, key string, r *model.Room) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "chat.SetToCache")
	defer span.Finish()

	c.Lock()
	defer c.Unlock()
	c.data[key] = r
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "chat.DeleteFromCache")
	defer span.Finish()

	c.Lock()
	defer c.Unlock()
	delete(c.data, key)
	return nil
}
