package service

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/jelinden/stock-portfolio/app/util"
)

func GetQuotes(symbols ...string) []Quote {
	quoteData := util.Get(`https://api.iextrading.com/1.0/stock/market/batch?symbols=`+strings.Join(symbols, ",")+`&types=quote`, 8)
	if quoteData == nil {
		return []Quote{}
	}
	return MarshalQuotes(quoteData)
}

func GetDividends(symbols ...string) []Dividend {
	dividends := []Dividend{}
	for _, symbol := range symbols {
		if div := getStockDividends(symbol); div != nil {
			dividends = append(dividends, div...)
		}
		time.Sleep(1200 * time.Millisecond)
	}
	return dividends
}

func getStockDividends(symbol string) []Dividend {
	rawDividend := []rawDividend{}
	dividend := util.Get(`https://api.iextrading.com/1.0/stock/`+symbol+`/dividends/5y`, 8)
	err := json.Unmarshal(dividend, &rawDividend)
	if err != nil {
		log.Println("Getting dividends for", symbol, "failed")
		return nil
	}
	dividends := []Dividend{}
	for i := range rawDividend {
		div := Dividend{Symbol: symbol}
		div.Amount = rawDividend[i].Amount
		div.Type = rawDividend[i].Type
		exDate, err := time.Parse("2006-01-02", rawDividend[i].ExDate)
		if err != nil {
			log.Println("err")
			return nil
		}
		paymentDate, err := time.Parse("2006-01-02", rawDividend[i].PaymentDate)
		if err != nil {
			log.Println("err")
			return nil
		}
		div.ExDate = exDate.Unix() * 1000
		div.PaymentDate = paymentDate.Unix() * 1000
		dividends = append(dividends, div)
	}
	return dividends
}

func MarshalQuotes(q []byte) []Quote {
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(q, &objmap)
	if err != nil {
		log.Println(err)
	}
	quotes := []Quote{}
	for r := range objmap {
		quote := Q{}
		err = json.Unmarshal(*objmap[string(r)], &quote)
		if err != nil {
			log.Println(err.Error())
		}
		quotes = append(quotes, quote.Quote)
	}
	return quotes
}

type Q struct {
	Quote Quote `json:"quote"`
}

type Quote struct {
	Symbol        string  `json:"symbol"`
	CompanyName   string  `json:"companyName"`
	Sector        string  `json:"sector"`
	Close         float64 `json:"close"`
	CloseTime     int     `json:"closeTime"`
	LatestPrice   float64 `json:"latestPrice"`
	LatestUpdate  int     `json:"latestUpdate"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	PERatio       float64 `json:"peRatio"`
}

type rawDividend struct {
	Symbol      string  `json:"symbol"`
	ExDate      string  `json:"exDate"`
	PaymentDate string  `json:"paymentDate"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
}

type Dividend struct {
	Symbol      string  `json:"symbol"`
	ExDate      int64   `json:"exDate"`
	PaymentDate int64   `json:"paymentDate"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
}

// https://api.iextrading.com/1.0/stock/OHI/dividends/5y
/*
[
	{
		exDate: "2018-01-30",
		paymentDate: "2018-02-15",
		recordDate: "2018-01-31",
		declaredDate: "2018-01-16",
		amount: 0.66,
		flag: "FI",
		type: "Dividend income",
		qualified: "",
		indicated: ""
	},
*/
// https://api.iextrading.com/1.0/stock/market/batch?symbols=bns,xom,adp,ohi&types=quote
