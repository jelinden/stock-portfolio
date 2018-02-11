package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/cznic/ql"
	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/jelinden/stock-portfolio/app/service"
)

const dbFileName = "./ql.db"

var db *sql.DB

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
	portfolioid string,
	userid string,
	symbol string,
	price float64,
	amount int,
	date string,
	commission float64
);
CREATE INDEX IF NOT EXISTS portfolioStocksUserIdIndex ON portfoliostocks (portfolioid);

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
CREATE INDEX IF NOT EXISTS divSymbolIndex ON dividend (symbol);`

func Init() {
	var err error
	ql.RegisterDriver()
	db, err = sql.Open("ql", dbFileName)
	if err != nil {
		log.Fatalf("ql.OpenFile() failed with '%s'\n", err.Error())
	}
	log.Println("db file", dbFileName, "opened")
	populateDatabase()
	go doEvery(time.Second*60, getQuotes)
}

func populateDatabase() {
	err := exec(createTables)
	if err != nil {
		log.Printf("failed creating user table '%s'", err.Error())
	}
}

func exec(command string, args ...interface{}) error {
	//log.Println(command, args)
	tx, err := db.Begin()
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
	return err
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
	quotes := service.GetQuotes(GetPortfolioSymbols()...)
	if len(quotes) > 0 {
		log.Println("got quotes", quotes[0].CompanyName)
	}
	SaveQuotes(quotes)
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func After() {
	log.Println("closing db connection")
	err := db.Close()
	if err != nil {
		log.Println("closing db connection error", err.Error())
	}
}
