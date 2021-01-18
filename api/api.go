package api

import (
	"context"
	"github.com/htw-swa-jk-nk-ns/service-raw-data/db"
	"github.com/htw-swa-jk-nk-ns/service-raw-data/vote"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

// StartAPI starts the API.
func StartAPI() {
	// initialize database
	_, err := db.GetDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start api")
	}

	e := echo.New()

	e.POST("/vote", postVote)

	e.GET("/all", all)

	if viper.GetString("api.certfile") != "" && viper.GetString("api.keyfile") != "" {
		e.Logger.Fatal(e.StartTLS(":"+viper.GetString("api.port"), viper.GetString("api.certfile"), viper.GetString("api.keyfile")))
	} else {
		e.Logger.Fatal(e.Start(":" + viper.GetString("api.port")))
	}
}

func postVote(ctx echo.Context) error {
	v := vote.Vote{}
	if err := ctx.Bind(&v); err != nil {
		return getApiResponse(ctx, http.StatusBadRequest, newOutputError(errors.Wrap(err, "failed to bind input")))
	}
	v.Date = time.Now().Unix()

	database, err := db.GetDatabase()
	if err != nil {
		return getApiResponse(ctx, http.StatusBadGateway, newOutputError(errors.Wrap(err, "failed to get database")))
	}
	err = database.InsertVote(context.TODO(), v)
	if err != nil {
		return getApiResponse(ctx, http.StatusBadGateway, newOutputError(errors.Wrap(err, "failed to insert vote into database")))
	}

	return getApiResponse(ctx, http.StatusOK, nil)
}

func all(ctx echo.Context) error {
	database, err := db.GetDatabase()
	if err != nil {
		return getApiResponse(ctx, http.StatusBadGateway, newOutputError(errors.Wrap(err, "failed to get database")))
	}
	votes, err := database.GetAllVotes(context.TODO())
	if err != nil {
		return getApiResponse(ctx, http.StatusBadGateway, newOutputError(errors.Wrap(err, "failed to get all votes")))
	}
	if votes == nil {
		votes = vote.Votes{}
	}
	return getApiResponse(ctx, http.StatusOK, votes)
}

func getApiResponse(ctx echo.Context, statusCode int, response interface{}) error {
	switch format := viper.GetString("api.format"); format {
	case "json":
		return ctx.JSON(statusCode, response)
	case "xml":
		return ctx.XML(statusCode, response)
	default:
		return ctx.String(http.StatusInternalServerError, "invalid output format '"+format+"'")
	}
}

type OutputError struct {
	message string
}

func newOutputError(err error) OutputError {
	return OutputError{
		message: err.Error(),
	}
}
