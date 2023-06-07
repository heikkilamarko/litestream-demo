package adapters

import (
	"api/internal/application"
	"api/internal/application/command"
	"api/internal/domain"
	"encoding/json"
	"net/http"

	"github.com/heikkilamarko/goutils"
	"golang.org/x/exp/slog"
)

const (
	errCodeInvalidRequestBody = "invalid_request_body"
)

const (
	fieldRequestBody = "request_body"
)

type HTTPHandlers struct {
	app    *application.Application
	logger *slog.Logger
}

func NewHTTPHandlers(app *application.Application, logger *slog.Logger) *HTTPHandlers {
	return &HTTPHandlers{app, logger}
}

func (h *HTTPHandlers) GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.app.Queries.GetItems.Handle(r.Context())
	if err != nil {
		h.logger.Error(err.Error())
		goutils.WriteInternalError(w, nil)
		return
	}

	goutils.WriteOK(w, items, nil)
}

func (h *HTTPHandlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	c, err := parseCreateItemCommand(r)
	if err != nil {
		h.logger.Error(err.Error())
		goutils.WriteValidationError(w, err)
		return
	}

	if err := h.app.Commands.CreateItem.Handle(r.Context(), c); err != nil {
		h.logger.Error(err.Error())
		goutils.WriteInternalError(w, nil)
		return
	}

	goutils.WriteCreated(w, c.Item, nil)
}

func parseCreateItemCommand(r *http.Request) (*command.CreateItem, error) {
	errs := map[string][]string{}

	item := &domain.Item{}
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		errs[fieldRequestBody] = []string{errCodeInvalidRequestBody}
	}

	if 0 < len(errs) {
		return nil, goutils.ValidationError{Errors: errs}
	}

	return &command.CreateItem{
		Item: item,
	}, nil
}
