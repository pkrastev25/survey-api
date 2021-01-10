package poll

import validation "github.com/go-ozzo/ozzo-validation/v4"

type PollCreate struct {
	Content    string             `json:"content"`
	Options    []PollCreateOption `json:"options"`
	Visibility PollVisibility     `json:"visibility"`
}

type PollCreateOption struct {
	Content string `json:"content,omitempty"`
}

type PollDetails struct {
	PollList
	Options []PollOption `json:"options"`
}

type PollList struct {
	Id           string `json:"id"`
	Content      string `json:"content"`
	Participants int    `json:"participants"`
	Created      string `json:"created"`
	Closed       string `json:"closed,omitempty"`
}

type PollVote struct {
	PollId string `json:"poll_id"`
	Index  string `json:"index"`
}

func (pollCreate PollCreate) Validate() error {
	return validation.ValidateStruct(&pollCreate,
		validation.Field(&pollCreate.Content, validation.Required),
		validation.Field(&pollCreate.Visibility, validation.Required, validation.In(Public)),
		validation.Field(&pollCreate.Options, validation.Required, validation.Length(2, 8)),
	)
}

func (pollCreateOption PollCreateOption) Validate() error {
	return validation.ValidateStruct(&pollCreateOption,
		validation.Field(&pollCreateOption.Content, validation.Required),
	)
}
