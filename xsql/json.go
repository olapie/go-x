package xsql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"go.olapie.com/x/xconv"
)

func JSON(v any) any {
	if v == nil {
		return nil
	}
	switch val := reflect.ValueOf(v); val.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Array:
		if val.IsNil() {
			return nil
		}
	default:
		break
	}
	return &jsonHolder{v: v}
}

type jsonHolder struct {
	v any
}

var _ driver.Valuer = (*jsonHolder)(nil)
var _ sql.Scanner = (*jsonHolder)(nil)

func (j *jsonHolder) Scan(src any) error {
	if src == nil {
		return nil
	}

	b, err := xconv.ToBytes(src)
	if err != nil {
		return fmt.Errorf("parse bytes: %w", err)
	}

	if len(b) == 0 {
		return nil
	}

	err = json.Unmarshal(b, j.v)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

func (j *jsonHolder) Value() (driver.Value, error) {
	return json.Marshal(j.v)
}
