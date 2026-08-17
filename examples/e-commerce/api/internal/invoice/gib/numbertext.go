package gib

import (
	"strings"
	"unicode"
)

var ONES_WORDS = []string{"", "BİR", "İKİ", "ÜÇ", "DÖRT", "BEŞ", "ALTI", "YEDİ", "SEKİZ", "DOKUZ"}
var TENS_WORDS = []string{"", "ON", "YİRMİ", "OTUZ", "KIRK", "ELLİ", "ALTMIŞ", "YETMİŞ", "SEKSEN", "DOKSAN"}

func convertHundreds(n int) string {
	var parts []string
	h := n / 100
	t := (n % 100) / 10
	o := n % 10
	if h == 1 {
		parts = append(parts, "YÜZ")
	} else if h > 1 {
		parts = append(parts, ONES_WORDS[h]+" YÜZ")
	}
	if t > 0 {
		parts = append(parts, TENS_WORDS[t])
	}
	if o > 0 {
		parts = append(parts, ONES_WORDS[o])
	}
	return strings.Join(parts, " ")
}

func numberToTurkishWords(n int) string {
	if n == 0 {
		return "SIFIR"
	}
	var parts []string
	milyar := n / 1_000_000_000
	milyon := (n % 1_000_000_000) / 1_000_000
	bin := (n % 1_000_000) / 1_000
	rest := n % 1_000
	if milyar > 0 {
		parts = append(parts, convertHundreds(milyar)+" MİLYAR")
	}
	if milyon > 0 {
		parts = append(parts, convertHundreds(milyon)+" MİLYON")
	}
	if bin == 1 {
		parts = append(parts, "BİN")
	} else if bin > 1 {
		parts = append(parts, convertHundreds(bin)+" BİN")
	}
	if rest > 0 {
		parts = append(parts, convertHundreds(rest))
	}
	return strings.Join(parts, " ")
}

// Turkish keeps the dot on İ/i and drops it on I/ı; Go's default casing maps both to ASCII I,
// which would print "IkiBin" instead of "İkiBin" on a legally-printed invoice amount.
var turkishLower = strings.NewReplacer("I", "ı", "İ", "i")

func titleCaseWord(word string) string {
	if word == "" {
		return word
	}
	lowered := []rune(strings.ToLower(turkishLower.Replace(word)))
	first := lowered[0]
	switch first {
	case 'i':
		first = 'İ'
	case 'ı':
		first = 'I'
	default:
		first = unicode.ToUpper(first)
	}
	return string(first) + string(lowered[1:])
}

// convertPriceToText renders a TRY amount as the Turkish invoice note GIB requires (e.g. "YALNIZ YüzTürkLirası SıfırKr'dir.").
func ConvertPriceToText(price float64) string {
	mainPart := int(price)
	kurusPart := int((price-float64(mainPart))*100 + 0.5)

	mainWords := joinTitleCased(numberToTurkishWords(mainPart))
	var kurusWords string
	if kurusPart == 0 {
		kurusWords = "SIFIR KR'DIR."
	} else {
		kurusWords = joinTitleCased(numberToTurkishWords(kurusPart)) + "KR'DIR."
	}
	return "YALNIZ " + mainWords + "TürkLirası " + kurusWords
}

func joinTitleCased(words string) string {
	fields := strings.Fields(words)
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		out = append(out, titleCaseWord(f))
	}
	return strings.Join(out, "")
}
