package domain

import (
	"fmt"
)

const Admin = "admin"
const Normal = "normal"

type UserList struct {
	Users []User `json:"users"`
}

type User struct {
	ID                      string `json:"-"`
	Email                   string `json:"email"`
	Username                string `json:"username"`
	Password                string `json:"-"`
	RoleName                string `json:"role"`
	EmailVerified           bool   `json:"verified"`
	EmailVerificationString string `json:"-"`
	EmailVerifiedDate       string `json:"verifieddate"`
	CreateDate              string `json:"createdate"`
	ModifyDate              string `json:"modifydate"`
}

type CustomError struct {
	Type    string
	Message string
}

func (e CustomError) Error() string {
	return fmt.Sprintf("%v: %v", e.Type, e.Message)
}

type PortfolioStocks struct {
	Stocks        []PortfolioStock `json:"stocks,omitempty"`
	PortfolioName string           `json:"portfolioName"`
}

type PortfolioStock struct {
	Portfolioid   string   `json:"portfolioid"`
	Symbol        string   `json:"symbol"`
	Price         float64  `json:"price"`
	Amount        int      `json:"amount"`
	Commission    float64  `json:"commission"`
	Date          string   `json:"date"`
	CompanyName   *string  `json:"companyName,omitempty"`
	LatestPrice   *float64 `json:"latestPrice,omitempty"`
	LatestUpdate  *float64 `json:"latestUpdate,omitempty"`
	Close         *float64 `json:"close,omitempty"`
	CloseTime     *int64   `json:"closeTime,omitempty"`
	PERatio       *float64 `json:"peRatio,omitempty"`
	Change        *float64 `json:"change,omitempty"`
	ChangePercent *float64 `json:"changePercent,omitempty"`
}

type Portfolios struct {
	Portfolios []Portfolio `json:"portfolios,omitempty"`
}

type Portfolio struct {
	Portfolioid string `json:"portfolioid"`
	Userid      string `json:"userid"`
	Name        string `json:"name"`
}
