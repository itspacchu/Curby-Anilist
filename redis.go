package main

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

func UpdateEpoch(rdb *redis.Client, key string, newbs basicStore) error {
	data, _ := json.Marshal(newbs)
	err := rdb.Set(key, data, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
