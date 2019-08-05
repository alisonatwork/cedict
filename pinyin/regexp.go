package pinyin

import "regexp"

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
