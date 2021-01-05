package model

import (
	"strconv"
	"survey-api/pkg/dtime"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Public PollVisibility = "public"
)

type PollVisibility string

type Poll struct {
	Id           primitive.ObjectID   `bson:"_id,omitempty"`
	CreatorId    primitive.ObjectID   `bson:"creator_id"`
	Content      string               `bson:"content"`
	Options      []PollOption         `bson:"options"`
	Visibility   PollVisibility       `bson:"visibility"`
	VoterIds     []primitive.ObjectID `bson:"voter_ids"`
	Created      primitive.DateTime   `bson:"created"`
	Closed       primitive.DateTime   `bson:"closed"`
	LastModified primitive.DateTime   `bson:"last_modified"`
}

type PollOption struct {
	Index   string `bson:"index",json:"index"`
	Content string `bson:"content",json:"content"`
	Count   int    `bson:"count",json:"count"`
}

type CreatePoll struct {
	Content    string             `json:"content"`
	Options    []CreatePollOption `json:"options"`
	Visibility PollVisibility     `json:"visibility"`
}

type CreatePollOption struct {
	Content string `json:"content,omitempty"`
}

type PollClient struct {
	Id           string       `json:"id"`
	Content      string       `json:"content"`
	Options      []PollOption `json:"options"`
	Participants int          `json:"participants"`
	Created      string       `json:"created"`
	Closed       string       `json:"closed,omitempty"`
}

type PollVote struct {
	PollId string `json:"poll_id"`
	Index  string `json:"index"`
}

func (createPoll CreatePoll) ToPoll(userId string) (Poll, error) {
	var poll Poll
	creatorId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return poll, err
	}

	pollOptions := make([]PollOption, len(createPoll.Options))

	for index, item := range createPoll.Options {
		pollOptions[index] = item.ToPollOption(strconv.Itoa(index))
	}

	return Poll{
		CreatorId:  creatorId,
		Content:    createPoll.Content,
		Options:    pollOptions,
		Visibility: createPoll.Visibility,
		Created:    dtime.DateTimeNow(),
	}, nil
}

func (createPollOption CreatePollOption) ToPollOption(index string) PollOption {
	return PollOption{
		Index:   index,
		Content: createPollOption.Content,
		Count:   0,
	}
}

func (poll Poll) ToPollClient() PollClient {
	return PollClient{
		Id:           poll.Id.Hex(),
		Content:      poll.Content,
		Options:      poll.Options,
		Participants: len(poll.VoterIds),
		Created:      dtime.DateTimeToISO(poll.Created),
		Closed:       dtime.DateTimeToISO(poll.Closed),
	}
}

func (createPoll CreatePoll) Validate() error {
	return validation.ValidateStruct(&createPoll,
		validation.Field(&createPoll.Content, validation.Required),
		validation.Field(&createPoll.Visibility, validation.Required, validation.In(Public)),
		validation.Field(&createPoll.Options, validation.Required, validation.Length(2, 8)),
	)
}

func (createPollOption CreatePollOption) Validate() error {
	return validation.ValidateStruct(&createPollOption,
		validation.Field(&createPollOption.Content, validation.Required),
	)
}

func (pollVote PollVote) Validate() error {
	return validation.ValidateStruct(&pollVote,
		validation.Field(&pollVote.PollId, validation.Required),
		validation.Field(&pollVote.Index, validation.Required),
	)
}
