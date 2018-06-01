package facal

import (
	"reflect"
	"testing"
	"time"
)

func parseDateTime(s string) time.Time {
	t, _ := time.Parse(time.RFC822, s)
	return t
}

func TestNextWorkingDay(t *testing.T) {
	defaultWorkingTimes := DefaultWorkingTimes{
		Monday:    ParseWorkingTimes("07:00-18:00", []string{"12:00-13:00"}),
		Tuesday:   ParseWorkingTimes("07:00-18:00", []string{"12:00-13:00"}),
		Wednesday: ParseWorkingTimes("07:00-18:00", []string{"12:00-13:00"}),
		Thursday:  ParseWorkingTimes("07:00-18:00", []string{"12:00-13:00"}),
		Friday:    ParseWorkingTimes("07:00-18:00", []string{"12:00-13:00", "16:00-16:30"}),
		Saturday:  nil,
		Sunday:    nil,
	}

	workingTimesExceptions := []DateWorkingTimes{
		{
			Date:            Date{2018, time.May, 31},
			WorkingDayTimes: ParseWorkingTimes("07:00-20:00", []string{"13:00-14:00", "15:45-16:00"}),
		},
		{
			Date:            Date{2018, time.June, 1},
			WorkingDayTimes: nil,
		},
	}

	federalHolidays := []Date{
		{2018, time.January, 1},
		{2018, time.January, 2},
		{2018, time.January, 3},
		{2018, time.January, 4},
		{2018, time.January, 5},
		{2018, time.January, 6},
		{2018, time.January, 7},
		{2018, time.January, 8},
		{2018, time.February, 23},
		{2018, time.February, 24},
		{2018, time.February, 25},
		{2018, time.March, 8},
		{2018, time.March, 9},
		{2018, time.March, 10},
		{2018, time.March, 11},
		{2018, time.April, 29},
		{2018, time.April, 30},
		{2018, time.May, 1},
		{2018, time.May, 2},
		{2018, time.May, 9},
		{2018, time.June, 10},
		{2018, time.June, 11},
		{2018, time.June, 12},
		{2018, time.November, 3},
		{2018, time.November, 4},
		{2018, time.November, 5},
		{2018, time.December, 31},
	}

	for _, fh := range federalHolidays {
		workingTimesExceptions = append(workingTimesExceptions, DateWorkingTimes{
			Date:            fh,
			WorkingDayTimes: nil,
		})
	}

	oksarCalendar := New(defaultWorkingTimes, workingTimesExceptions)

	testCases := []struct {
		t    time.Time
		want Date
	}{
		{parseDateTime("02 Jan 18 15:04 UTC"), Date{2018, time.January, 9}},
		{parseDateTime("04 Jan 18 15:04 UTC"), Date{2018, time.January, 9}},
		{parseDateTime("08 May 18 19:04 UTC"), Date{2018, time.May, 10}},
		{parseDateTime("31 May 18 19:30 UTC"), Date{2018, time.May, 31}},
		{parseDateTime("31 May 18 05:30 UTC"), Date{2018, time.May, 31}},
		{parseDateTime("31 May 18 20:00 UTC"), Date{2018, time.June, 4}},
		{parseDateTime("01 Jun 18 12:04 UTC"), Date{2018, time.June, 4}},
		{parseDateTime("02 Jun 18 18:04 UTC"), Date{2018, time.June, 4}},
		{parseDateTime("03 Jun 18 22:04 UTC"), Date{2018, time.June, 4}},
		{parseDateTime("28 Dec 18 16:00 UTC"), Date{2018, time.December, 28}},
		{parseDateTime("28 Dec 18 19:00 UTC"), Date{2019, time.January, 1}},
		{parseDateTime("31 Dec 18 16:00 UTC"), Date{2019, time.January, 1}},
		{parseDateTime("31 Dec 18 22:04 UTC"), Date{2019, time.January, 1}},
	}

	for _, testCase := range testCases {
		res := oksarCalendar.NextWorkingDay(testCase.t)
		if !reflect.DeepEqual(res.Date, testCase.want) {
			t.Errorf("NextWorkingDay() want %v, but %v given", testCase.want, res.Date)
		}
	}
}

func TestParseDaytimePeriod(t *testing.T) {
	testCases := []struct {
		s    string
		want DaytimePeriod
	}{
		{"08:00-18:00", DaytimePeriod{
			Start:    Time{8, 0},
			End:      Time{18, 0},
			Duration: 10 * time.Hour,
		}},
		{"09:15-17:15", DaytimePeriod{
			Start:    Time{9, 15},
			End:      Time{17, 15},
			Duration: 8 * time.Hour,
		}},
	}

	for _, testCase := range testCases {
		res := ParseDaytimePeriod(testCase.s)
		if !res.Start.Equal(testCase.want.Start) || !res.End.Equal(testCase.want.End) || res.Duration != testCase.want.Duration {
			t.Errorf("ParseDaytimePeriod() wantErr %v, %v", res, testCase.want)
		}
	}
}