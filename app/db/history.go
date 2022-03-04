package db

import (
	"log"
	"time"

	"github.com/jelinden/stock-portfolio/app/service"
)

func GetHistory(portfolioid string) []service.ClosePrice {
	timeFrom := time.Now()
	var closePrices = []service.ClosePrice{}

	rows, err := mdb.Query(`SELECT h.symbol, h.closePrice, h.epoch
		FROM history AS h,
			(SELECT symbol, amount, epoch as minDate
				FROM portfoliostocks
				WHERE portfolioid = $1) AS p
		WHERE h.symbol = p.symbol
		AND h.epoch >= p.minDate`, portfolioid)
	if err != nil {
		log.Printf("failed with '%v'\n", err)
		return closePrices
	}
	defer rows.Close()

	for rows.Next() {
		closePrice := service.ClosePrice{}
		err := rows.Scan(
			&closePrice.Symbol,
			&closePrice.ClosePrice,
			&closePrice.Epoch)
		if err != nil {
			log.Println("scanning row failed", err.Error())
		}
		closePrices = append(closePrices, closePrice)
	}
	log.Println("get history took", time.Since(timeFrom))
	return closePrices
}

func SaveHistory(closePrices []service.ClosePrice) {
	for _, c := range closePrices {
		if !isClosePrice(c.Symbol, c.ClosePriceDate) {
			d, err := time.Parse("01/02/2006", c.ClosePriceDate)
			if err != nil {
				log.Println("SaveHistory failed", err)
			} else {
				err := exec(`INSERT INTO history (symbol, closePrice, closePriceDate, epoch) 
				VALUES ($1,$2,$3,$4);`,
					c.Symbol,
					c.ClosePrice,
					c.ClosePriceDate,
					d.Unix()*1000,
				)
				if err != nil {
					log.Printf("failed with '%s' %s\n", err.Error(), c.Symbol)
				}
			}
		}
	}
}

func isClosePrice(symbol string, date string) bool {
	log.Println(symbol, date)
	row, err := db.Query(`select symbol from history where symbol = $1 and closePriceDate = $2;`, symbol, date)
	if err != nil {
		log.Println(err)
		return false
	}
	defer row.Close()
	var s = ""
	if row.Next() {
		row.Scan(&s)
	}
	if s != "" {
		return true
	}
	return false
}
