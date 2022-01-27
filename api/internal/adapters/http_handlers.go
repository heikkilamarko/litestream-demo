package adapters

import (
	"api/internal/application"
	"api/internal/application/command"
	"api/internal/domain"
	"encoding/json"
	"net/http"

	"github.com/heikkilamarko/goutils"
	"github.com/rs/zerolog"
)

const (
	errCodeInvalidRequestBody = "invalid_request_body"
)

const (
	fieldRequestBody = "request_body"
)

type HTTPHandlers struct {
	app    *application.Application
	logger *zerolog.Logger
}

func NewHTTPHandlers(app *application.Application, logger *zerolog.Logger) *HTTPHandlers {
	return &HTTPHandlers{app, logger}
}

func (h *HTTPHandlers) GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.app.Queries.GetItems.Handle(r.Context())
	if err != nil {
		h.logError(err)
		goutils.WriteInternalError(w, nil)
		return
	}

	goutils.WriteOK(w, items, nil)
}

func (h *HTTPHandlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	c, err := parseCreateItemCommand(r)
	if err != nil {
		h.logError(err)
		goutils.WriteValidationError(w, err)
		return
	}

	if err := h.app.Commands.CreateItem.Handle(r.Context(), c); err != nil {
		h.logError(err)
		goutils.WriteInternalError(w, nil)
		return
	}

	goutils.WriteCreated(w, c.Item, nil)
}

func (h *HTTPHandlers) logError(err error) {
	h.logger.Error().Err(err).Send()
}

func parseCreateItemCommand(r *http.Request) (*command.CreateItem, error) {
	errorMap := map[string]string{}

	item := &domain.Item{}
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		errorMap[fieldRequestBody] = errCodeInvalidRequestBody
	}

	if 0 < len(errorMap) {
		return nil, goutils.NewValidationError(errorMap)
	}

	return &command.CreateItem{
		Item: item,
	}, nil
}
