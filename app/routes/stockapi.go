package routes

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jelinden/stock-portfolio/app/db"
	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/julienschmidt/httprouter"
)

func AddStock(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := getUser(r)
	if user.ID != "" {
		symbol := r.FormValue("symbol")
		portfolioid := r.FormValue("portfolioid")
		amount := r.FormValue("amount")
		price := r.FormValue("price")
		date := r.FormValue("date")
		commission := r.FormValue("commission")
		if verifyString(symbol) {
			stock := domain.PortfolioStock{}
			stock.Portfolioid = portfolioid
			stock.Symbol = symbol
			var err error
			stock.Amount, err = strconv.Atoi(amount)
			p, err := strconv.ParseFloat(price, 64)
			c, err := strconv.ParseFloat(commission, 64)
			if err != nil {
				log.Println("converting stock items failed", err.Error())
				ok(w, []byte(`{"error": "`+err.Error()+`"}`))
				return
			}
			stock.Date = date
			stock.Price = p
			stock.Commission = c
			db.AddStock(stock)
			w.Header().Add("Location", "/portfolio/"+portfolioid)
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(302)
			return
		}
	}
	w.Header().Add("Location", "/")
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(302)
}

func RemoveStock(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := getUser(r)
	if user.ID != "" {
		symbol := p.ByName("symbol")
		portfolioid := p.ByName("portfolioid")
		if verifyString(symbol) {
			err := db.RemoveStock(portfolioid, symbol)
			if err != nil {
				log.Println(err.Error())
				ok(w, []byte(`{"error": "`+err.Error()+`"}`))
			}
			ok(w, []byte(`{"removed": "`+p.ByName("symbol")+`"}`))
			return
		}
	}
	ok(w, []byte(`{"error": "Not logged in"}`))
}

func ok(w http.ResponseWriter, content []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.Write(content)
}

func verifyString(v string) bool {
	return len(v) > 0
}

func verifySymbol(s string) bool {
	re, err := regexp.Compile(`^[a-zA-ZöäåÖÄÅ ]+$`)
	if err != nil {
		log.Println(err.Error())
	}
	return re.MatchString(s)
}
