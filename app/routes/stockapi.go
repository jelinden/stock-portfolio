package routes

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jelinden/stock-portfolio/app/db"
	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/jelinden/stock-portfolio/app/util"
	"github.com/julienschmidt/httprouter"
)

func AddStock(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := getUser(r)
	if user.ID != "" {
		symbol := strings.ToUpper(r.FormValue("symbol"))
		portfolioid := r.FormValue("portfolioid")
		amount := r.FormValue("amount")
		price := r.FormValue("price")
		date := r.FormValue("date")
		commission := r.FormValue("commission")
		if verifySymbol(symbol) && isPortfolioOwner(user.ID, portfolioid) {
			stock := domain.PortfolioStock{Portfolioid: portfolioid, Symbol: symbol}

			a, err := checkAmount(amount)
			if err {
				addStockResponse(w, portfolioid+"?amountMsg=true")
				return
			}
			p, err := checkPrice(price)
			if err {
				addStockResponse(w, portfolioid+"?priceMsg=true")
				return
			}
			c, err := checkCommission(commission)
			if err {
				addStockResponse(w, portfolioid+"?commissionMsg=true")
				return
			}
			d, err := checkDate(date)
			if err {
				addStockResponse(w, portfolioid+"?dateMsg=true")
				return
			}
			stock.Date = *d
			stock.Amount = *a
			stock.Price = *p
			stock.Commission = *c
			db.AddStock(stock)
			addStockResponse(w, portfolioid)
			return
		}
		addStockResponse(w, portfolioid+"?symbolMsg=true")
		return
	}
	w.Header().Add("Location", "/login")
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(302)
}

func RemoveStock(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := getUser(r)
	if user.ID != "" {
		symbol := p.ByName("symbol")
		portfolioid := p.ByName("portfolioid")
		if verifySymbol(symbol) && isPortfolioOwner(user.ID, portfolioid) {
			err := db.RemoveStock(portfolioid, symbol)
			if err != nil {
				log.Println(err.Error())
				ok(w, []byte(`{"error": "`+err.Error()+`"}`))
				return
			}
			ok(w, []byte(`{"removed": "`+p.ByName("symbol")+`"}`))
			return
		}
	}
	ok(w, []byte(`{"error": "Not logged in"}`))
}

func isPortfolioOwner(userID, portfolioID string) bool {
	if util.IsID(portfolioID) {
		p := db.GetPortfolios(userID)
		for _, item := range p.Portfolios {
			if item.Portfolioid == portfolioID {
				return true
			}
		}
	}
	return false
}

func addStockResponse(w http.ResponseWriter, path string) {
	w.Header().Add("Location", "/portfolio/"+path)
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(302)
}

func ok(w http.ResponseWriter, content []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.Write(content)
}

func verifySymbol(s string) bool {
	re, err := regexp.Compile(`^[a-zA-ZöäåÖÄÅ, ]+$`)
	if err != nil {
		log.Println(err.Error())
	}
	return re.MatchString(s)
}

func checkAmount(amount string) (*int, bool) {
	a, err := strconv.Atoi(amount)
	if err != nil {
		return nil, true
	}
	return &a, false
}

func checkPrice(price string) (*float64, bool) {
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return nil, true
	}
	return &p, false
}

func checkCommission(commission string) (*float64, bool) {
	if commission == "" {
		com := 0.0
		return &com, false
	}
	c, err := strconv.ParseFloat(commission, 64)
	if err != nil {
		return nil, true
	}
	return &c, false
}

func checkDate(date string) (*string, bool) {
	_, err := time.Parse("01/02/2006", date)
	if err != nil {
		return nil, true
	}
	return &date, false
}
