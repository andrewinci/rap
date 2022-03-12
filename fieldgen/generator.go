package fieldgen

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// generic data generator
type Gen func(seed int64) interface{}

func NewGenerator(rawPattern string) Gen {
	pattern := parsePattern(rawPattern)
	fmt.Println(pattern)
	return func(seed int64) interface{} {
		r := rand.New(rand.NewSource(seed))
		var res string
		for _, c := range pattern.content {
			for i := 0; i < c.count; i++ {
				j := r.Intn(len(c.options))
				res += c.options[j]
			}
		}
		fmt.Println(res)
		switch pattern.type_ {
		case "string":
			return res
		case "int":
			number, _ := strconv.Atoi(res)
			return number
		default:
			return nil
		}
	}
}

type patternOption struct {
	options []string
	count   int
}

type pattern struct {
	type_   string
	content []patternOption
}

// parse the pattern and returns the type and the content
func parsePatternType(p string) (string, string, bool) {
	var re = regexp.MustCompile(`(?m)\{(boolean|int|long|float|double|bytes|string)\}(.*)`)
	var matches = re.FindAllStringSubmatch(p, -1)
	if matches == nil || len(matches) != 1 {
		// invalid pattern
		return "", "", false
	}
	return matches[0][1], matches[0][2], true
}

func validatePatternContent(rawContent string) bool {
	// validate the content
	var re = regexp.MustCompile(`^(\[([^\]]+)\]\{(\d+)\})+$`)
	var matches = re.FindAllStringSubmatch(rawContent, -1)
	return matches != nil
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
	fmt.Printf("Pattern type %s\n", patternType)

	if !validatePatternContent(rawContent) {
		return nil
	}

	var re = regexp.MustCompile(`\[([^\]]+)\]\{(\d+)\}`)
	var matches = re.FindAllStringSubmatch(rawContent, -1)
	content := []patternOption{}
	if matches == nil {
		// invalid pattern
		return nil
	}
	for _, c := range matches {
		count, _ := strconv.Atoi(c[2])
		var options []string
		for _, c := range strings.Split(c[1], "|") {
			switch c {
			case "a-z":
				options = append(options, getLowerCaseLetters()...)
			case "A-Z":
				options = append(options, getUpperCaseLetters()...)
			case "0-9":
				options = append(options, getDigits()...)
			case "a-Z":
				options = append(getUpperCaseLetters(), getLowerCaseLetters()...)
			default:
				// constant case
				options = append(options, c)
			}
		}

		content = append(content, patternOption{options: options, count: count})
	}

	return &pattern{
		type_:   patternType,
		content: content,
	}
}
