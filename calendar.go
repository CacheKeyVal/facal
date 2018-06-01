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

func (c *Calendar) IsHoliday(d Date) bool {
	return false
}

func (c *Calendar) NextWorkingDay(initialDate time.Time) DateWorkingTimes {
	resultDate := initialDate
	t := TimeOf(resultDate)

	var isInitDate = true

	for {
		d := DateOf(resultDate)
		if e, ok := c.exceptions.Get(d); ok {
			if e.WorkingDayTimes != nil && ((isInitDate && e.WorkingDayTimes.Workday.End.GreaterThan(t)) || !isInitDate) {
				return e
			}
		} else {
			switch resultDate.Weekday() {
			case time.Monday:
				if c.workingTimes.Monday == nil {
					break
				}
				if (isInitDate && c.workingTimes.Monday.Workday.End.GreaterThan(t)) || !isInitDate {
					return makeWorkingTimes(resultDate, c.workingTimes.Monday)
				}
			case time.Tuesday:
				if c.workingTimes.Tuesday == nil {
					break
				}
				if (isInitDate && c.workingTimes.Tuesday.Workday.End.GreaterThan(t)) || !isInitDate {
					return makeWorkingTimes(resultDate, c.workingTimes.Tuesday)
				}
			case time.Wednesday:
				if c.workingTimes.Wednesday == nil {
					break
				}
				if (isInitDate && c.workingTimes.Wednesday.Workday.End.GreaterThan(t)) || !isInitDate {
					return makeWorkingTimes(resultDate, c.workingTimes.Wednesday)
				}
			case time.Thursday:
				if c.workingTimes.Thursday == nil {
					break
				}
				if (isInitDate && c.workingTimes.Thursday.Workday.End.GreaterThan(t)) || !isInitDate {
					return makeWorkingTimes(resultDate, c.workingTimes.Thursday)
				}
			case time.Friday:
				if c.workingTimes.Friday == nil {
					break
				}
				if (isInitDate && c.workingTimes.Friday.Workday.End.GreaterThan(t)) || !isInitDate {
					return makeWorkingTimes(resultDate, c.workingTimes.Friday)
				}
			case time.Saturday:
				if c.workingTimes.Saturday == nil {
					break
				}
				if (isInitDate && c.workingTimes.Saturday.Workday.End.GreaterThan(t)) || !isInitDate {
					return makeWorkingTimes(resultDate, c.workingTimes.Saturday)
				}
			case time.Sunday:
				if c.workingTimes.Sunday == nil {
					break
				}
				if (isInitDate && c.workingTimes.Sunday.Workday.End.GreaterThan(t)) || !isInitDate {
					return makeWorkingTimes(resultDate, c.workingTimes.Sunday)
				}
			}
		}

		if isInitDate {
			isInitDate = false
		}
		resultDate = resultDate.Add(24 * time.Hour)
	}
}

func makeWorkingTimes(dt time.Time, wt *WorkingTimes) DateWorkingTimes {
	return DateWorkingTimes{
		Date:            DateOf(dt),
		WorkingDayTimes: wt,
	}
}

type DaytimePeriod struct {
	Start    Time
	End      Time
	Duration time.Duration
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

type DefaultWorkingTimes struct {
	Monday    *WorkingTimes
	Tuesday   *WorkingTimes
	Wednesday *WorkingTimes
	Thursday  *WorkingTimes
	Friday    *WorkingTimes
	Saturday  *WorkingTimes
	Sunday    *WorkingTimes
}

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
		from, to       Time
		fromSec, toSec int
		err            error
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

	fromSec = from.Hour*60*60 + from.Minute*60
	toSec = to.Hour*60*60 + to.Minute*60

	p := DaytimePeriod{
		Start:    from,
		End:      to,
		Duration: time.Duration(toSec-fromSec) * time.Second,
	}
	return p
}
