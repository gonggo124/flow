package lexer

import (
	"bufio"
	"fmt"
	"strings"

	utils "flow/Utils"
	cmdvh "flow/VersionHandler"
	node "flow/omce/Node"
	token "flow/omce/Tokenizer"
	anonyfunc "flow/omfn/Anonyfunc"
	fnreader "flow/omfn/Reader"
)

// gpt야 고마워!!!!!!
// func preprocess(input string) string {
// 	// 1) key에 쌍따옴표 붙이기 (Tags:, block_state:, Name: 등)
// 	reKey := regexp.MustCompile(`(\w+):`)
// 	input = reKey.ReplaceAllString(input, `"$1":`)

// 	// 2) 문자열 값에 쌍따옴표가 붙어있지만, 혹시 빠진 곳 있으면 넣어주는 간단 처리 (필요시)
// 	// (여기선 이미 문자열은 " "로 되어 있으므로 스킵)

// 	// 3) 숫자 뒤에 붙은 f를 제거 (0f -> 0)
// 	reFloatF := regexp.MustCompile(`(\d+(\.\d+)?)[fF]`)
// 	input = reFloatF.ReplaceAllString(input, `$1`)

// 	return input
// }

// 중괄호 내부 문자열 불러오기.
// func itsBoring(tokens []token.Token, index *int) string {
// 	brace_stack := 0
// 	var code strings.Builder

// 	length := len(tokens)
// 	for {
// 		switch tokens[*index].Type {
// 		case token.LBRACE:
// 			brace_stack += 1
// 			code.WriteString("{")
// 		case token.RBRACE:
// 			brace_stack -= 1
// 			code.WriteString("}")
// 		case token.ESCAPE:
// 			*index++
// 			fmt.Println("'"+code.String()+"'", "으악!")
// 			if *index < length && tokens[*index].Type != token.LBRACE {
// 				code.WriteString("\\")
// 			}
// 			code.WriteString(tokens[*index].Value + " ")
// 		case token.NEWLINE:
// 			if code.Len() > 0 {
// 				code.WriteString("\n")
// 			}
// 		default:
// 			code.WriteString(tokens[*index].Value + " ")
// 		}
// 		*index++
// 		if !(*index < length && brace_stack > 0) {
// 			*index--
// 			break
// 		}
// 	}
// 	new_code := code.String()
// 	new_code = strings.TrimPrefix(new_code, "{\n")
// 	new_code = strings.TrimSuffix(new_code, "}")
// 	return new_code
// }

func Lexicalize(tokens []token.Token) node.Node {
	ast := node.NewNode(node.PROGRAM)

	index := 0
	length := len(tokens)

	consume := func(errMsg string, expected ...token.Symbol) *token.Token {
		index++
		if index >= length {
			utils.SetLine(tokens[index-1].Line)
			utils.Panic(errMsg)
		}
		t := &tokens[index]
		utils.SetLine(t.Line)

		is_ok := false
		for _, ex := range expected {
			if t.Type == ex {
				is_ok = true
			}
		}
		if !is_ok {
			utils.Panic(errMsg)
		}
		return t
	}

	for index < length {
		tok := tokens[index]
		utils.SetLine(tok.Line)
		switch tok.Type {
		case token.NEWLINE:
			// i don't do anything haha
		case token.INCLUDE:
			new_extender := node.NewNode(node.INCLUDE_DECLARATION)
			new_identifier := node.NewNode(node.IDENTIFIER)
			new_identifier.Name = consume("#include 다음에는 method set id가 와야합니다", token.TEXT).Value
			new_extender.Body = append(new_extender.Body, new_identifier)
			ast.Body = append(ast.Body, new_extender)
		case token.FUNCTION:
			new_method := node.NewNode(node.METHOD_DECLARATION)
			method_id := consume("'function' 키워드 다음에는 함수 식별자가 와야합니다", token.TEXT)
			new_identifier := node.NewNode(node.IDENTIFIER)
			new_identifier.Name = method_id.Value
			new_method.Id = &new_identifier
			consume("함수 식별자 다음에는 '{' 또는 함수 경로가 와야합니다", token.COMPOUND, token.TEXT)

			new_func_path := node.NewNode(node.IDENTIFIER)
			switch tokens[index].Type {
			case token.TEXT:
				new_func_path.Name = tokens[index].Value
			case token.COMPOUND:
				var writer strings.Builder

				scanner := bufio.NewScanner(strings.NewReader(tokens[index].Value))

				var reader fnreader.Reader
				reader.SetScanner(scanner)

				for reader.Scan() {
					line := reader.Text()
					line.Text = strings.TrimSpace(line.Text)
					writer.WriteString(cmdvh.ParseCmd(&reader, line) + "\n")
				}

				anonyfunc_path := anonyfunc.New(writer.String())
				new_func_path.Name = anonyfunc_path
			default:
				utils.Panic("알 수 없는 코드입니다")
			}
			new_method.Body = append(new_method.Body, new_func_path)
			ast.Body = append(ast.Body, new_method)
		default:
			utils.Panic(fmt.Sprintf("예상치 못한 토큰: '%s'", tokens[index].Value))
		}
		index++
	}

	return ast
}
