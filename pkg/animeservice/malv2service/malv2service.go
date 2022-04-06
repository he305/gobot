package malv2service

import (
	"encoding/json"
	"fmt"
	as "gobot/pkg/animeservice"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const basePath string = "https://api.myanimelist.net/v2"
const clientId string = "6114d00ca681b7701d1e15fe11a4987e"
const timeLayout string = "2021-03-26"

var animeListRequestFields = [7]string{
	"alternative_titles", "broadcast", "status", "start_date", "end_date", "my_list_status", "num_episodes",
}

var headers = map[string][]string{
	"Content-Type":    {"application/x-www-form-urlencoded"},
	"X-MAL-Client-ID": {clientId},
	"User-Agent":      {"NineAnimator/2 CFNetwork/976 Darwin/18.2.0"},
	"Connection":      {"Keep-Alive"},
}

var airingStatusMap = map[string]uint8{
	"currently_airing": as.Airing,
	"finished_airing":  as.CompletedAiring,
	"not_yet_aired":    as.NotStarted,
}

var listStatusMap = map[string]uint8{
	"plan_to_watch": as.PlannedToWatch,
	"completed":     as.Completed,
	"watching":      as.Watching,
	"dropped":       as.Dropped,
}

type malv2service struct {
	tokenInfo TokenAuthResponse
	username  string
	password  string
	client    *resty.Client
}

var _ as.AnimeService = (*malv2service)(nil)

func (serv *malv2service) GetAnimeByTitle(title string) *as.AnimeStruct {
	if err := serv.verifyToken(); err != nil {
		fmt.Println(err)
	}

	resp, err := serv.client.R().
		SetAuthToken(serv.tokenInfo.AccessToken).
		SetHeaderMultiValues(headers).
		SetQueryParam("q", title).
		Get(basePath + "/anime")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	return nil
}

func (serv *malv2service) GetUserAnimeList() []*as.AnimeStruct {
	if err := serv.verifyToken(); err != nil {
		fmt.Println(err)
	}

	fieldStr := strings.Join(animeListRequestFields[:], ",")

	resp, err := serv.client.R().
		SetAuthToken(serv.tokenInfo.AccessToken).
		SetHeaderMultiValues(headers).
		SetQueryParams(map[string]string{
			"limit":  "999",
			"offset": "0",
			"fields": fieldStr,
			"sort":   "anime_title",
		}).
		Get(basePath + "/users/@me/animelist")

	if err != nil {
		fmt.Println(err)
	}

	var respJson AnimeListResponse
	if err := json.Unmarshal(resp.Body(), &respJson); err != nil {
		fmt.Println(err)
	}

	var animeList []*as.AnimeStruct
	for _, v := range respJson.Data {
		entry := v.AnimeEntry

		var synonyms []string
		for _, altTitle := range entry.AlternativeTitles.Synonyms {
			if altTitle != "" {
				synonyms = append(synonyms, altTitle)
			}
		}
		if entry.AlternativeTitles.En != "" {
			synonyms = append(synonyms, entry.AlternativeTitles.En)
		}
		if entry.AlternativeTitles.Ja != "" {
			synonyms = append(synonyms, entry.AlternativeTitles.Ja)
		}

		layout := time.RFC3339[:len(entry.StartDate)]
		parsedStartTime, err := time.Parse(layout, strings.TrimSpace(entry.StartDate))
		if err != nil {
			parsedStartTime = time.Now()
		}

		parsedEndTime, err := time.Parse(layout, strings.TrimSpace(entry.EndDate))
		if err != nil {
			parsedEndTime = time.Now()
		}

		an := as.NewAnimeStruct(
			entry.ID,
			entry.Title,
			synonyms,
			parsedStartTime,
			parsedEndTime,
			float64(entry.MyListStatus.Score),
			airingStatusMap[entry.Status],
			listStatusMap[entry.MyListStatus.Status],
			entry.MainPicture.Large,
		)
		animeList = append(animeList, an)
	}

	return animeList
}

func NewMalv2Service(username string, password string) as.AnimeService {
	return &malv2service{username: username, password: password, tokenInfo: TokenAuthResponse{}, client: resty.New()}
}

func (serv *malv2service) verifyToken() error {
	if serv.tokenInfo.AccessToken == "" || serv.tokenInfo.RefreshToken == "" {
		return serv.getToken()
	}

	resp, err := serv.client.R().
		SetBody(fmt.Sprintf(`client_id=%s&grant_type=refresh_token&refresh_token=%s`, clientId, serv.tokenInfo.RefreshToken)).
		SetHeaderMultiValues(headers).
		Post("https://myanimelist.net/v1/oauth2/token")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Couldn't get a token in malv2service, status code %d", resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &serv.tokenInfo); err != nil {
		return err
	}
	fmt.Println("Parsed refresh Token!")
	return nil
}

func (serv *malv2service) getToken() error {
	resp, err := serv.client.R().
		SetBody(fmt.Sprintf(`client_id=%s&grant_type=password&password=%s&username=%s`, clientId, serv.password, serv.username)).
		SetHeaderMultiValues(headers).
		Post(basePath + "/auth/token")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Couldn't get a token in malv2service, status code %d", resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &serv.tokenInfo); err != nil {
		return fmt.Errorf("Couldn't parse token auth response")
	}
	fmt.Println("Parsed Token!")

	return nil
}
