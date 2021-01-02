package dtime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ISOFormat = "2006-01-02T15:04:05.000Z"
)

var (
	NilTime    = primitive.DateTime(0).Time().UTC()
	NilTimeISO = TimeToISO(NilTime)
)

func TimeNow() time.Time {
	return time.Now().UTC()
}

func DateTimeNow() primitive.DateTime {
	return primitive.NewDateTimeFromTime(TimeNow())
}

func DateTimeToISO(dateTime primitive.DateTime) string {
	if NilTime.Equal(dateTime.Time().UTC()) {
		return ""
	}

	return TimeToISO(dateTime.Time())
}

func TimeToISO(time time.Time) string {
	return time.UTC().Format(ISOFormat)
}

func ISOToDateTime(iso string) (primitive.DateTime, error) {
	time, err := ISOToTime(iso)
	if err != nil {
		return 0, err
	}

	return primitive.NewDateTimeFromTime(time), nil
}

func ISOToTime(iso string) (time.Time, error) {
	time, err := time.Parse(ISOFormat, iso)
	if err != nil {
		return NilTime, err
	}

	return time.UTC(), nil
}
