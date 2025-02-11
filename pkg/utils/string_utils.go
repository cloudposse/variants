package utils

import (
	"encoding/csv"
	"strings"
)

// UniqueStrings returns a unique subset of the string slice provided
func UniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// SplitStringByDelimiter splits a string by the delimiter, not splitting inside quotes
func SplitStringByDelimiter(str string, delimiter rune) ([]string, error) {
	r := csv.NewReader(strings.NewReader(str))
	r.Comma = delimiter

	parts, err := r.Read()
	if err != nil {
		return nil, err
	}

	return parts, nil
}

// SplitStringAtFirstOccurrence splits a string into two parts at the first occurrence of the separator
func SplitStringAtFirstOccurrence(s string, sep byte) [2]string {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return [2]string{s[:i], s[i+1:]}
		}
	}
	return [2]string{s, ""}
}
