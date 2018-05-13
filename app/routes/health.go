package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jelinden/stock-portfolio/app/service"
	"github.com/julienschmidt/httprouter"
)

func Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	j, err := json.Marshal(service.Health)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(j)
}
