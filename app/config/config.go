package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var Config struct {
	VerifyURL          string
	FromEmail          string
	EmailSendingPasswd string
	AdminUser          string
}

func SetConfig(file string) {
	conf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = json.Unmarshal(conf, &Config)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
