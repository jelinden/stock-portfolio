package db

import (
	"log"
	"reflect"

	"github.com/jelinden/stock-portfolio/app/service"
)

func GetDividend(symbols string) []service.Dividend {
	m := getQuery(`select symbol, max(paymentDate) as maxPaymentDate from dividend where symbol in (` + symbols + `) group by symbol order by maxPaymentDate desc`)
	var divs = []service.Dividend{}
	for i := range m {
		if reflect.TypeOf(m[i]["symbol"]).String() == "int64" {
			return []service.Dividend{}
		}
		symbol := m[i]["symbol"].(string)
		maxPaymentDate := m[i]["maxPaymentDate"].(int64)

		dividends := getQuery(`select
			symbol,
			amount,
			type,
			paymentDate,
			exDate
		from dividend
		where paymentDate = $1
		and symbol = $2`, maxPaymentDate, symbol)

		for i := range dividends {
			var div = service.Dividend{}
			if dividends[i]["amount"] != nil {
				div.Amount = dividends[i]["amount"].(float64)
				div.Symbol = dividends[i]["symbol"].(string)
				div.Type = dividends[i]["type"].(string)
				div.PaymentDate = dividends[i]["paymentDate"].(int64)
				div.ExDate = dividends[i]["exDate"].(int64)
			}
			divs = append(divs, div)
		}
	}
	return divs
}

func saveDividends(dividends []service.Dividend) {
	for _, item := range dividends {
		div := getQuery(`select 
							symbol, 
							amount, 
							type, 
							paymentDate, 
							exDate 
						from dividend 
						where symbol = $1 
							and type = $2
							and paymentDate = $3`,
			item.Symbol, item.Type, item.PaymentDate)

		if div == nil {
			log.Println("saving", item.Symbol, item.Amount, item.Type, "to database")
			err := exec("insert into dividend (symbol, amount, type, paymentDate, exDate) values ($1, $2, $3, $4, $5)",
				item.Symbol, item.Amount, item.Type, item.PaymentDate, item.ExDate)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func getDividends() {
	dividends := service.GetDividends(GetPortfolioSymbols()...)
	if len(dividends) > 0 {
		log.Printf("got %v dividends\n", len(dividends))
		saveDividends(dividends)
	}
}
