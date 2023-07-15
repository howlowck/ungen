package internal

import "regexp"

type DetectResult int

const (
	NotDetected DetectResult = iota
	DetectedDefault
	DetectedGpt
)

func Detect(line string) (DetectResult, string) {
	defaultPipeline := []*regexp.Regexp{
		regexp.MustCompile(`^\s*[\/]?[\/|#] (UNGEN: .*)$`),
		regexp.MustCompile(`^\s*\[\/\/\]: \# \'(UNGEN: .*)\'$`),
		regexp.MustCompile(`^\s*\/\* (UNGEN: .*)? \*\/\s*$`),
		regexp.MustCompile(`^\s*<!-- (UNGEN: .*)? -->\s*$`),
		regexp.MustCompile(`^\s*\{\/\* (UNGEN: .*)? \*\/\}\s*$`),
	}
	for _, regex := range defaultPipeline {
		if regex.MatchString(line) {
			submatch := regex.FindStringSubmatch(line)
			return DetectedDefault, submatch[1]
		}
	}
	return NotDetected, ""
}
