package authoritative

import (
	"fmt"
	"strings"
)

var Root = Scope{
	Action:  "*",
	Object:  "*",
	Subject: "*",
}

type Scope struct {
	Action  string
	Object  string
	Subject string
}

func (s1 Scope) Equals(s2 Scope) bool {
	return s1.Action == s2.Action && s1.Object == s2.Object && s1.Subject == s2.Subject
}

func (s Scope) String() string {
	if s == (Scope{}) {
		return ""
	}

	a := s.Action
	if a == "" {
		a = "*"
	}

	o := s.Object
	if o == "" {
		o = "*"
	}

	u := s.Subject
	if u == "" {
		u = "*"
	}

	return fmt.Sprintf("%s:%s:%s", a, o, u)
}

func NewScope(str string) Scope {
	if str == "" {
		return Scope{}
	}

	arr := strings.Split(str, ":")
	if len(arr) > 3 {
		return Scope{}
	}

	a := "*"
	if len(arr) >= 1 {
		a = arr[0]
	}

	o := "*"
	if len(arr) >= 2 {
		o = arr[1]
	}

	s := "*"
	if len(arr) == 3 {
		s = arr[2]
	}

	return Scope{Action: a, Object: o, Subject: s}
}

func (s Scope) IsSubsetOf(scp Scope) bool {
	if (scp.Object == "*" || scp.Object == s.Object) && (scp.Subject == "*" || scp.Subject == s.Subject) {
		return true
	}
	return false
}

func (s Scope) Contains(scp Scope) bool {
	var a, o, u bool

	if s.Action == scp.Action || s.Action == "*" {
		a = true
	}

	if s.Object == scp.Object || s.Object == "*" {
		o = true
	}

	if s.Subject == scp.Subject || s.Subject == "*" {
		u = true
	}

	return a && o && u
}

func (scp Scope) BelongsTo(scopes Scopes) bool {
	result := false
	for _, scope := range scopes {
		if (scope.Action == "*" || scp.Action == scope.Action) && (scope.Object == "*" || scp.Object == scope.Object) && (scope.Subject == "*" || scp.Subject == scope.Subject) {
			result = true
		}

		if (scope.Action == "-*" || "-"+scp.Action == scope.Action) && (scope.Object == "*" || scp.Object == scope.Object) && (scope.Subject == "*" || scp.Subject == scope.Subject) {
			result = false
		}
	}
	return result
}
