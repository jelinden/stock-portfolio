package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQuotes(t *testing.T) {
	quotes := GetQuotes("XOM")
	assert.True(t, quotes[0].Symbol == "XOM", "symbol should be XOM")
}

func TestGetDividends(t *testing.T) {
	dividends := GetDividends("XOM")
	assert.True(t, dividends[0].Symbol == "XOM", "symbol should be XOM")
}

func TestGetClosePrices(t *testing.T) {
	closePrices := GetClosePrices("ADP")
	assert.True(t, closePrices[0].Symbol == "ADP", "symbol should be ADP")
}
