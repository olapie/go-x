package xpostgres

import (
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapToHstore(m map[string]string) pgtype.Hstore {
	var h pgtype.Hstore
	for k, v := range m {
		h[k] = &v
	}
	return h
}

func HstoreToMap(h pgtype.Hstore) map[string]string {
	m := make(map[string]string, len(h))
	for k, v := range h {
		if v != nil {
			m[k] = *v
		}
	}
	return m
}

type hstoreScanner struct {
	m *map[string]string
}

var _ sql.Scanner = (*hstoreScanner)(nil)

func (hs *hstoreScanner) Scan(src any) error {
	if src == nil {
		return nil
	}
	var h pgtype.Hstore
	err := h.Scan(src)
	if err != nil {
		return err
	}
	m := HstoreToMap(h)
	*hs.m = m
	return nil
}
