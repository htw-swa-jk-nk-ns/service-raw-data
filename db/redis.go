package db

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/htw-swa-jk-nk-ns/service-raw-data/vote"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type redisDatabase struct {
	db *redis.Client
}

func (c *redisDatabase) insertVote(ctx context.Context, vote vote.Vote) error {
	jsonVote, err := json.Marshal(vote)
	if err != nil {
		return errors.Wrap(err, "failed to marshal vote to json")
	}
	c.db.Set(ctx, vote.ID, jsonVote, 0)
	return nil
}

func (c *redisDatabase) getAllVotes(ctx context.Context) (vote.Votes, error) {
	voteIDs, err := c.db.Keys(ctx, "*").Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan redis database")
	}
	if len(voteIDs) == 0 {
		return vote.Votes{}, nil
	}
	jsonVotes, err := c.db.MGet(ctx, voteIDs...).Result()

	var votes vote.Votes
	for k, jsonVote := range jsonVotes {
		var v vote.Vote
		err := json.Unmarshal([]byte(jsonVote.(string)), &v)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal json vote for vote id '%s'", voteIDs[k])
		}
		votes = append(votes, v)
	}

	return votes, nil
}

func (c *redisDatabase) initialize() error {
	c.db = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	_, err := c.db.Ping(context.TODO()).Result()
	if err != nil {
		return errors.Wrap(err, "failed to connect to redis")
	}
	return nil
}

func (c *redisDatabase) validateAvailableInformation() error {
	return nil
}
