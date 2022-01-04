package application

import (
	"api/internal/application/command"
	"api/internal/application/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateItem *command.CreateItemHandler
}

type Queries struct {
	GetItems *query.GetItemsHandler
}
