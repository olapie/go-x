package xmobile

import (
	"fmt"

	"go.olapie.com/times"
)

func GetDateTimeString(t int64) string {
	tm := times.TimeWithUnix(t)
	return tm.Date().PrettyText() + " " + tm.TimeTextWithZero()
}

func GetRelativeDateTimeString(t int64) string {
	tm := times.TimeWithUnix(t)
	return tm.RelativeDateTimeText()
}

type Time = times.Time

func NowTime() *Time {
	return (*Time)(times.NewTime())
}

func TimeWithUnix(seconds int64) *Time {
	return times.TimeWithUnix(seconds)
}

func TimerText(elapse int64) string {
	h := elapse / 3600
	elapse %= 3600
	m := elapse / 60
	s := elapse % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
