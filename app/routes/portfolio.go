package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jelinden/stock-portfolio/app/db"
	"github.com/jelinden/stock-portfolio/app/util"
	"github.com/julienschmidt/httprouter"
)

func AddPortfolio(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := getUser(r)
	if user.ID != "" {
		name := r.FormValue("name")
		if verifyString(name) {
			portfolioid := util.GetID()
			db.AddPortfolio(portfolioid, user.ID, name)
		}
	}
	w.Header().Add("Location", "/")
	w.WriteHeader(302)
	w.Write(nil)
}

func GetPortfolios(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var marshalled = []byte(`{"response": "failed"}`)
	var err error
	user := getUser(r)
	if user.ID != "" {
		portfolios := db.GetPortfolios(user.ID)
		marshalled, err = json.Marshal(portfolios)
		if err != nil {
			log.Println(err)
		}
	}
	ok(w, marshalled)
}

func GetPortfolio(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var marshalled = []byte(`{"response": "failed"}`)
	var err error
	user := getUser(r)
	if user.ID != "" {
		portfolio := db.GetPortfolio(p.ByName("id"))
		marshalled, err = json.Marshal(portfolio)
		if err != nil {
			log.Println(err)
		}
	}
	ok(w, marshalled)
}
