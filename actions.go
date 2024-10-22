package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/switchupcb/disgo"
)

func StartPollingUser(dbKey string, wg *sync.WaitGroup, bot *disgo.Client, rdb *redis.Client) {
	defer wg.Done()
	val, err := rdb.Get(dbKey).Result()
	if err != nil {
		log.Println("Unable to fetch user key data")
	}
	bs := basicStore{}
	err = json.Unmarshal([]byte(val), &bs)
	fmt.Printf("[%s] Fetching for epoch %d\n", dbKey, bs.Epoch)
	useEpoch := bs.Epoch
	useEpoch, err = FetchAndSendLatestActivityFromUsernameAfterEpoch(bs.AnilistUsername, useEpoch, bot, "1297163463583465482")
	fmt.Printf("[%s] Got Updated epoch %d\n", dbKey, useEpoch)
	time.Sleep(time.Second * time.Duration(5))
	bs.Epoch = useEpoch
	UpdateEpoch(rdb, dbKey, bs)
}

func FetchAndSendLatestActivityFromUsernameAfterEpoch(username string, epoch int64, bot *disgo.Client, channelID string) (int64, error) {

	userId, _, err := UserIdFetchFromUsername(username)
	if err != nil {
		log.Printf("%s\nError with fetching details username : %s\n", err, username)
		return epoch, err
	}
	activityID, newEpoch, err := FetchUserActivityIDAfterEpoch(userId, epoch)
	if err != nil {
		log.Printf("%s\nError with fetching activity with id : %d\n", err, activityID)
		return epoch, err
	}
	actdet := ActivityDetails{}
	_ = GetActivityDetailsForPosting(activityID, &actdet)
	getChannelRequest := disgo.GetChannel{ChannelID: channelID}
	_, err = getChannelRequest.Send(bot)
	if err != nil {
		log.Printf("error occurred getting channel %q: %v\n", channelID, err)
		return epoch, err
	}
	anilistActivityURL := fmt.Sprintf("https://anilist.co/activity/%d", activityID)
	Status, Progress, Eng_title, Rom_title, CoverImageURL := "", "", "", "", ""
	if actdet.Data.Activity.Progress != nil {
		Progress = *actdet.Data.Activity.Progress
	}
	if actdet.Data.Activity.Status != nil {
		Status = *actdet.Data.Activity.Status
	}
	if actdet.Data.Activity.Media.Title.English != nil {
		Eng_title = *actdet.Data.Activity.Media.Title.English
	}
	if actdet.Data.Activity.Media.Title.Romaji != nil {
		Rom_title = *actdet.Data.Activity.Media.Title.Romaji
	}
	if actdet.Data.Activity.Media.CoverImage.Large != nil {
		CoverImageURL = *actdet.Data.Activity.Media.CoverImage.Large
	}
	title := fmt.Sprintf("%s %s %s", strings.Title(username), Status, Progress)
	desc := fmt.Sprintf("%s (%s)", Eng_title, Rom_title)
	thumb := &disgo.EmbedThumbnail{URL: CoverImageURL}

	embedData := &disgo.Embed{}
	embedData.Title = &title
	embedData.Description = &desc
	embedData.Thumbnail = thumb
	embedData.URL = &anilistActivityURL
	// embedData.Image.URL = CoverImageURL

	createMessageRequest := &disgo.CreateMessage{
		ChannelID:        channelID,
		Content:          nil,
		Nonce:            nil,
		TTS:              nil,
		Embeds:           []*disgo.Embed{embedData},
		AllowedMentions:  nil,
		MessageReference: nil,
		Components:       nil,
		StickerIDS:       nil,
		Files:            nil,
		Attachments:      nil,
		Flags:            nil,
	}
	message, err := createMessageRequest.Send(bot)
	log.Printf("Successfully sent message with ID %q\n", message.ID)
	return newEpoch, nil
}
