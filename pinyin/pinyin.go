package pinyin

import (
	"regexp"
	"strconv"
	"strings"
)

var lookup = map[string][]string{
	"a": []string{"ā", "á", "ǎ", "à"},
	"e": []string{"ē", "é", "ě", "è"},
	"i": []string{"ī", "í", "ǐ", "ì"},
	"o": []string{"ō", "ó", "ǒ", "ò"},
	"u": []string{"ū", "ú", "ǔ", "ù"},
	"ü": []string{"ǖ", "ǘ", "ǚ", "ǜ"},
	"A": []string{"Ā", "Á", "Ǎ", "À"},
	"E": []string{"Ē", "É", "Ě", "È"},
	"I": []string{"Ī", "Í", "Ǐ", "Ì"},
	"O": []string{"Ō", "Ó", "Ǒ", "Ò"},
	"U": []string{"Ū", "Ú", "Ǔ", "Ù"},
	"Ü": []string{"Ǖ", "Ǘ", "Ǚ", "Ǜ"},
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
