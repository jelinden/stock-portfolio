package service

import (
	"fmt"
	"testing"
)

func TestGetQuotes(t *testing.T) {
	fmt.Println(GetQuotes("ADP", "OHI", "O", "XOM"))
}
