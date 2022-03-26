package helper

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const loc = "America/New_York"

func GetEDTTime() time.Time {
	location, err := time.LoadLocation(loc)
	if err != nil {
		fmt.Println(err)
	}
	return time.Now().In(location)
}

func GetNewDateTime() primitive.DateTime {
	return primitive.NewDateTimeFromTime(GetEDTTime())
}

func ISODateFormat(dateTime time.Time) string {
	return dateTime.Format(time.RFC3339)
}

func GetUTCDate(t time.Time, hoursToAdd int) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, hoursToAdd, 0, 0, 0, time.UTC)
}
