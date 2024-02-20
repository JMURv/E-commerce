package memory

import (
	"context"
	"github.com/JMURv/e-commerce/items/internal/repository"
	repo "github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
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

func (r *Repository) ListItem(_ context.Context) ([]*model.Item, error) {
	res := make([]*model.Item, 0, len(r.itemData))
	for _, v := range r.itemData {
		res = append(res, v)
	}
	return res, nil
}

func (r *Repository) GetItemByID(_ context.Context, id uint64) (*model.Item, error) {
	r.RLock()
	i, ok := r.itemData[id]
	r.RUnlock()
	if !ok {
		return nil, repository.ErrNotFound
	}
	return i, nil
}

func (r *Repository) ListUserItemsByID(ctx context.Context, userID uint64) ([]*model.Item, error) {
	return nil, nil
}

func (r *Repository) CreateItem(_ context.Context, i *model.Item) (*model.Item, error) {
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

func (r *Repository) UpdateItem(_ context.Context, itemID uint64, newData *model.Item) (*model.Item, error) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.itemData[itemID]; ok {
		r.itemData[itemID] = newData
		return newData, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) DeleteItem(_ context.Context, itemID uint64) error {
	r.Lock()
	delete(r.categoryData, itemID)
	r.Unlock()
	return nil
}

func (r *Repository) GetAllCategories(_ context.Context) ([]*model.Category, error) {
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

func (r *Repository) CreateCategory(_ context.Context, c *model.Category) (*model.Category, error) {
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

func (r *Repository) UpdateCategory(_ context.Context, categoryID uint64, newData *model.Category) (*model.Category, error) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.categoryData[categoryID]; ok {
		r.categoryData[categoryID] = newData
		return newData, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) DeleteCategory(_ context.Context, categoryID uint64) error {
	r.Lock()
	delete(r.categoryData, categoryID)
	r.Unlock()
	return nil
}

func (r *Repository) ListTags(_ context.Context) ([]*model.Tag, error) {
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

func (r *Repository) CreateTag(_ context.Context, t *model.Tag) (*model.Tag, error) {
	if t.Name == "" {
		return nil, repo.ErrTagNameIsRequired
	}

	r.Lock()
	r.tagsData[t.Name] = t
	r.Unlock()
	return t, nil
}

func (r *Repository) DeleteTag(_ context.Context, tagName string) error {
	r.Lock()
	delete(r.tagsData, tagName)
	r.Unlock()
	return nil
}
