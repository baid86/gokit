package utils

import "strings"

// StringArrayContains find in string array
//
// This routine looks for a word in string array. Case sensitive and insensitive
// search can be done
//
func StringArrayContains(s []string, e string, casecare bool) bool {
	if !casecare {
		e = strings.ToLower(e)
	}
	for _, a := range s {
		if !casecare {
			a = strings.ToLower(a)
		}
		if a == e {
			return true
		}
	}
	return false
}
