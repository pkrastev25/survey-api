package poll

import (
	"survey-api/pkg/db"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Public PollVisibility = "public"
)

const (
	propertCreatorId = "creator_id"
	propertyVoterIds = "voter_ids"
	propertOptions   = "options"
	propertyCount    = "count"
)

type PollVisibility string

type Poll struct {
	db.BaseModel `bson:",inline"`
	CreatorId    primitive.ObjectID   `bson:"creator_id"`
	Content      string               `bson:"content"`
	Options      []PollOption         `bson:"options"`
	VoterIds     []primitive.ObjectID `bson:"voter_ids"`
	Visibility   PollVisibility       `bson:"visibility"`
	Closed       primitive.DateTime   `bson:"closed,omitempty"`
}

type PollOption struct {
	Index   string `bson:"index",json:"index"`
	Content string `bson:"content",json:"content"`
	Count   int    `bson:"count",json:"count"`
}

func NewPoll() Poll {
	return Poll{
		BaseModel: db.NewBaseModel(),
	}
}

func (poll *Poll) Init() {
	poll.BaseModel = db.NewBaseModel()
	poll.VoterIds = []primitive.ObjectID{}
}

func (pollVote PollVote) Validate() error {
	return validation.ValidateStruct(&pollVote,
		validation.Field(&pollVote.PollId, validation.Required),
		validation.Field(&pollVote.Index, validation.Required),
	)
}
