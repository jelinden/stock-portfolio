package db

import (
	"log"

	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/jelinden/stock-portfolio/app/service"
)

const symbolQuery = `SELECT distinct symbol FROM portfoliostocks;`

func AddStock(stock domain.PortfolioStock) bool {
	err := exec(`INSERT INTO portfoliostocks (portfolioid, symbol, amount, price, date, commission) VALUES ($1,$2,$3,$4,$5,$6);`,
		stock.Portfolioid,
		stock.Symbol,
		stock.Amount,
		stock.Price,
		stock.Date,
		stock.Commission)
	if err != nil {
		log.Printf("failed with '%s'\n", err.Error())
		return false
	}
	return true
}

func RemoveStock(portfolioid, symbol string) error {
	err := exec(`DELETE FROM portfoliostocks WHERE portfolioid = $1 AND symbol = $2;`,
		portfolioid,
		symbol)
	if err != nil {
		log.Printf("failed with '%s'\n", err.Error())
		return err
	}
	return nil
}

func isQuote(symbol string) bool {
	row, err := db.Query(`select symbol from quotes where symbol = $1;`, symbol)
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

func SaveQuotes(quotes []service.Quote) {
	for _, q := range quotes {
		if isQuote(q.Symbol) {
			err := exec(`UPDATE quotes SET 
					latestPrice = $1,
					latestUpdate = $2,
					close = $3, 
					closeTime = $4, 
					change = $5, 
					changePercent = $6,
					PERatio = $7
					WHERE symbol = $8;`,
				q.LatestPrice,
				q.LatestUpdate,
				q.Close,
				q.CloseTime,
				q.Change,
				q.ChangePercent,
				q.PERatio,
				q.Symbol)
			if err != nil {
				log.Printf("failed with '%s' %s\n", err.Error(), q.Symbol)
			}
		} else {
			err := exec(`INSERT INTO quotes (symbol,companyName,sector,latestPrice,latestUpdate,close,closeTime,change,changePercent,PERatio) 
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`,
				q.Symbol,
				q.CompanyName,
				q.Sector,
				q.LatestPrice,
				q.LatestUpdate,
				q.Close,
				q.CloseTime,
				q.Change,
				q.ChangePercent,
				q.PERatio)
			if err != nil {
				log.Printf("failed with '%s' %s\n", err.Error(), q.Symbol)
			}
		}
	}
}

func GetPortfolioSymbols() []string {
	return queryPortfolioSymbols()
}

func AddPortfolio(portfolioid, userid, name string) bool {
	err := exec(`INSERT INTO portfolio (portfolioid, userid, name) VALUES ($1,$2,$3);`,
		portfolioid,
		userid,
		name)
	if err != nil {
		log.Printf("failed with '%s'\n", err.Error())
		return false
	}
	return true
}

func GetPortfolios(userid string) domain.Portfolios {
	rows, err := db.Query(`select portfolioid, userid, name from portfolio where userid like ($1);`, userid)
	if err != nil {
		log.Printf("failed with '%s'\n", err.Error())
	}
	defer rows.Close()
	var portfolios domain.Portfolios
	for rows.Next() {
		portfolio := domain.Portfolio{}
		rows.Scan(&portfolio.Portfolioid, &portfolio.Userid, &portfolio.Name)
		portfolios.Portfolios = append(portfolios.Portfolios, portfolio)
	}
	return portfolios
}

func GetPortfolio(portfolioid string) domain.PortfolioStocks {
	rows, err := db.Query(`SELECT 
			p.portfolioid,
			po.name,
			quotes.companyName,
			p.symbol,
			p.price,
			p.commission,
			p.amount,
			quotes.latestPrice,
			quotes.latestUpdate,
			quotes.close, 
			quotes.closeTime, 
			quotes.PERatio,
			quotes.change,
			quotes.changePercent
		FROM portfolio AS po,
			(SELECT portfolioid,
				symbol,
				sum(price*float64(amount)) as price,
				sum(commission) as commission,
				sum(amount) as amount
				FROM portfoliostocks
				GROUP BY symbol, portfolioid) AS p
			LEFT JOIN quotes on p.symbol = quotes.symbol 
		WHERE p.portfolioid LIKE ($1)
		AND po.portfolioid = p.portfolioid
		ORDER BY p.symbol ASC;`, portfolioid)
	if err != nil {
		log.Printf("failed with '%s'\n", err)
		return domain.PortfolioStocks{}
	}
	defer rows.Close()
	var portfolioStocks domain.PortfolioStocks
	for rows.Next() {
		stock := domain.PortfolioStock{}
		err := rows.Scan(&stock.Portfolioid,
			&portfolioStocks.PortfolioName,
			&stock.CompanyName,
			&stock.Symbol,
			&stock.Price,
			&stock.Commission,
			&stock.Amount,
			&stock.LatestPrice,
			&stock.LatestUpdate,
			&stock.Close,
			&stock.CloseTime,
			&stock.PERatio,
			&stock.Change,
			&stock.ChangePercent)
		if err != nil {
			log.Println("scanning row failed", err.Error())
		}
		portfolioStocks.Stocks = append(portfolioStocks.Stocks, stock)
	}
	return portfolioStocks
}

func GetTransactions(portfolioid string) domain.PortfolioStocks {
	log.Println("portfolioid", portfolioid)
	rows, err := db.Query(`SELECT
			p.portfolioid,
			po.name,
			quotes.companyName,
			p.symbol,
			p.price,
			p.commission,
			p.amount,
			p.date
		FROM portfolio AS po,
			portfoliostocks AS p
			LEFT JOIN quotes on p.symbol = quotes.symbol
		WHERE p.portfolioid LIKE ($1)
		AND po.portfolioid = p.portfolioid
		ORDER BY p.symbol, p.date ASC;`, portfolioid)
	if err != nil {
		log.Printf("failed with '%s'\n", err)
		return domain.PortfolioStocks{}
	}
	defer rows.Close()
	var portfolioStocks domain.PortfolioStocks
	for rows.Next() {
		stock := domain.PortfolioStock{}
		err := rows.Scan(&stock.Portfolioid,
			&portfolioStocks.PortfolioName,
			&stock.CompanyName,
			&stock.Symbol,
			&stock.Price,
			&stock.Commission,
			&stock.Amount,
			&stock.Date)
		if err != nil {
			log.Println("scanning row failed", err.Error())
		}
		portfolioStocks.Stocks = append(portfolioStocks.Stocks, stock)
	}
	return portfolioStocks
}
