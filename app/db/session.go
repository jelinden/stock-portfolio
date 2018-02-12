package db

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var boltDB *bolt.DB

const bucketName = "Sessions"

func InitBolt() {
	var err error
	boltDB, err = bolt.Open("/tmp/sessions.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err.Error())
	}
	boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func CloseBolt() {
	boltDB.Close()
}

func PutSession(key string, value string) {
	boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Put([]byte(key), []byte(value))
	})
}

func GetSession(key string) (val string) {
	boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		val = string(b.Get([]byte(key)))
		return nil
	})
	return
}

func RemoveSession(key string) {
	boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete([]byte(key))
	})
}
