package avrogen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hamba/avro"
)

type patternOption struct {
	options []func() []string
	count   int
}

type pattern struct {
	type_   string
	content []patternOption
}

func getUpperCaseLetters() []string {
	return []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
}
func getLowerCaseLetters() []string {
	return []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
}
func getDigits() []string {
	return []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
}

func parsePattern(p string) *pattern {
	patternType, rawContent, ok := parsePatternType(p)
	if !ok {
		return nil
	}

	var re = regexp.MustCompile(`\[([^\]]+)\]\{(\d+)\}`)
	var matches = re.FindAllStringSubmatch(rawContent, -1)
	content := []patternOption{}
	for _, c := range matches {
		count, _ := strconv.Atoi(c[2])
		var options []func() []string
		for _, c := range strings.Split(c[1], "|") {
			switch strings.Trim(c, " ") {
			case "a-z":
				options = append(options, getLowerCaseLetters)
			case "A-Z":
				options = append(options, getUpperCaseLetters)
			case "0-9":
				options = append(options, getDigits)
			case "a-Z":
				options = append(options, getUpperCaseLetters, getLowerCaseLetters)
			case "uuid()":
				options = append(options, func() []string { return []string{uuid.New().String()} })
			case "timestamp_ms()":
				options = append(options, func() []string {
					return []string{
						fmt.Sprintf("%d", time.Now().UnixMilli())}
				})
			default:
				// constant case
				// important: need to reassign the c into s
				// otherwise the closure `func` will alway point
				// to the latest c value
				s := strings.Trim(c, " ")
				options = append(options, func() []string { return []string{s} })
			}
		}

		content = append(content, patternOption{options: options, count: count})
	}

	return &pattern{
		type_:   patternType,
		content: content,
	}
}

// parse the pattern and returns the type and the content
func parsePatternType(p string) (string, string, bool) {
	types := strings.Join([]string{
		string(avro.Boolean),
		string(avro.Int),
		string(avro.Long),
		string(avro.Float),
		string(avro.Double),
		string(avro.String)}, "|")
	regex := fmt.Sprintf(`^\{(%s)\}((\[([^\]]+)\]\{(\d+)\})+)$`, types)
	var re = regexp.MustCompile(regex)
	var matches = re.FindAllStringSubmatch(p, -1)
	if matches == nil || len(matches) != 1 {
		// invalid pattern
		return "", "", false
	}
	return matches[0][1], matches[0][2], true
}
