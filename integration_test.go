package main

// This test needs a redis database to be running at localhost:6793 with database 1 unused and no entries

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/htw-swa-jk-nk-ns/service-raw-data/cmd"
	"github.com/htw-swa-jk-nk-ns/service-raw-data/db"
	"github.com/htw-swa-jk-nk-ns/service-raw-data/vote"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	viper.Set("api.port", 8887)
	viper.Set("api.format", "json")
	viper.Set("redis.addr", "localhost:6379")
	viper.Set("db.type", "redis")
	viper.Set("redis.db", 1)

	_, err := db.GetDatabase()
	if !assert.NoError(t, err, "failed to get db, redis is probably not running on localhost:6379") {
		return
	}

	go cmd.Execute()
	time.Sleep(1 * time.Second)

	client := resty.New()

	v1 := generateRandomVote()
	v2 := generateRandomVote()
	v3 := generateRandomVote()
	v4 := generateRandomVote()
	v5 := generateRandomVote()

	voteIDs, err := insertVotes(client, v1, v2, v3, v4, v5)
	if !assert.NoError(t, err, "failed to insert votes") {
		return
	}

	votes, err := getAllVotes(client)
	if !assert.NoError(t, err, "failed to get all votes") {
		return
	}

	for _, voteID := range voteIDs {
		found := false
		for _, v := range votes {
			if v.ID == voteID {
				found = true
				break
			}
		}
		if assert.True(t, found, "one vote is missing! id: %s", voteID) {
			return
		}
	}
}

func insertVotes(client *resty.Client, votes ...vote.Vote) ([]string, error) {
	var voteIDs []string
	for _, v := range votes {
		request := client.R()
		request.SetHeader("Content-Type", "application/"+viper.GetString("api.format"))
		jsonVote, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		request.SetBody(string(jsonVote))
		result, err := request.Post(fmt.Sprintf("http://localhost:%d/vote", viper.GetInt("api.port")))
		if err != nil {
			return nil, errors.Wrap(err, "failed to post api request")
		}
		if result.IsError() {
			return nil, fmt.Errorf("put request returned an error: %s", string(result.Body()))
		}
		var voteID string
		err = json.Unmarshal(result.Body(), &voteID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal api response")
		}
		voteIDs = append(voteIDs, voteID)
	}
	return voteIDs, nil
}

func getAllVotes(client *resty.Client) (vote.Votes, error) {
	request := client.R()
	request.SetHeader("Content-Type", "application/"+viper.GetString("api.format"))
	result, err := request.Get(fmt.Sprintf("http://localhost:%d/all", viper.GetInt("api.port")))
	if err != nil {
		return nil, errors.Wrap(err, "failed to post api request")
	}
	if result.IsError() {
		return nil, fmt.Errorf("put request returned an error: %s", string(result.Body()))
	}
	var votes vote.Votes
	err = json.Unmarshal(result.Body(), &votes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal result into vote.Votes")
	}
	return votes, nil
}

func generateRandomVote() vote.Vote {
	return vote.Vote{
		ID:        "",
		Name:      xid.New().String(),
		Country:   xid.New().String(),
		Candidate: xid.New().String(),
		Date:      time.Now().Unix(),
	}
}
