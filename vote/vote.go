package vote

import (
	"errors"
)

type Vote struct {
	Name      string `json:"name" xml:"name"`
	Country   string `json:"country" xml:"country"`
	Candidate string `json:"candidate" xml:"candidate"`
	Date      int64  `json:"date" xml:"date"`
}

func (v *Vote) Validate() error {
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
