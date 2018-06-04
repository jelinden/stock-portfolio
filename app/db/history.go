package db

import (
	"log"

	"github.com/jelinden/stock-portfolio/app/service"
)

func isClosePrice(symbol string, date int64) bool {
	row, err := db.Query(`select symbol from history where symbol = $1 and closePriceDate = $2;`, symbol, date)
	if err != nil {
		log.Println(err)
	}
	var s string
	if row.Next() {
		row.Scan(&s)
	}
	row.Close()
	if s != "" {
		return true
	}
	return false
}

func SaveHistory(closePrices []service.ClosePrice) {
	for _, c := range closePrices {
		if isClosePrice(c.Symbol, c.ClosePriceDate) {
			err := exec(`UPDATE history SET 
					symbol = $1,
					closePrice = $2,
					closePriceDate = $3;`,
				c.Symbol,
				c.ClosePrice,
				c.ClosePriceDate)
			if err != nil {
				log.Printf("failed with '%s' %s\n", err.Error(), c.Symbol)
			}
		} else {
			err := exec(`INSERT INTO history (symbol, closePrice, closePriceDate) 
				VALUES ($1,$2,$3);`,
				c.Symbol,
				c.ClosePrice,
				c.ClosePriceDate)
			if err != nil {
				log.Printf("failed with '%s' %s\n", err.Error(), c.Symbol)
			}
		}
	}
}
