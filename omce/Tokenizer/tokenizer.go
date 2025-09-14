package tokenizer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Symbol uint

const (
	INCLUDE  Symbol = iota // '#include'
	FUNCTION               // 'function'
	TEXT                   // Unknown things
	ESCAPE                 // '\'
	NEWLINE                // 줄바꿈
	COMPOUND               // '{'랑 '}'로 묶인거. 안에는 메서드 코드가 들어감
)

type Token struct {
	Type  Symbol
	Value string
	Line  int
}

func (tok Token) String() string {
	return fmt.Sprintf("[%d]:%s", tok.Type, tok.Value)
}

func lineNumberAtIndex(s string, idx int) int {
	if idx < 0 || idx > len(s) {
		return -1
	}

	// \r\n, \n\r, \n, \r 모두 처리
	re := regexp.MustCompile(`\r\n|\n\r|\n|\r`)

	// idx 이전까지의 문자열에서 줄바꿈 개수 세기
	sub := s[:idx]
	matches := re.FindAllStringIndex(sub, -1)

	// 줄 번호 = 줄바꿈 개수 + 1
	return len(matches) + 1
}

// 중괄호 내부 문자열 불러오기.
func itsBoring(rune_code []rune, index *int) string {
	brace_stack := 0
	var code strings.Builder

	length := len(rune_code)
	for {
		switch rune_code[*index] {
		case '{':
			brace_stack += 1
			code.WriteString("{")
		case '}':
			brace_stack -= 1
			code.WriteString("}")
		case '\\':
			*index++
			if *index < length && (rune_code[*index] != '{' && rune_code[*index] != '}') {
				code.WriteString("\\")
			}
			code.WriteRune(rune_code[*index])
		default:
			code.WriteRune(rune_code[*index])
		}
		*index++
		if !(*index < length && brace_stack > 0) {
			*index--
			break
		}
	}
	new_code := code.String()
	new_code = strings.TrimPrefix(new_code, "{")
	new_code = strings.TrimSuffix(new_code, "}")
	new_code = strings.TrimSpace(new_code)
	// fmt.Println("FUCK:", new_code)
	return new_code
}

func Tokenize(code string) []Token {
	tokens := make([]Token, 0)

	rune_code := []rune(code)

	var words string
	var idx int

	newToken := func(symbol_name Symbol, words string) Token {
		var new_token Token
		new_token.Type = symbol_name
		new_token.Value = words
		new_token.Line = lineNumberAtIndex(code, idx)
		return new_token
	}
	empty_words := func() {
		if words != "" {
			tokens = append(tokens, newToken(TEXT, words))
			words = ""
		}
	}

	length := len(rune_code)
	for idx < length {
		rune_chr := rune_code[idx]
		chr := string(rune_chr)

		// 줄바꿈 처리
		if rune_chr == '\n' || rune_chr == '\r' {
			empty_words()
			// CRLF, LFCR 같이 두 개 묶음이면 하나로 처리
			if rune_chr == '\r' && idx+1 < len(code) && rune_code[idx+1] == '\n' {
				tokens = append(tokens, newToken(NEWLINE, ""))
				idx += 2
				continue
			}
			if rune_chr == '\n' && idx+1 < len(code) && rune_code[idx+1] == '\r' {
				tokens = append(tokens, newToken(NEWLINE, ""))
				idx += 2
				continue
			}
			tokens = append(tokens, newToken(NEWLINE, ""))
			idx++
			continue
		}

		if unicode.IsSpace(rune_chr) {
			empty_words()
			idx++
			continue
		}

		switch chr {
		case "{":
			empty_words()
			tokens = append(tokens, newToken(COMPOUND, itsBoring(rune_code, &idx)))
			idx++
			continue
		case "\\":
			empty_words()
			tokens = append(tokens, newToken(ESCAPE, chr))
			idx++
			continue
		}
		words += string(chr)
		words_len := len(words)
		if idx-words_len == -1 || (idx-words_len >= 0 && (rune_code[idx-words_len] == '\n' || rune_code[idx-words_len] == '\r')) {
			switch words {
			case "#include":
				tokens = append(tokens, newToken(INCLUDE, words))
				words = ""
			case "function":
				tokens = append(tokens, newToken(FUNCTION, words))
				words = ""
			}
		}
		idx++
	}
	empty_words()
	return tokens
}
