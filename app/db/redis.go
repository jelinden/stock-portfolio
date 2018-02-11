package db

/*
import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jelinden/stock-portfolio/app/domain"

	"gopkg.in/redis.v3"
)

type Redis struct {
	SessionClient     *redis.Client
	UserInfoClient    *redis.Client
	PortfolioClient   *redis.Client
	StocksClient      *redis.Client
	TransactionClient *redis.Client
}

func NewRedis() *Redis {
	return &Redis{}
}

func (r *Redis) Init() {
	r.SessionClient = r.createRedis(10)
	r.UserInfoClient = r.createRedis(12)
	r.PortfolioClient = r.createRedis(4)
	r.StocksClient = r.createRedis(5)
	r.TransactionClient = r.createRedis(6)
}

func (r *Redis) createRedis(db int64) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       db,
	})
	return client
}

func (r *Redis) AddStock(key string, value string, stock domain.Quote) {
	portfolio := r.GetStocks(key)
	for i := range portfolio.Quotes {
		if len(portfolio.Quotes) >= i+1 && portfolio.Quotes[i].Symbol == stock.Symbol {
			if len(portfolio.Quotes) == i+1 {
				portfolio.Quotes = portfolio.Quotes[:len(portfolio.Quotes)-1]
			} else {
				portfolio.Quotes = append(portfolio.Quotes[:i], portfolio.Quotes[i+1:]...)
			}
		}
	}

	quote, err := json.Marshal(portfolio.Quotes)
	if err != nil {
		fmt.Println("marshalling stock failed", err)
	}
	quoteString := string(quote)
	quoteString = strings.Replace(quoteString, "[", "", -1)
	quoteString = strings.Replace(quoteString, "]", "", -1)
	quoteString = strings.Replace(quoteString, "},{", "}{", -1)
	err = r.PortfolioClient.Set(key, quoteString, -1).Err() //key,value,expiration in time.Hour
	if err != nil {
		fmt.Println("PUT ERROR", err)
	}

	err = r.PortfolioClient.Append(key, value).Err()
	if err != nil {
		fmt.Println("PUT ERROR", err)
	}
}

func (r *Redis) GetStocks(key string) domain.Portfolio {
	val, err := r.PortfolioClient.Get(key).Result()
	if err == redis.Nil {
		fmt.Println(key, "does not exists")
		return domain.Portfolio{}
	} else if err != nil {
		fmt.Println("GET ERROR", err)
		return domain.Portfolio{}
	}
	stocks := []domain.Quote{}
	portfolio := domain.Portfolio{}
	err = json.Unmarshal([]byte("["+strings.Replace(val, "}{", "},{", -1)+"]"), &stocks)
	if err != nil {
		fmt.Println(err)
	}
	portfolio.Quotes = stocks
	return portfolio
}

func (r *Redis) Put(key string, value string) {
	err := r.SessionClient.Set(key, value, time.Duration(24*7*4)*time.Hour).Err() //key,value,expiration in time.Hour
	if err != nil {
		fmt.Println("PUT ERROR", err)
	}
}

func (r *Redis) AddNewUser(key string, value string) error {
	exists, err := r.UserInfoClient.Exists(key).Result()
	if err != nil {
		fmt.Println("ERROR", err)
		return &domain.CustomError{Type: "EmailExists", Message: err.Error()}
	}
	if exists == false {
		err := r.UserInfoClient.Set(key, value, 0).Err() //key,value,expiration in time.Hour
		if err != nil {
			fmt.Println("PUT ERROR", err)
			return err
		}
	} else {
		fmt.Println("Email already exists")
	}
	return nil
}

func (r *Redis) UpdateUser(key string, value string) error {
	err := r.UserInfoClient.Set(key, value, 0).Err() //key,value,expiration in time.Hour
	if err != nil {
		fmt.Println("PUT ERROR", err)
		return err
	}
	return nil
}

func (r *Redis) DbSize() int64 {
	return r.UserInfoClient.DbSize().Val()
}

func (r *Redis) GetSession(key string) string {
	return r.get(key, r.SessionClient)
}

func (r *Redis) RemoveSession(key string) {
	err := r.SessionClient.Del(key).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func (r *Redis) GetUser(key string) domain.User {
	user := domain.User{}
	json.Unmarshal([]byte(r.get(key, r.UserInfoClient)), &user)
	return user
}

func (r *Redis) ListUsers() []domain.Member {
	keys := r.UserInfoClient.Keys("*").Val()
	userList := []domain.Member{}

	for _, value := range keys {
		user := domain.Member{}
		json.Unmarshal([]byte(r.get(value, r.UserInfoClient)), &user)
		userList = append(userList, user)
	}
	return userList
}

func (r *Redis) get(key string, client *redis.Client) string {
	val, err := client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println(key, "does not exists")
		return ""
	} else if err != nil {
		fmt.Println("GET ERROR", err)
	}
	return val
}
*/
