package xtime

import (
	"testing"
)

func TestMonth_NumOfWeeks(t *testing.T) {
	tests := []struct {
		Month      *Month
		NumOfWeeks int
	}{
		{NewMonth(2020, 6),
			5,
		}, {NewMonth(2020, 7),
			5,
		}, {NewMonth(2020, 8),
			6,
		}, {NewMonth(2020, 9),
			5,
		},
	}

	for _, te := range tests {
		if te.NumOfWeeks != te.Month.NumOfWeeks() {
			t.Fatal(te.NumOfWeeks, te.Month.NumOfWeeks(), te.Month)
		}
	}
}

func TestMonth_GetCalendarDate(t *testing.T) {
	tests := []struct {
		Month *Month
		Week  int
		Day   int
		Date  *Date
	}{
		{NewMonth(2020, 9),
			1,
			3,
			NewDate(2020, 9, 1),
		}, {NewMonth(2020, 9),
			4,
			1,
			NewDate(2020, 9, 20),
		}, {NewMonth(2020, 9),
			5,
			4,
			NewDate(2020, 9, 30),
		}, {NewMonth(2020, 9),
			1,
			1,
			nil,
		}, {NewMonth(2020, 9),
			5,
			7,
			nil,
		},
	}

	for _, te := range tests {
		date := te.Month.GetCalendarDate(te.Week, te.Day)
		if date == te.Date {
			continue
		}
		if date == nil {
			t.Fatal(te.Month.String(), te.Week, te.Day)
		}
		if te.Date == nil {
			t.FailNow()
		}

		if !date.Equals(te.Date) {
			t.Fatal(date.String(), te.Date.String())
		}
	}
}
