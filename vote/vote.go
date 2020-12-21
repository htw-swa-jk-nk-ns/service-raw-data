package vote

import (
	"errors"
)

type Vote struct {
	ID        string
	Name      string
	Country   string
	Candidate string
	Date      int64
}

func (v *Vote) Validate() error {
	if v.ID == "" {
		return errors.New("id cannot be empty")
	}
	if v.Name == "" {
		return errors.New("name cannot be empty")
	}
	if v.Country == "" {
		return errors.New("country cannot be empty")
	}
	if v.Candidate == "" {
		return errors.New("candidate cannot be empty")
	}
	if v.Date <= 0 {
		return errors.New("invalid date")
	}

	return nil
}

type Votes []Vote
