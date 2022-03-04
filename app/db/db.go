package db

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/jelinden/stock-portfolio/app/service"
	"github.com/jelinden/stock-portfolio/app/util"
)

func Init() {
	initFileDB()
	initMemDatabase()

	populateMemoryDatabase()
	populateDatabase()

	//go util.DoEvery(time.Hour*12, getHistory)
	go util.DoEvery(time.Minute*15, getQuotes)
	go util.DoEvery(time.Hour*24, getDividends)
	if os.Getenv("divs") == "run" {
		getDividends()
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
	now := time.Now()
	// get quotes only when the stock exchange is open
	if int(now.Weekday()) != 0 && int(now.Weekday()) != 6 { // https://golang.org/pkg/time/#Weekday
		if now.Hour() >= 16 && now.Hour() < 24 {
			quotes := service.GetQuotes(GetPortfolioSymbols()...)
			if len(quotes) > 0 {
				log.Printf("got %v quotes in %v\n", len(quotes), now.Sub(now))
				SaveQuotes(quotes)
			}
		} else {
			log.Println("time", now.Hour(), "was not between 16-24")
		}
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
