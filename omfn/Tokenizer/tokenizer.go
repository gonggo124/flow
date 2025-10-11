package tokenizer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	TYPE TokenType = iota
	IDENTIFIER
	LBRACE
	RBRACE
	LPAREN
	RPAREN
	SEMICOLON
	COMMA
	STRING
	SELECTOR
	NUMBER
	EQUAL
	RAWLINE
	RETURN
	OPERATOR
	UNKNOWN
	EOF
)

type Token struct {
	Type  TokenType
	Begin int
	End   int
	Value string
}

func TokTypeToString(tok TokenType) string {
	switch tok {
	case TYPE:
		return "TYPE"
	case IDENTIFIER:
		return "IDENTIFIER"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case SEMICOLON:
		return "SEMICOLON"
	case COMMA:
		return "COMMA"
	case STRING:
		return "STRING"
	case SELECTOR:
		return "SELECTOR"
	case NUMBER:
		return "NUMBER"
	case EQUAL:
		return "EQUAL"
	case RAWLINE:
		return "RAWLINE"
	case EOF:
		return "EOF"
	case RETURN:
		return "RETURN"
	case OPERATOR:
		return "OPERATOR"
	default:
		return "UNEXPECTED TOKEN"
	}
}

func (tok Token) String() string {
	return fmt.Sprintf("[%d]:%s %d~%d", tok.Type, strings.ReplaceAll(tok.Value, "\n", "\\n"), tok.Begin, tok.End)
}

var selector_rule = regexp.MustCompile(`@([arpnes])(?:\[([^,]*=(?:{.*}|[^,]*)(?:,[^,]*=(?:{.*}|[^,]*))*)?\])?`)

// hint: target
func Tokenize(tCode string) []Token {
	token_list := make([]Token, 0)

	begin := 0
	// end := 0

	current := ""

	inDQuotes := false
	inQuotes := false
	inRawLine := false

	newToken := func(t TokenType, val string) Token {
		var nT Token
		nT.Type = t
		nT.Value = val
		nT.Begin = begin
		nT.End = begin + len(val)
		return nT
	}

	pushToken := func(t Token) {
		token_list = append(token_list, t)
		current = ""
		begin += len(t.Value)
	}

	emptyCurrent := func() {
		if current == "" {
			return
		}
		if selector_rule.MatchString(current) {
			pushToken(newToken(SELECTOR, current))
			return
		}
		if _, err := strconv.Atoi(current); err == nil {
			pushToken(newToken(NUMBER, current))
			return
		}
		switch current {
		case "int", "void", "selector":
			pushToken(newToken(TYPE, current))
		case "return":
			pushToken(newToken(RETURN, current))
		default:
			pushToken(newToken(IDENTIFIER, current))
		}
		begin += len(current)
	}

	for idx, chr := range tCode {
		// end = idx

		if inRawLine {
			if chr == '\n' || chr == '\r' {
				inRawLine = false
				pushToken(newToken(RAWLINE, current))
				continue
			}
			fmt.Println(string(chr), chr, "FUCK?")
			current += string(chr)
			continue
		}
		if inQuotes || inDQuotes {
			if inQuotes && chr == '\'' {
				inQuotes = false
				pushToken(newToken(STRING, current))
				continue
			}
			if inDQuotes && chr == '"' {
				inDQuotes = false
				pushToken(newToken(STRING, current))
				continue
			}
			current += string(chr)
			continue
		}
		switch chr {
		case '`':
			emptyCurrent()
			inRawLine = true
		case '"':
			emptyCurrent()
			inDQuotes = true
		case '\'':
			emptyCurrent()
			inQuotes = true
		case '{':
			emptyCurrent()
			pushToken(newToken(LBRACE, string(chr)))
		case '}':
			emptyCurrent()
			pushToken(newToken(RBRACE, string(chr)))
		case '(':
			emptyCurrent()
			pushToken(newToken(LPAREN, string(chr)))
		case ')':
			emptyCurrent()
			pushToken(newToken(RPAREN, string(chr)))
		case ',':
			emptyCurrent()
			pushToken(newToken(COMMA, string(chr)))
		case '=':
			emptyCurrent()
			pushToken(newToken(EQUAL, string(chr)))
		case ';':
			emptyCurrent()
			pushToken(newToken(SEMICOLON, string(chr)))
		case '+', '-', '*', '/':
			emptyCurrent()
			pushToken(newToken(OPERATOR, string(chr)))
		default:
			if unicode.IsSpace(chr) {
				if current == "" {
					begin = idx
				} else {
					emptyCurrent()
				}
			} else {
				current += string(chr)
			}
		}
	}

	pushToken(newToken(EOF, ""))

	return token_list
}
