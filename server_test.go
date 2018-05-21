package main

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var connectionCount = 0
var randomSleepMillis = []int{55, 70, 75, 85, 110, 93, 89, 97, 82, 74, 210}
var urls = []string{
	"http://localhost:3300/login",
	"http://localhost:3300/signup",
	"http://localhost:3300/api/user",
}

func init() {
	go main()
	waitForConnection()
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(2*time.Second))
}

func httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: dialTimeout,
		},
	}
}

func TestTraffic(t *testing.T) {
	var counter = 0
	for {
		rand.Seed(time.Now().UnixNano())
		url := urls[rand.Int()%len(urls)]
		go get(t, url)
		log.Printf("getting %v", url)
		if counter == 60000 {
			break
		}
		counter++
		time.Sleep(time.Duration(randomSleepMillis[rand.Int()%len(randomSleepMillis)]) * time.Millisecond)
	}
}

func get(t *testing.T, url string) {
	resp, err := httpClient().Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	assert.True(t, resp.StatusCode == 200, "status code was not 200")
}

func waitForConnection() {
	if connectionCount < 10 {
		resp, err := http.Get("http://localhost:3300/login")
		if err != nil || resp.StatusCode != 200 {
			time.Sleep(1 * time.Second)
			connectionCount++
			waitForConnection()
		}
	} else {
		log.Println("Connecting to http://localhost:3300/login failed")
		os.Exit(1)
	}
}
