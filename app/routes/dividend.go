package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jelinden/stock-portfolio/app/db"
	"github.com/julienschmidt/httprouter"
)

func GetDividend(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var marshalled = []byte(`{"response": "failed"}`)
	var err error
	user := getUser(r)
	if user.ID != "" && verifySymbol(r.URL.Query().Get("symbols")) {
		s := strings.Split(r.URL.Query().Get("symbols"), ",")

		var symbols string
		for i, item := range s {
			if i == 0 {
				symbols = `"` + item + `"`
			} else {
				symbols = symbols + `,"` + item + `"`
			}
		}

		dividends := db.GetDividend(symbols)
		marshalled, err = json.Marshal(dividends)
		if err != nil {
			log.Println(err)
		}
	}
	ok(w, marshalled)
}
