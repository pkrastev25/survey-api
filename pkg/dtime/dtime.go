package dtime

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ISOFormat = "2006-01-02T15:04:05.000Z"
)

var (
	nilDateTime = primitive.DateTime(0)
)

func ConvertDateTimeToString(dt primitive.DateTime) string {
	if nilDateTime.Time().Equal(dt.Time()) {
		return ""
	}

	return dt.Time().Format(ISOFormat)
}
