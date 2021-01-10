package db

import (
	"context"
	"fmt"
	"github.com/htw-swa-jk-nk-ns/service-raw-data/vote"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"sync"
)

var db struct {
	sync.Once
	Database
}

type database interface {
	insertVote(context.Context, vote.Vote) error
	getAllVotes(context.Context) (vote.Votes, error)
	validateAvailableInformation() error
	initialize() error
}

type Database struct {
	database
}

func (d *Database) InsertVote(ctx context.Context, vote vote.Vote) error {
	err := vote.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid vote")
	}
	return d.insertVote(ctx, vote)
}

func (d *Database) GetAllVotes(ctx context.Context) (vote.Votes, error) {
	return d.getAllVotes(ctx)
}

func (d *Database) initialize() error {
	switch dbType := viper.Get("db.type"); dbType {
	case "redis":
		d.database = &redisDatabase{}
	case "mysql":
		d.database = &mysqlDatabase{}
	default:
		return fmt.Errorf("invalid database drivername '%s'", dbType)
	}
	err := d.validateAvailableInformation()
	if err != nil {
		return errors.Wrap(err, "there are missing information for the given db type")
	}
	err = d.database.initialize()
	if err != nil {
		return errors.Wrap(err, "failed to initialize db")
	}
	return nil
}

func GetDatabase() (*Database, error) {
	var err error
	db.Do(func() {
		err = db.initialize()
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize db")
	}
	return &db.Database, nil
}
