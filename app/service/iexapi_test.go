package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/jelinden/stock-portfolio/app/config"
	"github.com/stretchr/testify/assert"
)

func TestGetQuotes(t *testing.T) {
	quotes := GetQuotes("XOM")
	assert.True(t, quotes[0].Symbol == "XOM", "symbol should be XOM")
}

func TestGetDividends(t *testing.T) {
	config.Config.FromEmail = os.Getenv("FROMEMAIL")
	config.Config.EmailSendingPasswd = os.Getenv("EMAILSENDINGPASSWD")
	config.Config.AdminUser = os.Getenv("ADMINUSER")
	config.Config.Token = os.Getenv("IEXAPITOKEN")
	dividends := GetDividends("XOM")
	assert.True(t, dividends[0].Symbol == "XOM", "symbol should be XOM")
}

func TestGetDividendsTRI(t *testing.T) {
	config.Config.FromEmail = os.Getenv("FROMEMAIL")
	config.Config.EmailSendingPasswd = os.Getenv("EMAILSENDINGPASSWD")
	config.Config.AdminUser = os.Getenv("ADMINUSER")
	config.Config.Token = os.Getenv("IEXAPITOKEN")
	dividends := GetDividends("TRI")
	assert.True(t, dividends[0].Symbol == "TRI", "symbol should be TRI")
	fmt.Println(dividends[0])
	assert.True(t, dividends[0].Amount >= 0.571247)
	assert.True(t, dividends[0].Currency == "CAD")
	assert.True(t, dividends[0].Symbol == "TRI")
}

func TestGetClosePrices(t *testing.T) {
	closePrices := GetClosePrices("ADP")
	assert.True(t, closePrices[0].Symbol == "ADP", "symbol should be ADP")
}
