package xurl

import "testing"

func TestSetQuery(t *testing.T) {
	origin := "postgres:///ola?host=/var/run/postgresql/&search_path=default"
	got, err := SetQuery(origin, "search_path", "schema_a")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(got)
}
