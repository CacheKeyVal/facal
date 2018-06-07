package facal

import (
	"encoding/json"
	"time"
)

// A Date represents a date (year, month, day).
//
// This type does not include location information, and therefore does not
// describe a unique 24-hour timespan.
type Date struct {
	year  int        // year (e.g., 2014).
	month time.Month // month of the year (January = 1, ...).
	day   int        // day of the month, starting at 1.
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	date, err := ParseDate(str)
	if err != nil {
		return err
	}
	d.year = date.year
	d.month = date.month
	d.day = date.day
	return nil
}

// Equal returns true if d is equal with cd
func (d Date) Equal(cd Date) bool {
	return d.year == cd.year && d.month == cd.month && d.day == cd.day
}

// In returns the time corresponding to time 00:00:00 of the date in the location.
//
// In is always consistent with time.Date, even when time.Date returns a time
// on a different day. For example, if loc is America/Indiana/Vincennes, then both
//     time.Date(1955, time.May, 1, 0, 0, 0, 0, loc)
// and
//     Date{year: 1955, month: time.May, day: 1}.In(loc)
// return 23:00:00 on April 30, 1955.
//
// In panics if loc is nil.
func (d Date) In(loc *time.Location) time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, loc)
}

// DaysSince returns the signed number of days between the date and s, not including the end day.
func (d Date) DaysSince(s Date) (days int) {
	// We convert to Unix time so we do not have to worry about leap seconds:
	// Unix time increases by exactly 86400 seconds per day.
	deltaUnix := d.In(time.UTC).Unix() - s.In(time.UTC).Unix()
	return int(deltaUnix / 86400)
}

// DateOf returns the Date in which a time occurs in that time's location.
func DateOf(t time.Time) Date {
	var d Date
	d.year, d.month, d.day = t.Date()
	return d
}

// ParseDate parses a string in RFC3339 full-date format and returns the date value it represents.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return DateOf(t), nil
}

// A Time represents a time with nanosecond precision.
//
// This type does not include location information, and therefore does not
// describe a unique moment in time.
//
// This type exists to represent the TIME type in storage-based APIs like BigQuery.
// Most operations on Times are unlikely to be meaningful. Prefer the DateTime type.
type Time struct {
	Hour   int // The hour of the day in 24-hour format; range [0-23]
	Minute int // The minute of the hour; range [0-59]
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	pt, err := ParseTime(str)
	if err != nil {
		return err
	}
	t.Hour = pt.Hour
	t.Minute = pt.Minute
	return nil
}

// Equal returns true if t is equal with ct
func (t Time) Equal(ct Time) bool {
	return t.Hour == ct.Hour && t.Minute == ct.Minute
}

// TimeOf returns true if t is greater than ct
func (t Time) GreaterThan(ct Time) bool {
	return (t.Hour*60 + t.Minute) > (ct.Hour*60 + ct.Minute)
}

// TimeOf returns the Time representing the time of day in which a time occurs
// in that time's location. It ignores the date.
func TimeOf(t time.Time) Time {
	var tm Time
	tm.Hour, tm.Minute, _ = t.Clock()
	return tm
}

// ParseTime parses a string and returns the time value it represents.
// ParseTime accepts an extended form of the RFC3339 partial-time format. After
// the HH:MM:SS part of the string, an optional fractional part may appear,
// consisting of a decimal point followed by one to nine decimal digits.
// (RFC3339 admits only one digit after the decimal point).
func ParseTime(s string) (Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return Time{}, err
	}
	return TimeOf(t), nil
}

type DaytimePeriod struct {
	From Time `json:"d"`
	To   Time `json:"to"`
}
