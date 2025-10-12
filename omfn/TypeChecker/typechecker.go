package typechecker

import (
	lexer "flow/omfn/Lexer"
	perr "flow/omfn/ParseErrors"
	"fmt"
)

type Function struct {
	Type      string
	Definited bool
}

func Annotate(ast *lexer.Node) (returnErr error) {
	defer func() {
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(*perr.TokenError); ok {
				returnErr = e // paniced error를 반환
			}
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	readTree(ast, make(map[string]Function), make(map[string]string))

	return
}

func readTree(ast *lexer.Node, functions map[string]Function, variables map[string]string) {
	switch ast.Type {
	case lexer.FUNCTION_DEFINITION, lexer.FUNCTION_DECLARATION:
		if functions[ast.Name].Type != "" && (functions[ast.Name].Definited || ast.Type == lexer.FUNCTION_DECLARATION) {
			panic(&perr.TokenError{
				Begin: ast.Begin,
				End:   ast.End,
				Msg:   "함수를 중복 선언할 수 없습니다",
			})
		}
		something := functions[ast.Name]
		something.Type = ast.Value
		if ast.Type == lexer.FUNCTION_DEFINITION {
			something.Definited = true
		}
		functions[ast.Name] = something
		for i := range ast.Body {
			readTree(&ast.Body[i], functions, variables)
		}
	case lexer.VARIABLE_DECLARATION:
		if variables[ast.Name] != "" {
			fmt.Println(variables[ast.Name])
			panic(&perr.TokenError{
				Begin: ast.Begin,
				End:   ast.End,
				Msg:   "변수를 중복 선언할 수 없습니다",
			})
		}
		variables[ast.Name] = ast.Value
		for i := range ast.Body {
			readTree(&ast.Body[i], functions, variables)
		}
	case lexer.IDENTIFIER:
		datatype, ok := variables[ast.Name]
		if !ok {
			panic(&perr.TokenError{
				Begin: ast.Begin,
				End:   ast.End,
				Msg:   "아직 선언되지 않은 변수에 접근했습니다",
			})
		}
		ast.DataType = datatype
	case lexer.CALL_EXPRESSION:
		datatype, ok := functions[ast.Name]
		if !ok {
			panic(&perr.TokenError{
				Begin: ast.Begin,
				End:   ast.End,
				Msg:   "아직 선언되지 않은 함수에 접근했습니다",
			})
		}
		ast.DataType = datatype.Type
	case lexer.NUMBER:
		ast.DataType = "int"
		ast.DataType = "int"
	case lexer.STRING:
		ast.DataType = "string"
	case lexer.SELECTOR:
		ast.DataType = "selector"
	default:
		for i := range ast.Body {
			readTree(&ast.Body[i], functions, variables)
		}
	}
}
