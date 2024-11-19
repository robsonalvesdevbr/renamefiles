package renamer

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// normalizeUnicode removes diacritics and converts styled characters to their basic form.
func normalizeUnicode(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, norm.NFKD.String(input))
}

// toTitleCaseWithSeparator converts a string to Title Case while maintaining separators.
func toTitleCaseWithSeparator(input, separator string) string {
	if separator == "" {
		separator = " "
	}
	words := strings.Split(input, separator)
	for i, word := range words {
		if len(word) > 0 {
			if len(word) > 1 {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			} else {
				words[i] = strings.ToUpper(word)
			}
		}
	}
	return strings.Join(words, separator)
}
