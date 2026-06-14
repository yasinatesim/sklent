package rag

import (
	"regexp"
	"strings"
)

var (
	slugReplacer = strings.NewReplacer(
		"ç", "c", "ğ", "g", "ı", "i", "ö", "o", "ş", "s", "ü", "u",
		"Ç", "c", "Ğ", "g", "İ", "i", "Ö", "o", "Ş", "s", "Ü", "u",
	)
	nonSlug = regexp.MustCompile(`[^a-z0-9]+`)
)

func Slugify(s string) string {
	s = strings.ToLower(slugReplacer.Replace(s))
	s = nonSlug.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
