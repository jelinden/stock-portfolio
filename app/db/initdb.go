package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/cznic/ql/driver"
	"github.com/jelinden/stock-portfolio/app/util"
)

const dbFileName = "./ql.db"
const memDB = "mem.db"

var db *sql.DB
var mdb *sql.DB

const createTables = `CREATE TABLE IF NOT EXISTS user (
	id string,
	email string,
	username string,
	password string,
	rolename string,
	emailverified bool,
	emailverificationstring string,
	emailverifieddate string,
	createdate string,
	modifydate string
);
CREATE UNIQUE INDEX IF NOT EXISTS userIdIndex ON user (id);
CREATE UNIQUE INDEX IF NOT EXISTS emailIndex ON user (email);

CREATE TABLE IF NOT EXISTS portfolio (
	portfolioid string,
	userid string,
	name string
);
CREATE INDEX IF NOT EXISTS portfolioUserIdIndex ON portfolio (userid);
CREATE UNIQUE INDEX IF NOT EXISTS portfolioIdIndex ON portfolio (portfolioid);

ALTER TABLE portfoliostocks ADD epoch int64;
CREATE TABLE IF NOT EXISTS portfoliostocks (
	transactionid string,
	portfolioid string,
	userid string,
	symbol string,
	price float64,
	amount int,
	date string,
	epoch int64,
	commission float64
);
CREATE INDEX IF NOT EXISTS portfolioStocksTransactionIdIndex ON portfoliostocks (transactionid, portfolioid);
CREATE UNIQUE INDEX IF NOT EXISTS transactionIdIndex ON portfoliostocks (transactionid);
CREATE INDEX IF NOT EXISTS portfolioStocksEpochIndex ON portfoliostocks (epoch);
CREATE INDEX IF NOT EXISTS portfolioIdSymbolIndex ON portfoliostocks (portfolioid, symbol);
CREATE INDEX IF NOT EXISTS portfolioSymbolEpochIndex ON portfoliostocks (symbol, epoch);
CREATE INDEX IF NOT EXISTS xportfoliostocks_portfolioid ON portfoliostocks(portfolioid);

CREATE TABLE IF NOT EXISTS instrument (
	symbol string
);
CREATE UNIQUE INDEX IF NOT EXISTS instrumentSymbolIndex ON instrument (symbol);

CREATE TABLE IF NOT EXISTS quotes (
	symbol string,
	companyName string,
	sector string,
	close float64,
	closeTime int64,
	latestPrice float64,
	latestUpdate int64,
	change float64,
	changePercent float64,
	PERatio float64
);
CREATE UNIQUE INDEX IF NOT EXISTS quotesSymbolIndex ON quotes (symbol);

CREATE TABLE IF NOT EXISTS dividend (
	symbol string,
	exDate int64,
	paymentDate int64,
	amount float64,
	type string
);
CREATE INDEX IF NOT EXISTS divSymbolIndex ON dividend (symbol);
CREATE UNIQUE INDEX IF NOT EXISTS divSymbolPaymentIndex ON dividend (paymentDate, symbol, type);
CREATE INDEX IF NOT EXISTS divPaymentDateIndex ON dividend (paymentDate);

DROP TABLE history;
CREATE TABLE IF NOT EXISTS history (
	symbol string,
	closePriceDate string,
	epoch int64,
	closePrice float64
);
CREATE INDEX IF NOT EXISTS histSymbolIndex ON history (symbol);
CREATE INDEX IF NOT EXISTS histEpochIndex ON history (epoch);
CREATE INDEX IF NOT EXISTS histSymbolEpochIndex ON history (symbol, epoch);`

// ALTER TABLE history ADD epoch int64;
// ALTER TABLE portfoliostocks ADD transactionid string;

func initFileDB() {
	var err error
	db, err = sql.Open("ql", dbFileName)
	if err != nil {
		log.Fatal("fatal", err)
	}
	tx, _ := db.Begin()
	defer recoverFrom(tx)
	_, err = tx.Exec(createTables)
	if err != nil {
		log.Fatal("fatal error creating tables", err)
	}
	tx.Commit()
	log.Println("db file", dbFileName, "opened")
}

func initMemDatabase() {
	var err error
	mdb, err = sql.Open("ql-mem", memDB)
	if err != nil {
		log.Fatal(err)
	}
	tx, _ := mdb.Begin()
	defer recoverFrom(tx)
	_, err = tx.Exec(createTables)
	if err != nil {
		log.Fatal("fatal error creating tables", err)
	}
	tx.Commit()
}

func populateMemoryDatabase() {
	transactions := getQuery(`select * from portfoliostocks;`)
	tx, _ := mdb.Begin()
	defer recoverFrom(tx)
	for _, item := range transactions {
		// get all transactions to memory db
		_, err := tx.Exec(`insert into portfoliostocks (
					transactionid,
					portfolioid,
					userid,
					symbol,
					price,
					amount,
					date,
					epoch,
					commission) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			item["transactionid"],
			item["portfolioid"],
			item["userid"],
			item["symbol"],
			item["price"],
			item["amount"],
			item["date"],
			item["epoch"],
			item["commission"])
		if err != nil {
			log.Println(err)
		}
	}
	tx.Commit()
	log.Println("updated mdb portfoliostocks", len(transactions), "items")

	historyItems := getQuery(`select * from history;`)
	tx2, err := mdb.Begin()
	if err != nil {
		log.Println(err)
	}
	defer recoverFrom(tx2)
	for _, item := range historyItems {
		// update all history lines memory db
		_, err := tx2.Exec(`insert into history (symbol, closePriceDate, epoch, closePrice) values ($1,$2,$3,$4)`,
			item["symbol"],
			item["closePriceDate"],
			item["epoch"],
			item["closePrice"])
		if err != nil {
			log.Println(err)
		}
	}
	tx2.Commit()
	log.Println("updated mdb history", len(historyItems), "items")
}

func populateDatabase() {
	err := exec(createTables)
	if err != nil {
		log.Printf("failed creating user table '%s'", err.Error())
	}
	vals := getQuery(`select * from portfoliostocks where transactionid is null;`)
	for _, item := range vals {
		log.Println(item["portfolioid"], item["date"], item["symbol"], item["amount"], item["price"], item["transactionid"])
		// update all transactions to have a transactionid
		exec(`update portfoliostocks set transactionid = $1 where portfolioid = $2 and date = $3 and amount = $4 and price = $5 and symbol = $6`,
			util.GetTimeBasedID(), item["portfolioid"], item["date"], item["amount"], item["price"], item["symbol"])
		time.Sleep(10 * time.Millisecond)
	}

	transactionDates := getQuery(`select transactionid, date from portfoliostocks where epoch is null;`)
	for _, item := range transactionDates {
		// update all transactions to have an epoch timestamp
		d, err := time.Parse("01/02/2006", item["date"].(string))
		if err != nil {
			log.Println("parsing date failed", err)
		} else {
			epoch := d.Unix() * 1000
			exec(`update portfoliostocks set epoch = $1 where transactionid = $2`, epoch, item["transactionid"])
			time.Sleep(5 * time.Millisecond)
			log.Println("updated portfoliostocks", item["transactionid"], item["date"], epoch)
		}
	}

	historyDates := getQuery(`select symbol, closePriceDate from history where epoch is null;`)
	for _, item := range historyDates {
		// update all history lines to have an epoch timestamp
		d, err := time.Parse("01/02/2006", item["closePriceDate"].(string))
		if err != nil {
			log.Println("parsing closePriceDate failed", err)
		} else {
			epoch := d.Unix() * 1000
			exec(`update history set epoch = $1 where symbol = $2 and closePriceDate = $3`, epoch, item["symbol"], item["closePriceDate"])
			time.Sleep(5 * time.Millisecond)
			log.Println("updated history", item["symbol"], item["closePriceDate"], epoch)
		}
	}
}
