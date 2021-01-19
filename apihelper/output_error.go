package apihelper

import "github.com/rs/zerolog/log"

type OutputError struct {
	Message string
}

func NewOutputError(err error) OutputError {
	log.Error().Err(err).Msg("error during request")
	return OutputError{
		Message: err.Error(),
	}
}
