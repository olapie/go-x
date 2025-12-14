package naming

import (
	"fmt"
	"unicode"
)

func ToCamel(s string, options ...Option) string {
	input := []rune(s)
	type State int
	const (
		Start State = iota
		FirstUpper
		SecondUpper
		Lower
		FirstWordEnd
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
	for i, r := range input {
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
				state = FirstWordEnd
				return string(toCamel(input, wordStart, i, options))
			}
			if unicode.IsUpper(r) {
				state = SecondUpper
			} else {
				state = Lower
			}
		case SecondUpper:
			if isDelimiter(r) {
				state = FirstWordEnd
				return string(toCamel(input, wordStart, i, options))
			}
			if unicode.IsUpper(r) {
				state = UpperWord
			} else {
				state = Lower
			}
		case Lower:
			if isDelimiter(r) || unicode.IsUpper(r) {
				state = FirstWordEnd
				return string(toCamel(input, wordStart, i, options))
			}
		case UpperWord:
			if isDelimiter(r) {
				state = FirstWordEnd
				return string(toCamel(input, wordStart, i, options))
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
				state = FirstWordEnd
				return string(toCamel(input, wordStart, i-1, options))
			}
			state = Lower
		default:
			panic(fmt.Sprintf("invalid state: %d", state))
		}
	}

	if wordStart < 0 {
		return ""
	}
	n := len(input)
	for i := wordStart; i < n; i++ {
		input[i] = unicode.ToLower(input[i])
	}
	return string(input)
}

func toCamel(input []rune, firstWordStart, firstWordEnd int, options []Option) []rune {
	word := input[firstWordStart:firstWordEnd]
	for i := range word {
		word[i] = unicode.ToLower(word[i])
	}
	pascal := toPascal(input[firstWordEnd:], options...)
	return append(word, pascal...)
}
