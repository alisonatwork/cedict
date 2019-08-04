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

// replaceAllStringSubmatchFunc calls the replacement function passing all matched groups,
// similar to the behavior of JavaScript's String.replace(pattern, function)
// Issue: https://github.com/golang/go/issues/5690
// Workaround: https://gist.github.com/elliotchance/d419395aa776d632d897
// Blog: https://medium.com/@elliotchance/go-replace-string-with-regular-expression-callback-f89948bad0bb
func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

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
