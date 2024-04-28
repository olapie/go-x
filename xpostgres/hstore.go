package xpostgres

import (
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"go.olapie.com/x/xconv"
)

func MapToHstore[K, V ~string](m map[K]V) pgtype.Hstore {
	h := make(pgtype.Hstore, len(m))
	for k, v := range m {
		h[string(k)] = xconv.Pointer(string(v))
	}
	return h
}

func HstoreToMap[K, V ~string](h pgtype.Hstore) map[K]V {
	m := make(map[K]V, len(h))
	for k, v := range h {
		if v != nil {
			m[K(k)] = V(*v)
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
	m := HstoreToMap[string, string](h)
	*hs.m = m
	return nil
}
