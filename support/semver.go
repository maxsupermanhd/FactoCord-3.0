package support

import (
	"strconv"
	"strings"
)

type SemanticVersionT struct {
	Full  string
	Split []uint
}

func SemanticVersion(s string) (*SemanticVersionT, *error) {
	if s == "" {
		return &SemanticVersionT{}, nil
	}
	res := SemanticVersionT{
		Full: s,
	}
	for _, version := range strings.Split(s, ".") {
		num, err := strconv.ParseUint(version, 10, 0)
		if err != nil {
			return nil, &err
		}
		res.Split = append(res.Split, uint(num))
	}
	return &res, nil
}

func (v *SemanticVersionT) Equal(v2 *SemanticVersionT) bool {
	return v.Full == "" || v2.Full == "" || v.Full == v2.Full
}

func (v *SemanticVersionT) NewerThan(v2 *SemanticVersionT) bool {
	for z := 0; z < len(v.Split); z++ {
		if z >= len(v2.Split) {
			return true // 0.1.1 > 0.1
		}
		if v.Split[z] != v2.Split[z] {
			return v.Split[z] > v2.Split[z] // 1.2 > 1.1
		}
	}
	return false // 0.1 < 0.1.1
}
