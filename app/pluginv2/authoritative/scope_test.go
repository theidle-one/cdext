package authoritative

import "testing"

func TestScopeEquals(t *testing.T) {
	s1 := Scope{Action: "read", Object: "policy", Subject: "*"}
	s2 := Scope{Action: "read", Object: "policy", Subject: "*"}
	if !s1.Equals(s2) {
		t.Error("scope should equals")
	}

	s1 = Scope{Action: "*", Object: "*", Subject: "*"}
	s2 = Scope{Action: "*", Object: "*", Subject: "*"}
	if !s1.Equals(s2) {
		t.Error("scope should equals")
	}

	s1 = Scope{Action: "read", Object: "policy", Subject: "1"}
	s2 = Scope{Action: "read", Object: "policy", Subject: "1"}
	if !s1.Equals(s2) {
		t.Error("scope should equals")
	}

	s1 = Scope{Action: "read", Object: "*", Subject: "*"}
	s2 = Scope{Action: "read", Object: "policy", Subject: "*"}
	if s1.Equals(s2) {
		t.Error("scope should not equals")
	}

	s1 = Scope{Action: "*", Object: "policy", Subject: "*"}
	s2 = Scope{Action: "*", Object: "*", Subject: "*"}
	if s1.Equals(s2) {
		t.Error("scope should not equals")
	}

	s1 = Scope{Action: "write", Object: "policy", Subject: "1"}
	s2 = Scope{Action: "read", Object: "policy", Subject: "1"}
	if s1.Equals(s2) {
		t.Error("scope should not equals")
	}
}

func TestScopeToString(t *testing.T) {
	s1 := Scope{Action: "update", Object: "policy", Subject: "1"}
	if s1.String() != "update:policy:1" {
		t.Error("scope to string wrong")
	}

	s1 = Scope{Action: "update", Object: "policy"}
	if s1.String() != "update:policy:*" {
		t.Error("scope to string wrong")
	}

	s1 = Scope{Action: "update"}
	if s1.String() != "update:*:*" {
		t.Error("scope to string wrong")
	}

	s1 = Scope{}
	if s1.String() != "" {
		t.Error("scope to string wrong")
	}
}

func TestNewScopeFromString(t *testing.T) {
	s := NewScope("")
	if s != (Scope{}) {
		t.Error("invalid new scope")
	}

	s = NewScope("*")
	if s.Action != "*" || s.Object != "*" || s.Subject != "*" {
		t.Error("invalid new scope")
	}

	s = NewScope("update:*")
	if s.Action != "update" || s.Object != "*" || s.Subject != "*" {
		t.Error("invalid new scope")
	}

	s = NewScope("update:policy")
	if s.Action != "update" || s.Object != "policy" || s.Subject != "*" {
		t.Error("invalid new scope")
	}

	s = NewScope("update:policy:1")
	if s.Action != "update" || s.Object != "policy" || s.Subject != "1" {
		t.Error("invalid new scope")
	}
}

func TestScopeContains(t *testing.T) {
	s1 := Scope{Action: "*", Object: "*", Subject: "*"}
	s2 := Scope{Action: "read", Object: "policy", Subject: "1"}
	if !s1.Contains(s2) {
		t.Error("scope should contains")
	}

	s1 = Scope{Action: "*", Object: "policy", Subject: "*"}
	s2 = Scope{Action: "read", Object: "policy", Subject: "1"}
	if !s1.Contains(s2) {
		t.Error("scope should contains")
	}

	s1 = Scope{Action: "read", Object: "*", Subject: "*"}
	s2 = Scope{Action: "read", Object: "policy", Subject: "1"}
	if !s1.Contains(s2) {
		t.Error("scope should contains")
	}

	s1 = Scope{Action: "read", Object: "policy", Subject: "1"}
	s2 = Scope{Action: "read", Object: "policy", Subject: "1"}
	if !s1.Contains(s2) {
		t.Error("scope should contains")
	}

	s1 = Scope{Action: "*", Object: "policy", Subject: "*"}
	s2 = Scope{Action: "*", Object: "activity", Subject: "*"}
	if s1.Contains(s2) {
		t.Error("scope should not contains")
	}

	s1 = Scope{Action: "read", Object: "policy", Subject: "*"}
	s2 = Scope{Action: "*", Object: "policy", Subject: "*"}
	if s1.Contains(s2) {
		t.Error("scope should not contains")
	}
}

func TestScopeBelongsTo(t *testing.T) {
	scope := Scope{Action: "*", Object: "policy", Subject: "*"}
	scopes := Scopes{
		{Action: "*", Object: "*", Subject: "*"},
	}
	if !scope.BelongsTo(scopes) {
		t.Error("scope should belongs to")
	}

	scope = Scope{Action: "read", Object: "policy", Subject: "*"}
	scopes = Scopes{
		{Action: "read", Object: "*", Subject: "*"},
	}
	if !scope.BelongsTo(scopes) {
		t.Error("scope should belongs to")
	}

	scope = Scope{Action: "read", Object: "licence", Subject: "*"}
	scopes = Scopes{
		{Action: "*", Object: "*", Subject: "*"},
		{Action: "-*", Object: "licence", Subject: "*"},
	}
	if scope.BelongsTo(scopes) {
		t.Error("scope should not belongs to")
	}

	scope = Scope{Action: "read", Object: "licence", Subject: "1"}
	scopes = Scopes{
		{Action: "*", Object: "*", Subject: "*"},
		{Action: "-*", Object: "licence", Subject: "*"},
		{Action: "-*", Object: "system-admin", Subject: "*"},
	}
	if scope.BelongsTo(scopes) {
		t.Error("scope should not belongs to")
	}
}
