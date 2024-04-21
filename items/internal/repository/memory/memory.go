package memory

import (
	"context"
	"github.com/JMURv/e-commerce/items/internal/repository"
	repo "github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

type Repository struct {
	sync.RWMutex
	itemData     map[uint64]*model.Item
	categoryData map[uint64]*model.Category
	tagsData     map[string]*model.Tag
}

func New() *Repository {
	return &Repository{
		itemData:     map[uint64]*model.Item{},
		categoryData: map[uint64]*model.Category{},
		tagsData:     map[string]*model.Tag{},
	}
}

func (r *Repository) ListItem(ctx context.Context) ([]*model.Item, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.ListItem.repo")
	defer span.Finish()

	res := make([]*model.Item, 0, len(r.itemData))
	for _, v := range r.itemData {
		res = append(res, v)
	}
	return res, nil
}

func (r *Repository) GetItemByID(ctx context.Context, id uint64) (*model.Item, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.GetItemByID.repo")
	defer span.Finish()

	r.RLock()
	i, ok := r.itemData[id]
	r.RUnlock()
	if !ok {
		return nil, repository.ErrNotFound
	}
	return i, nil
}

func (r *Repository) ListUserItemsByID(ctx context.Context, userID uint64) ([]*model.Item, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.ListUserItemsByID.repo")
	defer span.Finish()

	return nil, nil
}

func (r *Repository) CreateItem(ctx context.Context, i *model.Item) (*model.Item, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.CreateItem.repo")
	defer span.Finish()

	i.ID = uint64(time.Now().Unix())

	if i.UserID == 0 {
		return nil, repo.ErrUserIDRequired
	}
	if i.CategoryID == 0 {
		return nil, repo.ErrCategoryIDRequired
	}

	if i.Name == "" {
		return nil, repo.ErrNameRequired
	}
	if i.Description == "" {
		return nil, repo.ErrDescriptionRequired
	}
	if i.Price == 0 {
		return nil, repo.ErrPriceRequired
	}
	if i.Quantity == 0 {
		i.Quantity = 1
	}

	if _, ok := r.categoryData[i.CategoryID]; !ok {
		return nil, repo.ErrNoSuchCategory
	}

	for idx := range i.Tags {
		currTagName := i.Tags[idx].Name
		if _, ok := r.tagsData[currTagName]; !ok {
			newTag := &model.Tag{Name: currTagName}
			r.tagsData[currTagName] = newTag
			i.Tags[idx] = newTag
		}
	}

	i.Status = "created"
	i.CreatedAt = *timestamppb.Now()

	r.Lock()
	r.itemData[i.ID] = i
	r.Unlock()
	return i, nil
}

func (r *Repository) UpdateItem(ctx context.Context, itemID uint64, newData *model.Item) (*model.Item, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.UpdateItem.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()
	if _, ok := r.itemData[itemID]; ok {
		r.itemData[itemID] = newData
		return newData, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) DeleteItem(ctx context.Context, itemID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.DeleteItem.repo")
	defer span.Finish()

	r.Lock()
	delete(r.categoryData, itemID)
	r.Unlock()
	return nil
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.GetAllCategories.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()

	res := make([]*model.Category, 0, len(r.categoryData))
	for i := range r.categoryData {
		if v, ok := r.categoryData[i]; ok {
			res = append(res, v)
		}
	}
	return res, nil
}

func (r *Repository) CreateCategory(ctx context.Context, c *model.Category) (*model.Category, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.CreateCategory.repo")
	defer span.Finish()

	c.ID = uint64(time.Now().Unix())

	if c.Name == "" {
		return nil, repo.ErrCategoryNameIsRequired
	}

	if c.Description == "" {
		return nil, repo.ErrCategoryDescriptionIsRequired
	}

	r.Lock()
	r.categoryData[c.ID] = c
	r.Unlock()
	return c, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, categoryID uint64, newData *model.Category) (*model.Category, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.UpdateCategory.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()
	if _, ok := r.categoryData[categoryID]; ok {
		r.categoryData[categoryID] = newData
		return newData, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) DeleteCategory(ctx context.Context, categoryID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.DeleteCategory.repo")
	defer span.Finish()

	r.Lock()
	delete(r.categoryData, categoryID)
	r.Unlock()
	return nil
}

func (r *Repository) ListTags(ctx context.Context) ([]*model.Tag, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.ListTags.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()

	res := make([]*model.Tag, 0, len(r.tagsData))
	for i := range r.tagsData {
		if v, ok := r.tagsData[i]; ok {
			res = append(res, v)
		}
	}
	return res, nil
}

func (r *Repository) CreateTag(ctx context.Context, t *model.Tag) (*model.Tag, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.CreateTag.repo")
	defer span.Finish()

	if t.Name == "" {
		return nil, repo.ErrTagNameIsRequired
	}

	r.Lock()
	r.tagsData[t.Name] = t
	r.Unlock()
	return t, nil
}

func (r *Repository) DeleteTag(ctx context.Context, tagName string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.DeleteTag.repo")
	defer span.Finish()

	r.Lock()
	delete(r.tagsData, tagName)
	r.Unlock()
	return nil
}
