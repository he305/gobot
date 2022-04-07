package malv2service

import "time"

type TokenAuthResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AnimeListResponse struct {
	Data []struct {
		AnimeEntry AnimeListEntry `json:"node"`
	} `json:"data"`
	Paging struct {
	} `json:"paging"`
}

type AnimeListEntry struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	MainPicture struct {
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"main_picture"`
	AlternativeTitles struct {
		Synonyms []string `json:"synonyms"`
		En       string   `json:"en"`
		Ja       string   `json:"ja"`
	} `json:"alternative_titles"`
	Broadcast struct {
		DayOfTheWeek string `json:"day_of_the_week"`
		StartTime    string `json:"start_time"`
	} `json:"broadcast"`
	Status       string `json:"status"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	MyListStatus struct {
		Status             string    `json:"status"`
		Score              int       `json:"score"`
		NumEpisodesWatched int       `json:"num_episodes_watched"`
		IsRewatching       bool      `json:"is_rewatching"`
		UpdatedAt          time.Time `json:"updated_at"`
	} `json:"my_list_status"`
	NumEpisodes int `json:"num_episodes"`
}

type AnimePlainResponse struct {
	Data []struct {
		AnimeEntry AnimePlainEntry `json:"node"`
	} `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

type AnimePlainEntry struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	MainPicture struct {
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"main_picture"`
}
