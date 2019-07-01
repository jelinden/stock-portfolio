package util

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func newTimeoutClient(readWriteTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(30*time.Second, readWriteTimeout),
		},
	}
}

func Get(url string, timeoutSeconds int) []byte {
	client := newTimeoutClient(time.Duration(timeoutSeconds) * time.Second)
	res, err := client.Get(url)
	if err != nil {
		log.Println("client get failed", err.Error())
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("reading body failed", err.Error())
	}
	return body
}
