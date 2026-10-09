package config

import "testing"

func TestReadOnlyAllows(t *testing.T) {
	defer SetReadOnly(false, false)

	cases := []struct {
		on, postIsQuery bool
		method          string
		want            bool
	}{
		{false, false, "DELETE", true},
		{true, false, "GET", true},
		{true, false, "POST", false},
		{true, true, "POST", true},
		{true, true, "PUT", false},
		{true, true, "PATCH", false},
		{true, true, "DELETE", false},
		{false, true, "POST", true},
	}
	for _, c := range cases {
		SetReadOnly(c.on, c.postIsQuery)
		if got := ReadOnlyAllows(c.method); got != c.want {
			t.Errorf("readOnly=%v postIsQuery=%v %s: got %v, want %v", c.on, c.postIsQuery, c.method, got, c.want)
		}
	}
}
