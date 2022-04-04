package malv2service

import (
	"encoding/json"
	"fmt"
	as "gobot/pkg/animeservice"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const basePath string = "https://api.myanimelist.net/v2"
const clientId string = "6114d00ca681b7701d1e15fe11a4987e"

var headers = map[string][]string{
	"Content-Type":    {"application/x-www-form-urlencoded"},
	"X-MAL-Client-ID": {clientId},
	"User-Agent":      {"NineAnimator/2 CFNetwork/976 Darwin/18.2.0"},
	"Connection":      {"Keep-Alive"},
}

type malv2service struct {
	tokenInfo TokenAuth
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

func NewMalv2Service(username string, password string) as.AnimeService {
	return &malv2service{username: username, password: password, tokenInfo: TokenAuth{}, client: resty.New()}
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
