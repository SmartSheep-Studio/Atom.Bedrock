package common

import (
	"fmt"
	"github.com/IGLOU-EU/go-wildcard"
)

func MatchTree(v map[string]any, requires ...string) error {
	for _, require := range requires {
		passed := false
		for key, val := range v {
			if wildcard.Match(key, require) && (val != nil || val != false) {
				passed = true
				break
			}
		}

		if !passed {
			return fmt.Errorf("missing node: %s", require)
		}
	}

	return nil
}

func MatchList(v []string, requires ...string) error {
	for _, require := range requires {
		passed := false
		for _, perm := range v {
			if wildcard.Match(perm, require) {
				passed = true
				break
			}
		}

		if !passed {
			return fmt.Errorf("missing item: %s", requires)
		}
	}

	return nil
}
