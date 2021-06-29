package poll

import (
	"survey-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	VisibilityPublic PollVisibility = "public"
)

const (
	StateOpen   PollState = "open"
	StateClosed PollState = "closed"
)

const (
	propertCreatorId = "creator_id"
	propertyVoterIds = "voter_ids"
	propertOptions   = "options"
	propertyCount    = "count"
	propertyState    = "state"
	propertyClosed   = "closed"
	propertyOpenTill = "open_till"
)

type PollVisibility string
type PollState string

type Poll struct {
	db.BaseModel `bson:",inline"`
	CreatorId    primitive.ObjectID   `bson:"creator_id"`
	Content      string               `bson:"content"`
	Options      []PollOption         `bson:"options"`
	VoterIds     []primitive.ObjectID `bson:"voter_ids"`
	Visibility   PollVisibility       `bson:"visibility"`
	State        PollState            `bson:"state"`
	OpenTill     primitive.DateTime   `bson:"open_till"`
	Closed       primitive.DateTime   `bson:"closed,omitempty"`
}

type PollOption struct {
	Index   string `bson:"index",json:"index"`
	Content string `bson:"content",json:"content"`
	Count   int    `bson:"count",json:"count"`
}

func (poll *Poll) Init() {
	poll.BaseModel = db.NewBaseModel()
	poll.VoterIds = []primitive.ObjectID{}
	poll.State = StateOpen
}
