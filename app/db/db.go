package db

import (
	"database/sql"
	"log"
	"reflect"
	"time"

	_ "github.com/cznic/ql/driver"

	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/jelinden/stock-portfolio/app/service"
	"github.com/jelinden/stock-portfolio/app/util"
)

const dbFileName = "./ql.db"

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

CREATE TABLE IF NOT EXISTS history (
	symbol string,
	closePriceDate string,
	epoch int64,
	closePrice float64
);
CREATE INDEX IF NOT EXISTS histSymbolIndex ON history (symbol);
CREATE INDEX IF NOT EXISTS histEpochIndex ON history (epoch);
CREATE INDEX IF NOT EXISTS histSymbolEpochIndex ON history (symbol, epoch);`

// ALTER TABLE portfoliostocks ADD epoch int64;
// ALTER TABLE history ADD epoch int64;
// ALTER TABLE portfoliostocks ADD transactionid string;

func init() {
	var err error
	db, err = sql.Open("ql", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db file", dbFileName, "opened")

	mdb, err = sql.Open("ql-mem", "mem.db")
	if err != nil {
		log.Fatal(err)
	}
	tx, _ := mdb.Begin()
	defer recoverFrom(tx)
	_, err = tx.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	initMemoryDatabase()

	if err != nil {
		log.Fatalf("ql.OpenFile() failed with '%s'\n", err.Error())
	}

	populateDatabase()
	go util.DoEvery(time.Hour*12, getHistory)
	go util.DoEvery(time.Second*20, getQuotes)
	go util.DoEvery(time.Minute*180, getDividends)
}

func initMemoryDatabase() {
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
	tx2, _ := mdb.Begin()
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

func exec(command string, args ...interface{}) error {
	tx, err := db.Begin()
	defer recoverFrom(tx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(command, args...)
	if err != nil {
		log.Println("failed executing", command, err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Println("commit msg", err)
	}
	return err
}

func recoverFrom(tx *sql.Tx) {
	if r := recover(); r != nil {
		log.Println("recovered from ", r)
		tx.Rollback()
	}
}

func execRow(command string, args ...interface{}) domain.User {
	row := db.QueryRow(command, args...)
	var user = domain.User{}
	err := row.Scan(&user.ID,
		&user.Email,
		&user.Username,
		&user.RoleName,
		&user.Password,
		&user.CreateDate,
		&user.EmailVerified,
		&user.EmailVerifiedDate,
		&user.EmailVerificationString,
		&user.ModifyDate,
	)
	if err != nil {
		log.Println(err)
	}
	return user
}

func getQuery(query string, args ...interface{}) []map[string]interface{} {
	//log.Println(query, args)
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	var values = make([]interface{}, len(columns))
	var vals []map[string]interface{}
	for i := range values {
		var temp interface{}
		values[i] = &temp
	}

	for rows.Next() {
		m := make(map[string]interface{})
		err := rows.Scan(values...)
		if err != nil {
			log.Println(err)
		}
		for i, colName := range columns {
			var rawValue = *(values[i].(*interface{}))
			var rawType = reflect.TypeOf(rawValue)

			if rawValue != nil && rawType.Name() != "float64" && rawType.Name() != "int64" {
				rawValue = uint8ToString(rawValue.([]uint8))
			}
			m[colName] = rawValue
		}
		vals = append(vals, m)
	}
	if vals != nil && vals[0] != nil {
		return vals
	}
	return nil
}

func uint8ToString(val []uint8) string {
	var valString = ""
	for _, item := range val {
		valString = valString + string(item)
	}
	return valString
}

func queryUser(command string, args ...interface{}) domain.User {
	row := db.QueryRow(command, args...)
	var user = domain.User{}
	err := row.Scan(&user.ID,
		&user.Email,
		&user.Username,
		&user.RoleName,
		&user.Password,
		&user.CreateDate,
		&user.EmailVerified,
		&user.EmailVerifiedDate,
		&user.EmailVerificationString,
		&user.ModifyDate,
	)
	if err != nil {
		log.Println(err)
	}
	return user
}

func queryAllUsers() domain.UserList {
	rows, err := db.Query(`select id, email, username, rolename, createdate, modifydate, emailverified, emailverifieddate, emailverificationstring from user;`)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var userList domain.UserList
	for rows.Next() {
		var user = domain.User{}
		err := rows.Scan(&user.ID,
			&user.Email,
			&user.Username,
			&user.RoleName,
			&user.CreateDate,
			&user.ModifyDate,
			&user.EmailVerified,
			&user.EmailVerifiedDate,
			&user.EmailVerificationString,
		)
		if err != nil {
			log.Println(err)
		}
		userList.Users = append(userList.Users, user)
	}
	return userList
}

func queryPortfolioSymbols() []string {
	rows, err := db.Query(`SELECT distinct symbol FROM portfoliostocks;`)
	if err != nil {
		log.Printf("failed with '%s'\n", err.Error())
	}
	defer rows.Close()
	var symbols []string
	for rows.Next() {
		var symbol string
		rows.Scan(&symbol)
		symbols = append(symbols, symbol)
	}
	return symbols
}

func getQuotes() {
	timeFrom := time.Now()
	quotes := service.GetQuotes(GetPortfolioSymbols()...)
	if len(quotes) > 0 {
		log.Printf("got %v quotes in %v\n", len(quotes), time.Now().Sub(timeFrom))
		SaveQuotes(quotes)
	}
}

func getHistory() {
	closePrices := service.GetClosePrices(GetPortfolioSymbols()...)
	if len(closePrices) > 0 {
		log.Printf("got %v closePrices\n", len(closePrices))
		SaveHistory(closePrices)
	}
}

func After() {
	log.Println("closing db connection")
	err := db.Close()
	if err != nil {
		log.Println("closing db connection error", err.Error())
	}
	log.Println("db connection closed")
}

type History struct {
	Date  int
	Open  float64
	High  float64
	Low   float64
	Close float64
}
