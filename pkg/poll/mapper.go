package poll

import (
	"strconv"
	"survey-api/pkg/dtime"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PollMapper struct {
}

func NewPollMapper() PollMapper {
	return PollMapper{}
}

func (mapper PollMapper) ToPoll(pollCreate PollCreate, userId string) (Poll, error) {
	var poll Poll
	creatorId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return poll, err
	}

	closedAt, err := dtime.ISOToDateTime(pollCreate.CloseAt)
	if err != nil {
		return poll, err
	}

	poll = Poll{
		CreatorId:  creatorId,
		Content:    pollCreate.Content,
		Options:    mapper.toPollOptions(pollCreate.Options),
		Visibility: pollCreate.Visibility,
		OpenTill:   closedAt,
	}
	poll.Init()
	return poll, nil
}

func (mapper PollMapper) toPollOptions(createPollOptions []PollCreateOption) []PollOption {
	pollOptions := make([]PollOption, len(createPollOptions))
	for index, item := range createPollOptions {
		pollOptions[index] = mapper.toPollOption(item, strconv.Itoa(index))
	}

	return pollOptions
}

func (mapper PollMapper) toPollOption(createPollOption PollCreateOption, index string) PollOption {
	return PollOption{
		Index:   index,
		Content: createPollOption.Content,
		Count:   0,
	}
}

func (mapper PollMapper) ToPollDetails(poll Poll) PollDetails {
	pollList := mapper.ToPollList(poll)
	return PollDetails{
		PollList: pollList,
		OpenTill: dtime.DateTimeToISO(poll.OpenTill),
		Options:  poll.Options,
	}
}

func (mapper PollMapper) ToPollLists(polls []Poll) []PollList {
	pollLists := make([]PollList, len(polls))
	for index := range polls {
		pollLists[index] = mapper.ToPollList(polls[index])
	}

	return pollLists
}

func (mapper PollMapper) ToPollList(poll Poll) PollList {
	return PollList{
		Id:           poll.Id.Hex(),
		Content:      poll.Content,
		Participants: len(poll.VoterIds),
		State:        string(poll.State),
		Created:      dtime.DateTimeToISO(poll.Created),
		Closed:       dtime.DateTimeToISO(poll.Closed),
	}
}
