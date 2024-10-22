package main

type secrets struct {
	DISCORD_TOKEN      string `yaml:"discord_token"`
	DISCORD_GUILD_ID   string `yaml:"discord_guild_id"`
	DISCORD_CHANNEL_ID string `yaml:"discord_channel_id"`
	REDIS_ENDPOINT     string `yaml:"redis_endpoint"`
	REDIS_TOKEN        string `yaml:"redis_token"`
}

type basicStore struct {
	AnilistUsername string `json:"anilist_username"`
	Epoch           int64  `json:"epoch"`
}

type UserDataID struct {
	Data struct {
		User struct {
			Id     int `json:"id"`
			Avatar struct {
				Medium string `json:"medium"`
			} `json:"avatar"`
		} `json:"User"`
	} `json:"data"`
}

type ActivityData struct {
	Data struct {
		Activity struct {
			Status    string `json:"status"`
			SiteUrl   string `json:"siteUrl"`
			CreatedAt int64  `json:"createdAt"`
		} `json:"Activity"`
	} `json:"data"`
}

type ActivityDetails struct {
	Data struct {
		Activity struct {
			Type  *string `json:"type"`
			Media struct {
				CoverImage struct {
					Large *string `json:"large"`
					Color string  `json:"color"`
				} `json:"coverImage"`
				Title struct {
					English *string `json:"english"`
					Romaji  *string `json:"romaji"`
				} `json:"title"`
				UpdatedAt int64 `json:"updatedAt"`
			} `json:"media"`
			Status   *string `json:"status"`
			Progress *string `json:"progress"`
		} `json:"Activity"`
	} `json:"data"`
}
