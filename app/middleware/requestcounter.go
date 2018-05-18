package middleware

import (
	"net/http"

	"github.com/jelinden/stock-portfolio/app/service"
	"github.com/julienschmidt/httprouter"
)

func RequestCounter(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		service.Requests++
		fn(w, r, ps)
	}
}
