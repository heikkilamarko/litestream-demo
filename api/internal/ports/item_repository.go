package ports

import (
	"api/internal/domain"
	"context"
)

type ItemRepository interface {
	GetItems(ctx context.Context) ([]*domain.Item, error)
	CreateItem(ctx context.Context, item *domain.Item) error
}
