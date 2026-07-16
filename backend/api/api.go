package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"

	"anidesk/backend/models"
)

const baseURL = "https://api.myanimelist.net/v2"

type Client struct {
	http *resty.Client
}

func New(token string) *Client {
	c := resty.New().
		SetBaseURL(baseURL).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetTimeout(15 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)
	return &Client{http: c}
}

func (c *Client) SetToken(token string) {
	c.http.SetAuthToken(token)
}

func (c *Client) GetAnime(id int64, fields string) (*models.MALAnime, error) {
	if fields == "" {
		fields = "id,title,main_picture,alternative_titles,synopsis,mean,rank,popularity,media_type,status,num_episodes,start_season,genres,studios,rating,pictures"
	}
	resp, err := c.http.R().
		SetQueryParam("fields", fields).
		Get(fmt.Sprintf("/anime/%d", id))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("MAL API error: %s", resp.String())
	}
	var anime models.MALAnime
	if err := json.Unmarshal(resp.Body(), &anime); err != nil {
		return nil, err
	}
	return &anime, nil
}

func (c *Client) SearchAnime(q string, limit int, offset int) ([]models.MALAnime, error) {
	return c.search("/anime", q, limit, offset)
}

func (c *Client) SearchManga(q string, limit int, offset int) ([]models.MALAnime, error) {
	return c.search("/manga", q, limit, offset)
}

func (c *Client) search(path, q string, limit, offset int) ([]models.MALAnime, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	resp, err := c.http.R().
		SetQueryParams(map[string]string{
			"q":      q,
			"limit":  fmt.Sprintf("%d", limit),
			"offset": fmt.Sprintf("%d", offset),
			"fields": "id,title,main_picture,alternative_titles,synopsis,mean,rank,popularity,media_type,status,num_episodes,start_season,genres,studios,rating",
		}).
		Get(path)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("search error: %s", resp.String())
	}
	var result struct {
		Data []struct {
			Node models.MALAnime `json:"node"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	out := make([]models.MALAnime, len(result.Data))
	for i, item := range result.Data {
		out[i] = item.Node
	}
	return out, nil
}

func (c *Client) GetUser() (*models.MALUser, error) {
	resp, err := c.http.R().
		SetQueryParam("fields", "id,name,picture").
		Get("/users/@me")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("user fetch error: %s", resp.String())
	}
	var user models.MALUser
	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) GetSeasonalAnime(year int, season string, limit int, offset int) ([]models.MALAnime, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	resp, err := c.http.R().
		SetQueryParams(map[string]string{
			"limit":  fmt.Sprintf("%d", limit),
			"offset": fmt.Sprintf("%d", offset),
			"fields": "id,title,main_picture,synopsis,mean,rank,popularity,media_type,status,num_episodes,genres,studios,rating",
		}).
		Get(fmt.Sprintf("/anime/season/%d/%s", year, season))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("seasonal error: %s", resp.String())
	}
	var result struct {
		Data []struct {
			Node models.MALAnime `json:"node"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	out := make([]models.MALAnime, len(result.Data))
	for i, item := range result.Data {
		out[i] = item.Node
	}
	return out, nil
}

func (c *Client) GetRanking(rankingType string, limit int, offset int) ([]models.MALAnime, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	resp, err := c.http.R().
		SetQueryParams(map[string]string{
			"ranking_type": rankingType,
			"limit":        fmt.Sprintf("%d", limit),
			"offset":       fmt.Sprintf("%d", offset),
			"fields":       "id,title,main_picture,synopsis,mean,rank,popularity,media_type,status,num_episodes,genres,studios,rating",
		}).
		Get("/anime/ranking")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("ranking error: %s", resp.String())
	}
	var result struct {
		Data []struct {
			Node models.MALAnime `json:"node"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	out := make([]models.MALAnime, len(result.Data))
	for i, item := range result.Data {
		out[i] = item.Node
	}
	return out, nil
}

func (c *Client) GetAnimeList(userName string, status string, limit int, offset int) ([]models.MALAnime, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	params := url.Values{}
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("offset", fmt.Sprintf("%d", offset))
	params.Set("fields", "id,title,main_picture,synopsis,mean,rank,popularity,media_type,status,num_episodes,genres,studios,rating,list_status{status,score,num_episodes_watched,is_rewatching,start_date,finish_date,priority,notes}")
	if status != "" {
		params.Set("status", status)
	}
	resp, err := c.http.R().
		SetQueryParamsFromValues(params).
		Get(fmt.Sprintf("/users/%s/animelist", userName))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("animelist error: %s", resp.String())
	}
	var result struct {
		Data []struct {
			Node models.MALAnime `json:"node"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	out := make([]models.MALAnime, len(result.Data))
	for i, item := range result.Data {
		out[i] = item.Node
	}
	return out, nil
}

func (c *Client) UpdateAnimeStatus(animeID int64, status string, score int, episodesWatched int) error {
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	if score > 0 {
		params.Set("score", fmt.Sprintf("%d", score))
	}
	if episodesWatched > 0 {
		params.Set("num_watched_episodes", fmt.Sprintf("%d", episodesWatched))
	}
	resp, err := c.http.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(params.Encode()).
		Patch(fmt.Sprintf("/anime/%d/my_list_status", animeID))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("update error: %s", resp.String())
	}
	return nil
}

func (c *Client) DeleteAnimeFromList(animeID int64) error {
	resp, err := c.http.R().Delete(fmt.Sprintf("/anime/%d/my_list_status", animeID))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("delete error: %s", resp.String())
	}
	return nil
}

func fetchImage(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
