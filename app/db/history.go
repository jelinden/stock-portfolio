package db

import (
	"log"

	"github.com/jelinden/stock-portfolio/app/service"
)

func GetHistory(portfolioid string) []service.ClosePrice {
	rows, err := db.Query(`SELECT h.symbol, h.closePrice, h.closePriceDate
		FROM history AS h,
			(SELECT symbol, min(date) as minDate
				FROM portfoliostocks
				WHERE portfolioid LIKE ($1)
				GROUP BY symbol) AS p
		WHERE h.symbol = p.symbol
		AND h.closePriceDate = p.minDate
		ORDER BY p.minDate ASC;`, portfolioid)
	if err != nil {
		log.Printf("failed with '%s'\n", err)
		return []service.ClosePrice{}
	}
	defer rows.Close()
	var closePrices []service.ClosePrice
	for rows.Next() {
		closePrice := service.ClosePrice{}
		err := rows.Scan(
			&closePrice.Symbol,
			&closePrice.ClosePrice,
			&closePrice.ClosePriceDate)
		if err != nil {
			log.Println("scanning row failed", err.Error())
		}
		closePrices = append(closePrices, closePrice)
	}
	return closePrices
}

func SaveHistory(closePrices []service.ClosePrice) {
	for _, c := range closePrices {
		if !isClosePrice(c.Symbol, c.ClosePriceDate) {
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

func isClosePrice(symbol string, date string) bool {
	log.Println(symbol, date)
	row, err := db.Query(`select symbol from history where symbol = $1 and closePriceDate = $2;`, symbol, date)
	if err != nil {
		log.Println(err)
	}
	var s string
	if row.Next() {
		row.Scan(&s)
	}
	log.Println("s1", s)
	row.Close()
	log.Println("s2", s)
	if s != "" {
		return true
	}
	return false
}
