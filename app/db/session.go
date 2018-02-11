package db

import (
	"github.com/orcaman/concurrent-map"
)

var cache = cmap.New()

func PutSession(key string, value string) {
	cache.Set(key, value)
}

func GetSession(key string) string {
	value, ok := cache.Get(key)
	if ok {
		return value.(string)
	}
	return ""
}

func RemoveSession(key string) {
	cache.Remove(key)
}
