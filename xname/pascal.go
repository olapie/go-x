package xname

import (
	"fmt"
	"unicode"
)

func ToPascal(s string, options ...Option) string {
	return string(toPascal([]rune(s), options...))
}

func toPascal(input []rune, options ...Option) []rune {
	type State int
	const (
		Start State = iota
		Lower
		Upper
		LowerWordEnd
		LastLowerWordEnd
		UpperStart
		End
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
	output := make([]rune, 0, len(input))
	n := len(input)
	lowerStart := -1
	for i := 0; i < n; i++ {
		r := input[i]
		switch state {
		case Start:
			switch r {
			case '_', '.', '-', ' ', '\t':
				break
			default:
				if unicode.IsUpper(r) {
					state = Upper
					output = append(output, r)
				} else {
					state = Lower
					lowerStart = i
				}
			}
		case Lower:
			switch r {
			case '_', '.', '-', ' ', '\t':
				state = LowerWordEnd
				output = appendWord(output, input[lowerStart:i], useAcronym)
				i-- // null transition
				state = Start
			default:
				if unicode.IsUpper(r) {
					state = UpperStart
					output = appendWord(output, input[lowerStart:i], useAcronym)
					i-- // null transition
					state = Upper
				}
			}
		case Upper:
			switch r {
			case '_', '.', '-', ' ', '\t':
				state = Start
			default:
				output = append(output, r)
			}
		default:
			panic(fmt.Sprintf("invalid state: %d", state))
		}
	}

	switch state {
	case Start:
		state = End
	case Upper:
		state = End
	case Lower:
		state = LastLowerWordEnd
		output = appendWord(output, input[lowerStart:], useAcronym)
		state = End
	}

	return output
}

func appendWord(s []rune, word []rune, useAcronym bool) []rune {
	if useAcronym {
		acronymsMapRWMutex.RLock()
		acronym := acronymsL2U[string(word)]
		acronymsMapRWMutex.RUnlock()
		if acronym != "" {
			return append(s, []rune(acronym)...)
		}
	}
	word[0] = unicode.ToUpper(word[0])
	return append(s, word...)
}
