package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var ctx = context.Background()

func main() {
	ExampleClient()
}

type basicStore struct {
	AnilistUsername string `json:"anilist_username"`
	Epoch           int64  `json:"epoch"`
}

func setAnilistUser(rdb *redis.Client, discordUsername string, anilistUsername string) error {
	if len(discordUsername) < 1 || len(anilistUsername) < 1 {
		return fmt.Errorf("oi put in some values")
	}
	m := basicStore{
		anilistUsername,
		time.Now().Unix(),
	}
	data, _ := json.Marshal(m)
	err := rdb.Set("discord:"+discordUsername, data, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-12572.c244.us-east-1-2.ec2.redns.redis-cloud.com:12572",
		Password: "yuZMviHTkIjmkRxkoxRtGOO5kwlldDab",
		DB:       0,
	})
	setAnilistUser(rdb, "jayanth2292", "Jayanthpavan")
	// UpdateEpoch(rdb, "discord:.pacchu", basicStore{
	// 	AnilistUsername: "pacchu",
	// 	Epoch:           1729337496,
	// })

	var cursor uint64
	pattern := "discord:*"
	for {
		keys, cursor, err := rdb.Scan(cursor, pattern, 10).Result()
		if err != nil {
			log.Fatalf("Could not retrieve keys: %v", err)
		}
		for _, key := range keys {
			if err != nil {
				fmt.Println("Unable to fetch Unmarshal basicStore data")
			}
			val, err := rdb.Get(key).Result()
			if err != nil {
				panic(err)
			}
			bs := basicStore{}
			json.Unmarshal([]byte(val), &bs)
			fmt.Printf("%s - %d\n", bs.AnilistUsername, bs.Epoch)
		}
		if cursor == 0 {
			break
		}
	}

}

func UpdateEpoch(rdb *redis.Client, key string, newbs basicStore) error {
	data, _ := json.Marshal(newbs)
	err := rdb.Set(key, data, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
