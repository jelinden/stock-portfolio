package service

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/jelinden/stock-portfolio/app/util"
)

func GetQuotes(symbols ...string) []Quote {
	quoteData := util.Get(`https://api.iextrading.com/1.0/stock/market/batch?symbols=`+strings.Join(symbols, ",")+`&types=quote`, 5)
	if quoteData == nil {
		return []Quote{}
	}
	return MarshalQuotes(quoteData)
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

type Dividend struct {
	Symbol      string  `json:"symbol"`
	ExDate      int     `json:"exDate"`
	PaymentDate int     `json:"paymentDate"`
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
