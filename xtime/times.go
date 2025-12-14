package xtime

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var timeFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-1-2",
	"20060102",
	"2006/1/2",
	"2/1/2006",
}

// ToDuration converts an interface to a time.Duration type.
func ToDuration[T string | int64](i T) (time.Duration, error) {
	switch v := any(i).(type) {
	case string:
		s := strings.ToLower(v)
		if strings.ContainsAny(s, "nsuµmh") {
			return time.ParseDuration(s)
		} else {
			return time.ParseDuration(s + "ns")
		}
	case int64:
		return time.Duration(v), nil
	default:
		return 0, fmt.Errorf("invalid type: %T", i)
	}
}

func ToTime[T ~string](s T) (time.Time, error) {
	return ToTimeInLocation(s, time.Local)
}

func ToTimeInLocation[T ~string](s T, loc *time.Location) (time.Time, error) {
	for _, df := range timeFormats {
		d, err := time.ParseInLocation(df, string(s), loc)
		if err == nil {
			return d, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot convert %#v to date", s)
}

func ToLocation(name string, offset int) *time.Location {
	// LoadLocation get failed on iOS
	loc, err := time.LoadLocation(name)
	if err == nil {
		return loc
	}
	log.Printf("Cannot load location %s: %v. Converted to a fixed zone", name, err)
	loc = time.FixedZone(name, offset)
	return loc
}

func IsDate(s string) bool {
	if len(s) > 10 {
		return false
	}
	_, err := ToTime(s)
	return err == nil
}

func ToDateString(s string) string {
	s = strings.Replace(s, "年", "-", 1)
	s = strings.Replace(s, "月", "-", 1)
	s = strings.Replace(s, "日", "", 1)
	s = strings.Replace(s, "o", "0", -1)
	s = strings.Replace(s, "O", "0", -1)
	s = strings.Replace(s, "l", "1", -1)
	s = strings.Replace(s, "I", "1", -1)
	if strings.Contains(s, "-") && len(s) >= 10 {
		s = s[:10]
	} else if len(s) >= 8 {
		s = s[:8]
	} else {
		return ""
	}

	t, err := ToTime(s)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}
