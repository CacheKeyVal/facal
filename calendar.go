package facal

import (
	"strings"
	"time"
)

type Calendar struct {
	workingTimes DefaultWorkingTimes
	exceptions   CalendarExceptions
}

// New creates a new factory calendar
func New(wt DefaultWorkingTimes, excepts CalendarExceptions) (c *Calendar) {
	c = &Calendar{
		workingTimes: wt,
		exceptions:   excepts,
	}
	return
}

type WorkingDay struct {
	Weekday time.Weekday
	Times   *WorkingTimes
}

func (c *Calendar) GetWorkingTimes(d Date) DaytimePeriod {
	return DaytimePeriod{}
}

func (c *Calendar) IsHoliday(t time.Time) bool {
	return c.GetWorkingDays(t, DateOf(t)) != 1
}

func (c *Calendar) GetNearestWorkingDay(dt time.Time) DateWorkingTimes {
	t := TimeOf(dt)

	var isInitDate = true

	// todo check for infinity
	for {
		d := DateOf(dt)
		if e, ok := c.exceptions.Get(d); ok {
			if e.WorkingDayTimes != nil && ((isInitDate && e.WorkingDayTimes.Workday.To.GreaterThan(t)) || !isInitDate) {
				return e
			}
		} else {
			wd := dt.Weekday()
			if c.workingTimes[wd] != nil {
				if (isInitDate && c.workingTimes[wd].Workday.To.GreaterThan(t)) || !isInitDate {
					return makeWorkingTimes(dt, c.workingTimes[wd])
				}
			}
		}

		if isInitDate {
			isInitDate = false
		}
		dt = dt.Add(24 * time.Hour)
	}
}

func (c *Calendar) GetWorkingDays(from time.Time, to Date) int {
	t := TimeOf(from)
	dest := to.In(time.UTC).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var isInitDate = true

	total := 0
	for from.Sub(dest) < 0 {
		d := DateOf(from)
		if e, ok := c.exceptions.Get(d); ok {
			if e.WorkingDayTimes != nil && ((isInitDate && e.WorkingDayTimes.Workday.To.GreaterThan(t)) || !isInitDate) {
				total++
			}
		} else {
			wd := from.Weekday()
			if c.workingTimes[wd] != nil {
				if (isInitDate && c.workingTimes[wd].Workday.To.GreaterThan(t)) || !isInitDate {
					total++
				}
			}
		}

		if isInitDate {
			isInitDate = false
		}
		from = from.Add(24 * time.Hour)
	}
	return total
}

func makeWorkingTimes(dt time.Time, wt *WorkingTimes) DateWorkingTimes {
	return DateWorkingTimes{
		Date:            DateOf(dt),
		WorkingDayTimes: wt,
	}
}

type WorkingTimes struct {
	// Breaks is time period in format 13:00-14:00
	Workday DaytimePeriod

	// Breaks is comma separated time periods in format 13:00-14:00
	Breaks []DaytimePeriod
}

func (wt *WorkingTimes) ToGetInTime(t Time) bool {

	return false
}

type DefaultWorkingTimes map[time.Weekday]*WorkingTimes

type CalendarExceptions []DateWorkingTimes

func (s CalendarExceptions) Get(d Date) (dt DateWorkingTimes, ok bool) {
	for _, e := range s {
		if e.Date.DaysSince(d) == 0 {
			return e, true
		}
	}
	return DateWorkingTimes{}, false
}

type DateWorkingTimes struct {
	Date            Date
	WorkingDayTimes *WorkingTimes
}

func ParseWorkingTimes(wt string, bt []string) *WorkingTimes {
	p := ParseDaytimePeriod(wt)

	t := &WorkingTimes{
		Workday: p,
		Breaks:  make([]DaytimePeriod, 0),
	}

	for _, r := range bt {
		b := ParseDaytimePeriod(r)
		t.Breaks = append(t.Breaks, b)
	}

	return t
}

// ParseDaytimePeriod parses string in format "08:00-18:00" to DaytimePeriod
func ParseDaytimePeriod(s string) DaytimePeriod {
	var (
		from, to Time
		err      error
	)
	sl := strings.Split(s, "-")
	from, err = ParseTime(sl[0])
	if err != nil {
		panic("invalid time")
	}
	to, err = ParseTime(sl[1])
	if err != nil {
		panic("invalid time")
	}

	if to.Hour < from.Hour || (to.Hour == from.Hour && to.Minute < from.Minute) {
		panic("invalid time")
	}

	p := DaytimePeriod{
		From: from,
		To:   to,
	}
	return p
}
