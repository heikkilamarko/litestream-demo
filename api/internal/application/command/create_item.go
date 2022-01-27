package command

import (
	"api/internal/domain"
	"api/internal/ports"
	"context"
)

type CreateItem struct {
	Item *domain.Item
}

type CreateItemHandler struct {
	r ports.ItemRepository
}

func NewCreateItemHandler(r ports.ItemRepository) *CreateItemHandler {
	return &CreateItemHandler{r}
}

func (h *CreateItemHandler) Handle(ctx context.Context, c *CreateItem) error {
	return h.r.CreateItem(ctx, c.Item)
}
