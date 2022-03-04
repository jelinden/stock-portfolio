package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jelinden/stock-portfolio/app/service"
)

func GetDividend(symbols string) []service.Dividend {

	var divs = []service.Dividend{}

	dividends := getQuery(`select
			symbol,
			amount,
			type,
			paymentDate,
			exDate
		from dividend
		where symbol in (` + symbols + `) order by paymentDate desc limit 50`)

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
	return divs
}

func saveDividends(dividends []service.Dividend) {
	for _, item := range dividends {
		div := getQuery(`select 
							symbol, 
							amount, 
							type, 
							paymentDate, 
							exDate,
							currency
						from dividend 
						where symbol = $1 
							and type = $2
							and paymentDate = $3`,
			item.Symbol, item.Type, item.PaymentDate)

		if div == nil {
			log.Println("saving", item.Symbol, item.Amount, item.Type, item.PaymentDate, item.ExDate, item.Currency, "to database")
			err := exec("insert into dividend (symbol, amount, type, paymentDate, exDate, currency) values ($1, $2, $3, $4, $5, $6)",
				item.Symbol, item.Amount, item.Type, item.PaymentDate, item.ExDate, item.Currency)
			if err != nil {
				log.Println(err)
			}
		} else {
			var amount = item.Amount
			var exDate = item.ExDate
			if i, ok := div[0]["amount"].(float64); ok {
				amount = float64(i)
			}
			if i, ok := div[0]["exDate"].(int64); ok {
				exDate = int64(i)
			}

			if amount != item.Amount ||
				exDate != item.ExDate ||
				fmt.Sprintf("%v", div[0]["currency"]) != item.Currency {
				log.Println("updating", item.Symbol, item.Amount, item.Type, item.PaymentDate, item.ExDate, item.Currency, "to database")
				err := exec(`update dividend 
							set amount=$1, exDate=$2, currency=$3
							where symbol = $4 and type = $5 and paymentDate = $6`,
					item.Amount, item.ExDate, item.Currency, item.Symbol, item.Type, item.PaymentDate)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func getDividends() {
	now := time.Now()
	// get dividends once per week
	if os.Getenv("divs") == "run" ||
		now.Weekday() == 1 ||
		now.Weekday() == 3 { // https://golang.org/pkg/time/#Weekday
		dividends := service.GetDividends(GetPortfolioSymbols()...)
		if len(dividends) > 0 {
			log.Printf("got %v dividends\n", len(dividends))
			saveDividends(dividends)
		} else {
			log.Printf("got zero dividends\n")
		}
	}
}
