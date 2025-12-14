package xtime

import (
	"fmt"
	"time"
)

func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func NumOfYearDays(year int) int {
	if IsLeap(year) {
		return 366
	}
	return 365
}

func GetDayTime(t time.Time) time.Duration {
	return time.Hour*time.Duration(t.Hour()) +
		time.Minute*time.Duration(t.Minute()) +
		time.Second*time.Duration(t.Second()) +
		time.Duration(t.Nanosecond())
}

func IsToday(t time.Time) bool {
	return DateWithTime(t).IsToday()
}

func IsYesterday(t time.Time) bool {
	return DateWithTime(t).IsYesterday()
}

func IsTomorrow(t time.Time) bool {
	return DateWithTime(t).IsTomorrow()
}

func BeginOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.Local)
}

func ToTimerText(seconds int64) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	if h < 100 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

func UnixSecP(sec *int64) *time.Time {
	if sec == nil {
		return nil
	}
	t := time.Unix(*sec, 0)
	return &t
}

func ToUnixSecP(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	sec := int64(t.Second())
	return &sec
}

const (
	Day  = 24 * time.Hour
	Week = 7 * 24 * time.Hour
)

var enWeekdaySymbols = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

func GetWeekdaySymbol(d int) string {
	d = d % 7
	return Localize(enWeekdaySymbols[d])
}

func abs[T ~int64](num T) T {
	if num < 0 {
		return -num
	}
	return num
}

func RepeatSchedule(repeat Repeat, start time.Time, before time.Time) []time.Time {
	var res []time.Time
	t := start
	for t.Before(before) {
		res = append(res, t)
		if repeat == Never {
			return res
		}
		t = RepeatNext(repeat, t)
	}
	return res
}

func RepeatNext(repeat Repeat, start time.Time) time.Time {
	switch repeat {
	case Daily:
		return start.AddDate(0, 0, 1)
	case Weekly:
		return start.AddDate(0, 0, 7)
	case BiWeekly:
		return start.AddDate(0, 0, 14)
	case Monthly:
		return start.AddDate(0, 1, 0)
	case Quarterly:
		return start.AddDate(0, 3, 0)
	case Yearly:
		return start.AddDate(1, 0, 0)
	default:
		panic(fmt.Sprintf("invalid repeat %d: %v", repeat, repeat))
	}
}
