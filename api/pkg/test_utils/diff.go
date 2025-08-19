package test_utils

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func trimSpaces(s string) string {
	split := strings.Split(s, "\n")
	var result []string
	for i := range split {
		trimmed := strings.TrimSpace(split[i])
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\\n\n")
}

func Diff(got string, want string) string {
	differ := diffmatchpatch.New()
	// Prepare the input strings for comparison

	diff := differ.DiffMain(trimSpaces(got), trimSpaces(want), true)

	if len(diff) > 1 || diff[0].Type != diffmatchpatch.DiffEqual {
		return differ.DiffPrettyText(diff)
	} else {
		return ""
	}
}
