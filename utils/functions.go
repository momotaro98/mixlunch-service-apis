package utils

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func MakeCorrectFormatDateTimeStr(year int, month int, day int, hour int, min int) string {
	// RFC3339 format is accepted.
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:00Z", year, month, day, hour, min)
}

func MakeDateTimeFromStringDate(strDate string) (time.Time, error) {
	d := strings.Split(strDate, "-")
	if len(d) != 3 {
		return time.Time{}, fmt.Errorf("cannot parse passed string")
	}
	year, err := strconv.Atoi(d[0])
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(d[1])
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(d[2])
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}

func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt32(i int32) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}
