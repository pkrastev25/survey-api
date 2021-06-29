package poll

import (
	"errors"
	"strconv"
	"survey-api/pkg/dtime"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	minPollOptions = 2
	maxPollOptions = 8
)

type PollCreate struct {
	Content    string             `json:"content"`
	Options    []PollCreateOption `json:"options"`
	Visibility PollVisibility     `json:"visibility"`
	CloseAt    string             `json:"close_at"`
}

type PollCreateOption struct {
	Content string `json:"content"`
}

type PollDetails struct {
	PollList `json:",omitempty"`
	OpenTill string       `json:"open_till"`
	Options  []PollOption `json:"options"`
}

type PollList struct {
	Id           string `json:"id"`
	Content      string `json:"content"`
	Participants int    `json:"participants"`
	State        string `json:"state"`
	Created      string `json:"created"`
	Closed       string `json:"closed,omitempty"`
}

type PollVote struct {
	Index string `json:"index"`
}

type PollModify struct {
	State string `json:"state"`
}

func (pollCreate PollCreate) Validate() error {
	return validation.ValidateStruct(&pollCreate,
		validation.Field(&pollCreate.Content, validation.Required, validation.Length(1, 500)),
		validation.Field(&pollCreate.Visibility, validation.Required, validation.In(VisibilityPublic)),
		validation.Field(&pollCreate.Options, validation.Required, validation.Length(minPollOptions, maxPollOptions)),
		validation.Field(&pollCreate.CloseAt, validation.By(validateCloseAt)),
	)
}

func validateCloseAt(value interface{}) error {
	closeAtString, ok := value.(string)
	if !ok {
		return errors.New("")
	}

	closeAt, err := dtime.ISOToTime(closeAtString)
	if err != nil {
		return err
	}

	if closeAt.After(dtime.TimeNow().Add(time.Hour * time.Duration(24))) {
		return errors.New("")
	}

	if closeAt.Before(dtime.TimeNow().Add(time.Minute * time.Duration(5))) {
		return errors.New("")
	}

	return nil
}

func (pollCreateOption PollCreateOption) Validate() error {
	return validation.ValidateStruct(&pollCreateOption,
		validation.Field(&pollCreateOption.Content, validation.Required, validation.Length(1, 250)),
	)
}

func (pollVote PollVote) Validate() error {
	return validation.ValidateStruct(&pollVote,
		validation.Field(&pollVote.Index, validation.Required, validation.By(validatePollVoteIndex)),
	)
}

func validatePollVoteIndex(value interface{}) error {
	indexString, ok := value.(string)
	if !ok {
		return errors.New("")
	}

	index, err := strconv.Atoi(indexString)
	if err != nil {
		return err
	}

	if index < 0 {
		return errors.New("")
	}

	if index > maxPollOptions {
		return errors.New("")
	}

	return nil
}

func (pollModify PollModify) Validate() error {
	return validation.ValidateStruct(&pollModify,
		validation.Field(&pollModify.State, validation.Required, validation.In(StateClosed)),
	)
}
