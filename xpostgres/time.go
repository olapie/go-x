package xpostgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func FromDBTime(t pgtype.Timestamp) *time.Time {
	if t.Valid {
		v := t.Time.UTC()
		return &v
	}
	return nil
}

func ToDBTime(t *time.Time) pgtype.Timestamp {
	pt := pgtype.Timestamp{}
	if t == nil {
		pt.Valid = false
		return pt
	}
	pt.Valid = true
	pt.Time = *t
	return pt
}

func FromDBTimeTz(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		v := t.Time.UTC()
		return &v
	}
	return nil
}

func ToDBTimeTz(t *time.Time) pgtype.Timestamptz {
	pt := pgtype.Timestamptz{}
	if t == nil {
		pt.Valid = false
		return pt
	}
	pt.Valid = true
	pt.Time = *t
	return pt
}
