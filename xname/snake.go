package naming

import (
	"fmt"
	"unicode"
)

func ToSnake(s string, options ...Option) string {
	input := []rune(s)
	n := len(input)
	output := make([]rune, 0, n+2)
	type State int
	const (
		Start State = iota
		FirstUpper
		SecondUpper
		Lower
		WordEnd
		UpperWord
		UpperWordEnd
	)

	useAcronym := false
	if options != nil {
		var opts Options
		for _, fn := range options {
			fn(&opts)
		}
		useAcronym = opts.useAcronym
	}

	state := Start
	wordStart := -1
	for i := 0; i < n; i++ {
		r := input[i]
		switch state {
		case Start:
			if isDelimiter(r) {
				break
			}
			wordStart = i
			if unicode.IsUpper(r) {
				state = FirstUpper
			} else {
				state = Lower
			}
		case FirstUpper:
			if isDelimiter(r) {
				state = WordEnd
				output = appendSnakeWord(output, input[wordStart:i])
				state = Start
				wordStart = -1
				break
			}
			if unicode.IsUpper(r) {
				state = SecondUpper
			} else {
				state = Lower
			}
		case SecondUpper:
			if isDelimiter(r) {
				state = WordEnd
				output = appendSnakeWord(output, input[wordStart:i])
				state = Start
				wordStart = -1
				break
			}
			if unicode.IsUpper(r) {
				state = UpperWord
			} else {
				state = Lower
			}
		case Lower:
			if isDelimiter(r) {
				output = appendSnakeWord(output, input[wordStart:i])
				state = Start
				wordStart = -1
				break
			}

			if unicode.IsUpper(r) {
				output = appendSnakeWord(output, input[wordStart:i])
				i--
				state = Start
				wordStart = -1
				break
			}
		case UpperWord:
			if isDelimiter(r) {
				output = appendSnakeWord(output, input[wordStart:i])
				state = Start
				wordStart = -1
				break
			}

			if unicode.IsUpper(r) {
				break
			}

			state = UpperWordEnd
			if !useAcronym {
				state = Lower
				break
			}
			acronymsMapRWMutex.RLock()
			acronym := acronymsU2L[string(input[wordStart:i-1])]
			acronymsMapRWMutex.RUnlock()
			if acronym != "" {
				output = appendSnakeWord(output, input[wordStart:i-1])
				i -= 2
				state = Start
				wordStart = -1
				break
			}
			state = Lower
		default:
			panic(fmt.Sprintf("invalid state: %d", state))
		}
	}

	if wordStart >= 0 {
		output = appendSnakeWord(output, input[wordStart:])
	}
	return string(output)
}

func appendSnakeWord(s []rune, word []rune) []rune {
	for i, r := range word {
		word[i] = unicode.ToLower(r)
	}
	if len(s) > 0 {
		s = append(s, '_')
	}
	return append(s, word...)
}
