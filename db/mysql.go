package db

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/htw-swa-jk-nk-ns/service-raw-data/vote"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var mysqlSchemaArr = []string{
	`DROP TABLE IF EXISTS votes;`,

	`CREATE TABLE votes (
		name varchar(255) NOT NULL,
		country varchar(255) NOT NULL,
		candidate varchar(255) NOT NULL,
		date datetime DEFAULT current_timestamp,
		);`,
}

type mysqlDatabase struct {
	db *sqlx.DB
}

func (d *mysqlDatabase) insertVote(ctx context.Context, vote vote.Vote) error {
	sb := sqlbuilder.NewInsertBuilder()
	sb.InsertInto("values")
	sb.Cols("name", "country", "candidate", "date")
	sb.Values(vote.Name, vote.Country, vote.Candidate, vote.Date)
	query, err := sqlbuilder.MySQL.Interpolate(sb.Build())
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}
	_, err = d.db.ExecContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}
	return nil
}

func (d *mysqlDatabase) getAllVotes(_ context.Context) (vote.Votes, error) {
	var votes vote.Votes
	err := d.db.Select(&votes, d.db.Rebind("SELECT * FROM votes"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute sql query")
	}

	return votes, nil
}

func (d *mysqlDatabase) initialize() error {
	var err error
	d.db, err = sqlx.ConnectContext(context.Background(), "mysql", viper.GetString("mysql.dataSourceName"))
	if err != nil {
		return errors.Wrap(err, "failed to connect to mysql database")
	}
	for _, query := range mysqlSchemaArr {
		_, err := d.db.Exec(query)
		if err != nil {
			_, _ = d.db.Exec(`DROP TABLE IF EXISTS cache;`)
			return errors.Wrap(err, "Could not set up database schema - query: "+query)
		}
	}
	return nil
}

func (d *mysqlDatabase) validateAvailableInformation() error {
	return nil
}
