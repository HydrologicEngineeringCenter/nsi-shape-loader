package util

import "sort"

// StrContains checks whether a string slice contains a specific string
func StrContains(sl []string, s string) bool {
	i := sort.SearchStrings(sl, s)
	return i < len(sl) && sl[i] == s
}

// StrContains checks whether an int slice contains a specific int value
func IntContains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
