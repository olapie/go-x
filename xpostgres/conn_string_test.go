package xpostgres

import "testing"

func TestConnStringBuilder_Build(t *testing.T) {
	t.Run("unix_domain_socket", func(t *testing.T) {
		actual := NewConnStringBuilder().User("ola").UseUnixDomainSocket(true).Schema("public").Build()
		expected := "user=ola port=5432 search_path=public"
		if actual != expected {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	})

	t.Run("regular", func(t *testing.T) {
		actual := NewConnStringBuilder().
			Host("server1.olapie.com").
			Port(6543).
			Password("pass123").
			DB("testdb").User("user1").
			WithQuery("key1", "value1").
			WithQuery("key2", "value2").
			Build()
		expected := "postgres://user1:pass123@server1.olapie.com:6543/testdb?key1=value1&key2=value2"
		if actual != expected {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	})
}

func TestSetParameterInConnString(t *testing.T) {
	t.Run("unix_domain_socket", func(t *testing.T) {
		actual := NewConnStringBuilder().User("ola").UseUnixDomainSocket(true).Schema("public").Build()
		actual = SetParameterInConnString(actual, "search_path", "test")
		expected := "user=ola port=5432 search_path=test"
		if actual != expected {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	})

	t.Run("regular", func(t *testing.T) {
		actual := "postgres://user1:pass123@server1.olapie.com:6543/testdb?key1=value1&key2=value2"
		actual = SetParameterInConnString(actual, "key1", "value111")
		actual = SetParameterInConnString(actual, "search_path", "test")
		expected := "postgres://user1:pass123@server1.olapie.com:6543/testdb?key1=value111&key2=value2&search_path=test"
		if actual != expected {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	})
}
