package xconv

import "time"

func Int64ToTimeP(sec int64) *time.Time {
	if sec == 0 {
		return nil
	}
	t := time.Unix(sec, 0)
	return &t
}

func Int64PToTimeP(sec *int64) *time.Time {
	if sec == nil {
		return nil
	}
	t := time.Unix(*sec, 0)
	return &t
}

func TimePToInt64P(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	i := t.Unix()
	return &i
}

func TimePToInt64(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.Unix()
}
