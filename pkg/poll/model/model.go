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
	OwnerId      primitive.ObjectID   `bson:"creator_id,omitempty"`
	Content      string               `bson:"content,omitempty"`
	Options      []PollOption         `bson:"options,omitempty"`
	Visibility   PollVisibility       `bson:"visibility,omitempty"`
	VoterIds     []primitive.ObjectID `bson:"voter_ids,omitempty"`
	Created      primitive.DateTime   `bson:"created,omitempty"`
	Closed       primitive.DateTime   `bson:"closed,omitempty"`
	LastModified primitive.DateTime   `bson:"last_modified,omitempty"`
}

type PollOption struct {
	Index   string `bson:"index,omitempty",json:"index"`
	Content string `bson:"content,omitempty",json:"content"`
	Count   int    `bson:"count,omitempty",json:"count"`
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

func (p *CreatePoll) ToPoll(userId string) (*Poll, error) {
	creatorId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	pollOptions := make([]PollOption, len(p.Options))

	for index, item := range p.Options {
		pollOptions[index] = *item.ToPollOption(strconv.Itoa(index))
	}

	poll := &Poll{
		Id:         primitive.NewObjectID(),
		OwnerId:    creatorId,
		Content:    p.Content,
		Options:    pollOptions,
		Visibility: p.Visibility,
		Created:    dtime.DateTimeNow(),
	}

	return poll, nil
}

func (po *CreatePollOption) ToPollOption(index string) *PollOption {
	return &PollOption{
		Index:   index,
		Content: po.Content,
		Count:   0,
	}
}

func (p *Poll) ToPollClient() *PollClient {
	return &PollClient{
		Id:           p.Id.Hex(),
		Content:      p.Content,
		Options:      p.Options,
		Participants: len(p.VoterIds),
		Created:      dtime.DateTimeToISO(p.Created),
		Closed:       dtime.DateTimeToISO(p.Closed),
	}
}

func (p CreatePoll) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Content, validation.Required),
		validation.Field(&p.Visibility, validation.Required, validation.In(Public)),
		validation.Field(&p.Options, validation.Required, validation.Length(2, 8)),
	)
}

func (po CreatePollOption) Validate() error {
	return validation.ValidateStruct(&po,
		validation.Field(&po.Content, validation.Required),
	)
}

func (pv PollVote) Validate() error {
	return validation.ValidateStruct(&pv,
		validation.Field(&pv.PollId, validation.Required),
		validation.Field(&pv.Index, validation.Required),
	)
}
