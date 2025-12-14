package xtime

//go:generate stringer -type Timespan,Repeat -trimprefix=Timespan -output=enum_string.go
type Timespan int

const (
	TimespanSecond Timespan = iota
	TimespanMinute
	TimespanHour
	TimespanDay
	TimespanWeek
	TimespanBiWeek
	TimespanMonth
	TimespanBiMonth
	TimespanQuarter
	TimespanHalfYear
	TimespanYear

	_TimespanCount
)

func (t Timespan) Valid() bool {
	return t >= TimespanSecond && t < _TimespanCount
}

type Repeat int

const (
	Never Repeat = iota
	Daily
	Weekly
	BiWeekly
	Monthly
	Quarterly
	Yearly

	_RepeatCount
)

func (r Repeat) Valid() bool {
	return r >= Never && r < _RepeatCount
}
