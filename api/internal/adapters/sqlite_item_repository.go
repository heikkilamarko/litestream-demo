package adapters

import (
	"api/internal/domain"
	"context"
	"database/sql"
)

type SQLiteItemRepository struct {
	db *sql.DB
}

func NewSQLiteItemRepository(db *sql.DB) *SQLiteItemRepository {
	return &SQLiteItemRepository{db}
}

func (r *SQLiteItemRepository) GetItems(ctx context.Context) ([]*domain.Item, error) {
	rows, err := r.db.QueryContext(ctx, "select id, name from items")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []*domain.Item{}

	for rows.Next() {
		item := &domain.Item{}

		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *SQLiteItemRepository) CreateItem(ctx context.Context, item *domain.Item) error {
	return r.db.QueryRowContext(ctx, "insert into items (name) values (?) returning id", item.Name).Scan(&item.ID)
}
