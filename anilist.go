package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const ANILIST_API string = "https://graphql.anilist.co"

func GetActivityDetailsForPosting(aid int, actDetails *ActivityDetails) error {
	QUERY := fmt.Sprintf("query GetActivityDetails {\n  Activity(id: %d) {\n    ... on ListActivity {\n      type\n      media {\n        coverImage {\n          large\n  color \n       }\n        title {\n          english\n          romaji\n        }\n        updatedAt\n      }\n      status\n      progress\n    }\n  }\n}", aid)
	body := map[string]interface{}{
		"query": QUERY,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.Post(ANILIST_API, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.StatusCode > 250 {
		response, _ := io.ReadAll(req.Body)
		fmt.Println(string(response))
		return fmt.Errorf("Unavailable? %d\n", req.StatusCode)

	}
	response, err := io.ReadAll(req.Body)
	json.Unmarshal(response, actDetails)
	return nil
}

// Returns ActivityID , ActivityEpoch, Errors if any
func FetchUserActivityIDAfterEpoch(uid int, epoch int64) (int, int64, error) {
	QUERY := fmt.Sprintf(
		"query GetRecentActivity {  Activity(userId: %d type_in: [ANIME_LIST, MANGA_LIST] createdAt_greater: %d) {... on ListActivity {status siteUrl createdAt}}}",
		uid,
		epoch,
	)
	body := map[string]interface{}{
		"query": QUERY,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, 0, err
	}
	req, err := http.Post(ANILIST_API, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, 0, err
	}
	defer req.Body.Close()
	if req.StatusCode > 250 {
		return 0, 0, fmt.Errorf("Unavailable? %d", req.StatusCode)
	}
	response, err := io.ReadAll(req.Body)
	adata := ActivityData{}
	json.Unmarshal(response, &adata)
	actList := strings.Split(adata.Data.Activity.SiteUrl, "/")
	actID, _ := strconv.Atoi(actList[len(actList)-1])
	return actID, adata.Data.Activity.CreatedAt, nil
}

func UserIdFetchFromUsername(username string) (int, string, error) {
	QUERY := `query GetUserID {User(name: "REPLACE_ME"){id avatar{medium}}}`
	QUERY = strings.ReplaceAll(QUERY, "REPLACE_ME", username)
	body := map[string]interface{}{
		"query": QUERY,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, "", err
	}
	req, err := http.Post(ANILIST_API, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, "", err
	}
	defer req.Body.Close()
	response, err := io.ReadAll(req.Body)
	respStruct := UserDataID{}
	json.Unmarshal(response, &respStruct)
	return respStruct.Data.User.Id, respStruct.Data.User.Avatar.Medium, nil
}
