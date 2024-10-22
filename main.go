package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/switchupcb/disgo"
	"gopkg.in/yaml.v3"
)

func main() {
	settingsPtr, err := os.Open("settings.yaml")
	defer settingsPtr.Close()
	if err != nil {
		log.Fatalf("Unable to find settings.yaml file")
	}
	settingsText, err := io.ReadAll(settingsPtr)

	SettingsYaml := secrets{}
	err = yaml.Unmarshal(settingsText, &SettingsYaml)
	if err != nil {
		log.Fatalf("%s", err)
	}

	bot := &disgo.Client{
		ApplicationID:  "",
		Authentication: disgo.BotToken(SettingsYaml.DISCORD_TOKEN),
		Config:         disgo.DefaultConfig(),
		Handlers:       new(disgo.Handlers),
		Sessions:       disgo.NewSessionManager(),
	}

	session := disgo.NewSession()

	if err := session.Connect(bot); err != nil {
		log.Printf("can't open websocket session to Discord Gateway: %v\n", err)
		return
	}
	log.Println("Connected to Discord")

	rdb := redis.NewClient(&redis.Options{
		Addr:     SettingsYaml.REDIS_ENDPOINT,
		Password: SettingsYaml.REDIS_TOKEN,
		DB:       0,
	})
	log.Println("Connected to Redis")

	var wg sync.WaitGroup
	var cursor uint64

	pattern := "discord:*"
	exitCode := make(chan os.Signal, 1)
	ticker := time.NewTicker(5 * time.Minute)

	for range ticker.C {
		for {
			keys, cursor, err := rdb.Scan(cursor, pattern, 10).Result()
			if err != nil {
				log.Fatalf("Could not retrieve keys: %v", err)
			}
			for _, key := range keys {
				if err != nil {
					log.Println("Unable to fetch Unmarshal basicStore data")
				}
				wg.Add(1)
				go StartPollingUser(key, &wg, bot, rdb)
				fmt.Printf("Goroutine Started for %s\n", key)
			}
			if cursor == 0 {
				break
			}
		}
	}
	signal.Notify(exitCode, os.Interrupt)

	<-exitCode
	log.Println("^C Recived exitting!")
}
