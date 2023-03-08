package internal

import "regexp"

func Detect(line string) (bool, string) {
	pipeline := []*regexp.Regexp{
		regexp.MustCompile(`^\s*[\/]?[\/|#] (UNGEN: .*)$`),
		regexp.MustCompile(`^\s*\[\/\/\]: \# \'(UNGEN: .*)\'$`),
		regexp.MustCompile(`^\s*\/\* (UNGEN: .*)? \*\/\s*$`),
		regexp.MustCompile(`^\s*<!-- (UNGEN: .*)? -->\s*$`),
		regexp.MustCompile(`^\s*\{\/\* (UNGEN: .*)? \*\/\}\s*$`),
	}
	for _, regex := range pipeline {
		if regex.MatchString(line) {
			submatch := regex.FindStringSubmatch(line)
			return true, submatch[1]
		}
	}
	return false, ""
}
