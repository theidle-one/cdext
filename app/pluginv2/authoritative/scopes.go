package authoritative

import "sort"

type Scopes []Scope

func NewScopesFromSlices(scps []string) Scopes {
	return NewScopes(scps...)
}

func NewScopes(strScopes ...string) Scopes {
	s := make(Scopes, 0)

	if len(strScopes) == 0 || strScopes == nil {
		return s
	}

	for _, str := range strScopes {
		scp := NewScope(str)
		if scp == (Scope{}) {
			continue
		}

		s = append(s, scp)
	}
	return s
}

func (scopes Scopes) Merge(scps Scopes) (result Scopes) {
	indexes := make([]int, 0)
	unique := make(map[string]int)
	m := make(map[int]string)

	nextIndex := func() int {
		return len(m)
	}

	exist := func(str string) bool {
		_, ok := unique[str]
		return ok
	}

	for _, scp := range scopes {
		scpStr := scp.String()
		if exist(scpStr) {
			continue
		}

		nextIndex := nextIndex()
		m[nextIndex] = scpStr
		unique[scpStr] = nextIndex
		indexes = append(indexes, nextIndex)
	}

	for _, scp := range scps {
		scpStr := scp.String()
		if exist(scpStr) {
			continue
		}

		nextIndex := nextIndex()
		m[nextIndex] = scpStr
		unique[scpStr] = nextIndex
		indexes = append(indexes, nextIndex)
	}

	result = make(Scopes, 0)
	sort.Ints(indexes)
	for _, index := range indexes {
		result = append(result, NewScope(m[index]))
	}

	return

	//result = scopes
	//for _, scp := range scps {
	//	if !IsIn(scp, result) {
	//		result = append(result, scp)
	//	}
	//	continue
	//}
	//return
}

func (scopes Scopes) Strings() []string {
	var result []string
	for _, scp := range scopes {
		result = append(result, scp.String())
	}
	return result
}

func (scopes Scopes) Equals(scps Scopes) bool {
	if len(scopes) != len(scps) {
		return false
	}

	matched := 0
	for _, s1 := range scopes {
		for _, s2 := range scps {
			if s1.Equals(s2) {
				matched++
			}
		}
	}
	return matched == len(scopes)
}

func (scopes Scopes) Contains(scp Scope) bool {
	return scp.BelongsTo(scopes)
}

func (scopes Scopes) DeleteIndex(i int) Scopes {
	return append(scopes[:i], scopes[i+1:]...)
}

func (scopes Scopes) Delete(str string) (result Scopes) {
	for _, s := range scopes {
		if NewScope(str).Equals(s) {
			continue
		}
		result = append(result, s)
	}
	return
}

func IsIn(scp Scope, scps Scopes) bool {
	for _, scope := range scps {
		if scp.Equals(scope) {
			return true
		}
	}
	return false
}
