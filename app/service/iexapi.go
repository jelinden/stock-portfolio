package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jelinden/stock-portfolio/app/config"
	"github.com/jelinden/stock-portfolio/app/util"
)

const httpTimeout = 20 // seconds
const iexBaseURL = "https://cloud.iexapis.com/stable/"

func GetQuotes(symbols ...string) []Quote {
	var quoteData = []byte{}
	if len(symbols) > 0 {
		quoteData = util.Get(iexBaseURL+`stock/market/batch?symbols=`+strings.Join(symbols, ",")+`&types=quote&token=`+config.Config.Token, httpTimeout)
		if quoteData == nil {
			return []Quote{}
		}
	}
	return MarshalQuotes(quoteData)
}

func GetDividends(symbols ...string) []Dividend {
	dividends := []Dividend{}
	for _, symbol := range symbols {
		if div := getStockDividends(symbol); div != nil {
			dividends = append(dividends, div...)
		}
		time.Sleep(2 * time.Second)
	}
	return dividends
}

func GetClosePrices(symbols ...string) []ClosePrice {
	closePrices := []ClosePrice{}
	for _, symbol := range symbols {
		if closePrice := getStockClosePrices(symbol); closePrice != nil {
			closePrices = append(closePrices, closePrice...)
		}
		time.Sleep(1200 * time.Millisecond)
	}
	return closePrices
}

// https://cloud.iexapis.com/stable/stock/xom/dividends/3m
func getStockDividends(symbol string) []Dividend {
	rawDividend := []rawDividend{}
	dividend := util.Get(iexBaseURL+`stock/`+symbol+`/dividends/3m?token=`+config.Config.Token, httpTimeout)
	log.Println(string(dividend))
	err := json.Unmarshal(dividend, &rawDividend)
	if err != nil {
		log.Println("Getting dividends for", symbol, "failed", err)
		return nil
	}
	dividends := []Dividend{}
	for i := range rawDividend {
		div := Dividend{Symbol: symbol}
		div.Amount = rawDividend[i].Amount
		div.Type = rawDividend[i].Type
		exDate, err := time.Parse("2006-01-02", rawDividend[i].ExDate)
		if err != nil {
			log.Println(err)
			return nil
		}
		paymentDate, err := time.Parse("2006-01-02", rawDividend[i].PaymentDate)
		if err != nil {
			log.Println("err")
			return nil
		}
		div.ExDate = exDate.Unix() * 1000
		div.PaymentDate = paymentDate.Unix() * 1000
		div.Currency = rawDividend[i].Currency

		if div.Currency == "CAD" {
			cadusd := util.Get(`https://www.bankofcanada.ca/valet/observations/FXCADUSD/json?recent=1`, httpTimeout)
			rate := CurrencyRate{}
			err := json.Unmarshal(cadusd, &rate)
			if err != nil {
				log.Println("err", err)
				return nil
			}
			if f, err := strconv.ParseFloat(rate.Observations[0].FXCADUSD.Value, 64); err == nil {
				fmt.Println("value cadusd", f)
				div.CurrencyRate = f
			}
		}

		dividends = append(dividends, div)
	}
	return dividends
}

// https://cloud.iexapis.com/stable/stock/aapl/chart/5y
func getStockClosePrices(symbol string) []ClosePrice {
	rawClosePrices := []rawClosePrice{}
	fetchedClosePrices := util.Get(iexBaseURL+`stock/`+symbol+`/chart/5y?token=`+config.Config.Token, httpTimeout)
	err := json.Unmarshal(fetchedClosePrices, &rawClosePrices)
	if err != nil {
		log.Println("Getting closePrices for", symbol, "failed")
		return nil
	}
	closePrices := []ClosePrice{}
	for i := range rawClosePrices {
		closePrice := ClosePrice{Symbol: symbol, ClosePrice: rawClosePrices[i].ClosePrice}
		date, err := time.Parse("2006-01-02", rawClosePrices[i].ClosePriceDate)
		if err != nil {
			log.Println("err")
			return nil
		}
		closePrice.ClosePriceDate = date.Format("01/02/2006")
		closePrices = append(closePrices, closePrice)
	}
	return closePrices
}

func MarshalQuotes(q []byte) []Quote {
	var objmap map[string]*json.RawMessage
	var quotes = []Quote{}
	if len(q) > 0 {
		err := json.Unmarshal(q, &objmap)
		if err != nil {
			log.Println("failed to unmarshal", string(q), err)
		}

		for r := range objmap {
			quote := Q{}
			err = json.Unmarshal(*objmap[string(r)], &quote)
			if err != nil {
				log.Println(err.Error(), string(*objmap[string(r)]))
			}
			quotes = append(quotes, quote.Quote)
		}
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
	PreviousClose float64 `json:"previousClose"`
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
	Currency    string  `json:"currency"`
}

type Dividend struct {
	Symbol       string  `json:"symbol"`
	ExDate       int64   `json:"exDate"`
	PaymentDate  int64   `json:"paymentDate"`
	Amount       float64 `json:"amount"`
	Type         string  `json:"type"`
	Currency     string  `json:"currency"`
	CurrencyRate float64 `json:"currencyRate"`
}

type rawClosePrice struct {
	ClosePriceDate string  `json:"date"`
	Epoch          int64   `json:"epoch"`
	ClosePrice     float64 `json:"close"`
}

type ClosePrice struct {
	Symbol         string  `sql:"symbol" json:"symbol"`
	ClosePriceDate string  `sql:"closePriceDate" json:"closePriceDate"`
	Epoch          int64   `sql:"epoch" json:"epoch"`
	ClosePrice     float64 `sql:"closePrice" json:"closePrice"`
}

type CurrencyRate struct {
	Symbol       string `json:"symbol"`
	Observations []Rate `json:"observations"`
}

type Rate struct {
	FXCADUSD FXCADUSD `json:"FXCADUSD"`
}

type FXCADUSD struct {
	Value string `json:"v"`
}
