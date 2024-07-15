package access

import (
	"errors"

	"github.com/layou233/zbproxy/v3/common/set"
)

var ErrRejected = errors.New("interrupted by access control")

// Check checks if item passes the access control.
func Check(lists []set.StringSet, mode string, item string) (hit bool) {
	for _, list := range lists {
		if hit = list.Has(item); hit {
			break
		}
	}
	switch mode {
	case AllowMode:
		if !hit {
			return false
		}
		return true
	case BlockMode:
		if hit {
			return false
		}
		return true
	}
	panic("bad access mode")
}
