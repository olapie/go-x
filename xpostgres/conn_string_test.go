package xpostgres

import "testing"

func TestSetParameterInConnString(t *testing.T) {
	s1 := NewConnStringBuilder().User("ola").UseUnixDomainSocket(true).Schema("public").Build()
	t.Log(s1)
	s2 := SetParameterInConnString(s1, "search_path", "test")
	if s2 != "host=/tmp port=5432 user=ola search_path=test" {
		t.Fatal(s2)
	}

	s3 := "host=/tmp port=5432 user=ola search_path="
	s4 := SetParameterInConnString(s3, "search_path", "test")
	if s4 != "host=/tmp port=5432 user=ola search_path=test" {
		t.Fatal(s4)
	}

	s5 := NewConnStringBuilder().User("ola").Schema("public").Build()
	t.Log(s5)
	s6 := SetParameterInConnString(s3, "search_path", "test")
	if s6 != "host=/tmp port=5432 user=ola search_path=test" {
		t.Fatal(s6)
	}
}
