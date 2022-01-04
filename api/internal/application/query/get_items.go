package query

import (
	"api/internal/domain"
	"api/internal/ports"
	"context"
)

type GetItemsHandler struct {
	r ports.ItemRepository
}

func NewGetItemsHandler(r ports.ItemRepository) *GetItemsHandler {
	return &GetItemsHandler{r}
}

func (h *GetItemsHandler) Handle(ctx context.Context) ([]*domain.Item, error) {
	return h.r.GetItems(ctx)
}
