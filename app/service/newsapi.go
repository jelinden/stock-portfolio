package service

import (
	"encoding/json"
	"log"

	"github.com/jelinden/stock-portfolio/app/domain"
	"github.com/jelinden/stock-portfolio/app/util"
)

func GetPortfolioNews(query string) domain.News {
	news := util.Get(`https://www.uutispuro.fi/api/news?q=`+query, 60)
	var marshalled domain.News
	err := json.Unmarshal(news, &marshalled)
	if err != nil {
		log.Println(err)
	}
	return marshalled
}
