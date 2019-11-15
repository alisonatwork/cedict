package pinyin

import (
	"regexp"
	"strconv"
	"strings"
)

var lookup = map[string][]string{
	"a": {"ā", "á", "ǎ", "à"},
	"e": {"ē", "é", "ě", "è"},
	"i": {"ī", "í", "ǐ", "ì"},
	"o": {"ō", "ó", "ǒ", "ò"},
	"u": {"ū", "ú", "ǔ", "ù"},
	"ü": {"ǖ", "ǘ", "ǚ", "ǜ"},
	"A": {"Ā", "Á", "Ǎ", "À"},
	"E": {"Ē", "É", "Ě", "È"},
	"I": {"Ī", "Í", "Ǐ", "Ì"},
	"O": {"Ō", "Ó", "Ǒ", "Ò"},
	"U": {"Ū", "Ú", "Ǔ", "Ù"},
	"Ü": {"Ǖ", "Ǘ", "Ǚ", "Ǜ"},
}

var pattern = regexp.MustCompile(`(([aeiouüv]|[uU]:){1,3})(n?g?r?)([12345])`)
var replacer = strings.NewReplacer("v", "ü", "u:", "ü", "V", "Ü", "U:", "Ü")

// NumberToMark replaces pinyin tone numbers with a mark, e.g. de2 -> dé
func NumberToMark(in string) string {
	return replaceAllStringSubmatchFunc(pattern, in, func(match []string) string {
		vowels := replacer.Replace(match[1])
		suffix := match[3]
		tone, err := strconv.Atoi(match[4])
		if err != nil || tone == 5 {
			return vowels + suffix
		}
		var vowel string
		if len([]rune(vowels)) == 1 || strings.Contains("aeoAEO", string([]rune(vowels)[0])) {
			vowel = string([]rune(vowels)[0])
		} else {
			vowel = string([]rune(vowels)[1])
		}
		return strings.Replace(vowels, vowel, lookup[vowel][tone-1], 1) + suffix
	})
}
