package authoritative

import (
	"testing"
)

func TestNewScopesFromString(t *testing.T) {
	scopes := NewScopes("*", "*:policy", "read:policy", "read:policy:1")

	if len(scopes) != 4 {
		t.Error("new scopes has wrong length")
	}

	if scopes[0].Action != "*" || scopes[0].Object != "*" || scopes[0].Subject != "*" {
		t.Error("new scopes create invalid")
	}

	if scopes[1].Action != "*" || scopes[1].Object != "policy" || scopes[1].Subject != "*" {
		t.Error("new scopes create invalid")
	}

	if scopes[2].Action != "read" || scopes[2].Object != "policy" || scopes[2].Subject != "*" {
		t.Error("new scopes create invalid")
	}

	if scopes[3].Action != "read" || scopes[3].Object != "policy" || scopes[3].Subject != "1" {
		t.Error("new scopes create invalid")
	}
}

func TestScopesEquals(t *testing.T) {
	scopes1 := NewScopes("*", "*:policy", "read:policy", "read:policy:1")
	scopes2 := NewScopes("*", "*:policy", "read:policy", "read:policy:1")
	if !scopes1.Equals(scopes2) {
		t.Error("scopes should be equal")
	}

	scopes1 = NewScopes("*:policy", "read:policy", "read:policy:1")
	scopes2 = NewScopes("*", "*:policy", "read:policy", "read:policy:1")
	if scopes1.Equals(scopes2) {
		t.Error("scopes should not be equal")
	}

	scopes1 = NewScopes("*", "*:policy", "write:policy", "read:policy:1")
	scopes2 = NewScopes("*", "*:policy", "read:policy", "read:policy:1")
	if scopes1.Equals(scopes2) {
		t.Error("scopes should not be equal")
	}
}

func TestScopesContains(t *testing.T) {
	scopes := NewScopes("*")
	if !scopes.Contains(NewScope("*")) {
		t.Error("scopes should contain")
	}

	if !scopes.Contains(NewScope("*:policy:1")) {
		t.Error("scopes should contain")
	}

	if !scopes.Contains(NewScope("read:policy:1")) {
		t.Error("scopes should contain")
	}

	scopes = NewScopes("*", "-*:licence")
	if !scopes.Contains(NewScope("read:policy:1")) {
		t.Error("scopes should contain")
	}

	scopes = NewScopes("-*:licence", "*")
	if !scopes.Contains(NewScope("read:licence")) {
		t.Error("scopes should not contain")
	}

	scopes = NewScopes("*:g", "-*:g", "*:g")
	if !scopes.Contains(NewScope("a:g")) {
		t.Error("scopes should contain")
	}

	scopes = NewScopes("a:b", "-a:c", "d:e", "d:f", "-a:g", "*:g")
	if !scopes.Contains(NewScope("a:g")) {
		t.Error("scopes should contain")
	}

	scopes = NewScopes("-do:b", "do:*", "*:b", "-*:a", "do:b")
	if scopes.Contains(NewScope("do:a")) {
		t.Error("scopes should contain")
	}

	if !scopes.Contains(NewScope("do:b")) {
		t.Error("scopes should contain")
	}

	if !scopes.Contains(NewScope("do:c")) {
		t.Error("scopes should contain")
	}

	scopes = NewScopes("do:a", "-do:b", "do:*", "*:b", "-*:a", "x:a")
	if !scopes.Contains(NewScope("do:b")) {
		t.Error("scopes should contain")
	}

	if !scopes.Contains(NewScope("x:b")) {
		t.Error("scopes should contain")
	}

	if scopes.Contains(NewScope("do:a")) {
		t.Error("scopes should contain")
	}

	if !scopes.Contains(NewScope("x:a")) {
		t.Error("scopes should contain")
	}

	scopes = NewScopes("*:*:*", "-*:malware:sandbox", "-*:category:*", "read:category:*", "-*:category-resource:*", "*:category-resource:blacklist", "*:category-resource:whitelist")
	if !scopes.Contains(NewScope("read:category-resource:blacklist")) {
		t.Error("scopes should contain")
	}
}

func TestMerge(t *testing.T) {
	s1 := NewScopes("*:*:*", "-*:malware:sandbox", "-*:category:*", "read:category:*", "-*:category-resource:*", "*:category-resource:blacklist", "*:category-resource:whitelist")
	s2 := NewScopes("*:*:*", "-*:user:*", "-*:groups:*", "-*:ssl:*")
	s3 := NewScopes("*:*:*", "-*:licence", "-*:system-admin")

	merged := s1.Merge(s2).Merge(s3)
	if !merged.Equals(NewScopes(
		"*:*:*", "-*:malware:sandbox", "-*:category:*", "read:category:*", "-*:category-resource:*", "*:category-resource:blacklist", "*:category-resource:whitelist",
		"-*:user:*", "-*:groups:*", "-*:ssl:*",
		"-*:licence", "-*:system-admin",
	)) {
		t.Error("scopes should equals")
	}
}
