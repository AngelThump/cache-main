package api

import (
	"encoding/json"
	"log"

	utils "github.com/angelthump/cache-main/utils"
	"github.com/go-resty/resty/v2"
)

type Stream struct {
	Created_at string `json:"createdAt"`
	User       User   `json:"user"`
	Ingest     struct {
		Server   string `json:"server"`
		Id       string `json:"id"`
		Mediamtx bool   `json:"mediamtx"`
	} `json:"ingest"`
}

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	StreamKey string `json:"stream_key"`
}

func Find() []Stream {
	client := resty.New()
	resp, _ := client.R().
		SetHeader("X-Api-Key", utils.Config.StreamsAPI.AuthKey).
		Get(utils.Config.StreamsAPI.Hostname + "/streams")

	statusCode := resp.StatusCode()
	if statusCode >= 400 {
		log.Printf("Unexpected status code, got %d %s", statusCode, string(resp.Body()))
		return nil
	}

	var streams []Stream
	err := json.Unmarshal(resp.Body(), &streams)
	if err != nil {
		log.Printf("Unmarshal Error %v", err)
		return nil
	}

	return streams
}

func GetUser(id string) *User {
	client := resty.New()
	resp, _ := client.R().
		SetHeader("X-Api-Key", utils.Config.StreamsAPI.AuthKey).
		Get(utils.Config.StreamsAPI.Hostname + "/users/" + id)

	statusCode := resp.StatusCode()
	if statusCode >= 400 {
		log.Printf("Unexpected status code, got %d %s", statusCode, string(resp.Body()))
		return nil
	}

	var user User
	err := json.Unmarshal(resp.Body(), &user)
	if err != nil {
		log.Printf("Unmarshal Error %v", err)
		return nil
	}

	return &user
}
