package adapters

import (
	"net/http"

	"github.com/heikkilamarko/goutils"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	goutils.WriteResponse(w, http.StatusNotFound, nil)
}
