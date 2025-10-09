package lexer

import (
	perr "flow/omfn/ParseErrors"
	tok "flow/omfn/Tokenizer"
)

type NodeType int

const (
	ROOT                 NodeType = iota
	FUNCTION_DEFINITION           // field: Body
	VARIABLE_DECLARATION          // field: Name, Value, Body
	CALL_EXPRESSION               // field: Body
	NUMBER                        // field: Value
	SELECTOR                      // field: Value
	STRING                        // field: Value
	IDENTIFIER                    // field: Name
	PARAM_LIST                    // field: Body
	PARAM                         // field: Name, Value as ParamType
	TYPE                          // field: Value
	COMPOUND_STATEMENT            // field: Body
	ARGUMENT_LIST                 // field: Body
	RAWLINE                       // field: Value
	RETURN_EXPRESSION             // field: Body
)

type Node struct {
	Type  NodeType
	Body  []Node
	Name  string
	Value string
	Begin int
	End   int
}

type Lexer struct {
	idx    int
	length int
	tokens []tok.Token
	ast    Node
}

func (lexer *Lexer) cur() tok.Token {
	return lexer.tokens[lexer.idx]
}

func (lexer *Lexer) next(i int) tok.Token {
	if lexer.idx+i >= lexer.length {
		return tok.Token{Type: tok.UNKNOWN, Value: ""}
	}
	return lexer.tokens[lexer.idx+i]
}

// func (lexer Lexer) getNext(idx int) (tok.Token, error) {
// 	if lexer.idx+idx >= lexer.length {
// 		return lexer.cur(), &perr.TokenError{
// 			Begin: lexer.cur().Begin,
// 			End:   lexer.cur().End,
// 			Msg:   "토큰을 예상했으나 EOF",
// 		}
// 	}
// 	lexer.idx += idx
// 	return lexer.cur(), nil
// }

func (lexer *Lexer) advance() {
	if lexer.idx+1 >= lexer.length {
		panic(&perr.TokenError{
			Begin: lexer.cur().Begin,
			End:   lexer.cur().End,
			Msg:   "토큰을 예상했으나 EOF",
		})
	}
	lexer.idx++
}

func (lexer *Lexer) eat(tt tok.TokenType) tok.Token {
	cur := lexer.cur()
	lexer.expect(tt)
	lexer.advance()
	return cur
}

func (lexer *Lexer) expect(t tok.TokenType) {
	if lexer.cur().Type == t {
		return
	}
	panic(&perr.TokenError{
		Begin: lexer.cur().Begin,
		End:   lexer.cur().End,
		Msg:   perr.UnexpectedToken(lexer.cur().Value, tok.TokTypeToString(t)),
	})
}

func (lexer *Lexer) argumentList() (returnNode Node, returnErr error) {
	newNode := Node{}
	newNode.Type = ARGUMENT_LIST

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	newNode.Begin = lexer.eat(tok.LPAREN).Begin

	for lexer.cur().Type != tok.RPAREN {

		switch lexer.cur().Type {
		case tok.IDENTIFIER:
			newNode.Body = append(newNode.Body, Node{Type: IDENTIFIER, Name: lexer.cur().Value})
		case tok.NUMBER:
			newNode.Body = append(newNode.Body, Node{Type: NUMBER, Value: lexer.cur().Value})
		case tok.STRING:
			newNode.Body = append(newNode.Body, Node{Type: STRING, Value: lexer.cur().Value})
		case tok.SELECTOR:
			newNode.Body = append(newNode.Body, Node{Type: SELECTOR, Value: lexer.cur().Value})
		default:
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   perr.UnexpectedToken(lexer.cur().Value, "'식별자', '숫자', '문자열', '선택자'를 예상함"),
			})
		}

		lexer.advance()
		if lexer.cur().Type == tok.RPAREN {
			break
		}
		if lexer.cur().Type == tok.EOF {
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   perr.UnexpectedToken(lexer.cur().Value, "')'를 예상함"),
			})
		}
		lexer.eat(tok.COMMA)
	}

	newNode.End = lexer.eat(tok.RPAREN).End

	return
}

func (lexer *Lexer) parameterList() (returnNode Node, returnErr error) {
	newNode := Node{}
	newNode.Type = PARAM_LIST

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	newNode.Begin = lexer.eat(tok.LPAREN).Begin

	for lexer.cur().Type != tok.RPAREN {

		paramType := lexer.eat(tok.TYPE)

		paramName := lexer.eat(tok.IDENTIFIER)

		newNode.Body = append(newNode.Body, Node{Type: PARAM, Name: paramName.Value, Value: paramType.Value})

		if lexer.cur().Type == tok.RPAREN {
			break
		}
		if lexer.cur().Type == tok.EOF {
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   perr.UnexpectedToken(lexer.cur().Value, "')'를 예상함"),
			})
		}
		lexer.eat(tok.COMMA)
	}

	newNode.Begin = lexer.eat(tok.RPAREN).Begin

	return
}

func (lexer *Lexer) returnExpression() (returnNode Node, returnErr error) {
	newNode := Node{}
	newNode.Type = RETURN_EXPRESSION

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	lexer.eat(tok.RETURN)

	for {
		if lexer.cur().Type == tok.SEMICOLON {
			newNode.End = lexer.cur().End
			lexer.advance()
			break
		}
		if lexer.cur().Type == tok.EOF {
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   perr.UnexpectedToken(lexer.cur().Value, "';'를 예상함"),
			})
		}

		switch lexer.cur().Type {
		case tok.IDENTIFIER:
			newNode.Body = append(newNode.Body, Node{Type: IDENTIFIER, Name: lexer.cur().Value})
		case tok.NUMBER:
			newNode.Body = append(newNode.Body, Node{Type: NUMBER, Value: lexer.cur().Value})
		case tok.STRING:
			newNode.Body = append(newNode.Body, Node{Type: STRING, Value: lexer.cur().Value})
		case tok.SELECTOR:
			newNode.Body = append(newNode.Body, Node{Type: SELECTOR, Value: lexer.cur().Value})
		default:
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   perr.UnexpectedToken(lexer.cur().Value, "'식별자', '숫자', '문자열', '선택자'를 예상함"),
			})
		}

		lexer.advance()
	}

	return
}

func (lexer *Lexer) callExpression() (returnNode Node, returnErr error) {
	newNode := Node{}
	newNode.Type = CALL_EXPRESSION

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	iddd := lexer.eat(tok.IDENTIFIER)
	newNode.Name = iddd.Value
	newNode.Begin = iddd.Begin
	newNode.End = iddd.End

	newArgsNode, err := lexer.argumentList()
	if err != nil {
		panic(err)
	}
	newNode.Body = append(newNode.Body, newArgsNode)

	return
}

func (lexer *Lexer) compoundStatement() (returnNode Node, returnErr error) {
	newNode := Node{}
	newNode.Type = COMPOUND_STATEMENT

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	newNode.Begin = lexer.eat(tok.LBRACE).Begin

	for {
		if lexer.cur().Type == tok.RBRACE {
			newNode.End = lexer.cur().End
			lexer.advance()
			break
		}
		if lexer.cur().Type == tok.EOF {
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   "CompoundStatement의 끝을 찾지 못 함",
			})
		}

		if lexer.cur().Type == tok.IDENTIFIER {
			newCallNode, err := lexer.callExpression()
			if err != nil {
				panic(err)
			}
			lexer.eat(tok.SEMICOLON)
			newNode.Body = append(newNode.Body, newCallNode)
		}

		if lexer.cur().Type == tok.TYPE {
			newTypeSpecifierNode, err := lexer.typeSpecifier()
			if err != nil {
				panic(err)
			}
			newNode.Body = append(newNode.Body, newTypeSpecifierNode)
		}

		if lexer.cur().Type == tok.RETURN {
			newReturnExpressionNode, err := lexer.returnExpression()
			if err != nil {
				panic(err)
			}
			newNode.Body = append(newNode.Body, newReturnExpressionNode)
		}

		if lexer.cur().Type == tok.RAWLINE {
			newNode.Body = append(newNode.Body, Node{Type: RAWLINE, Value: lexer.cur().Value})
			lexer.advance()
		}
	}

	return
	// COMPOUNDSTATEMENT:

	// 	if lexer.cur().Type != tok.RBRACE {
	// 		goto COMPOUNDSTATEMENT
	// 	}

}

func (lexer *Lexer) functionDefinition() (returnNode Node, returnErr error) {
	newNode := Node{}
	newNode.Type = FUNCTION_DEFINITION

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	lexer.expect(tok.TYPE)
	newNode.Begin = lexer.cur().Begin
	newNode.Body = append(newNode.Body, Node{Type: TYPE, Value: lexer.cur().Value}) // 함수 반환 타입
	lexer.advance()

	lexer.expect(tok.IDENTIFIER)
	newNode.End = lexer.cur().End
	newNode.Name = lexer.cur().Value // 함수 이름
	lexer.advance()

	paramNode, err := lexer.parameterList()
	if err != nil {
		panic(err)
	}
	newNode.Body = append(newNode.Body, paramNode)

	compoundNode, err := lexer.compoundStatement()
	if err != nil {
		panic(err)
	}
	newNode.Body = append(newNode.Body, compoundNode)
	return
}

func (lexer *Lexer) variableAssignment() (returnNode Node, returnErr error) {
	newNode := Node{}
	newNode.Type = VARIABLE_DECLARATION

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	newNode.Begin = lexer.cur().Begin
	newNode.Value = lexer.eat(tok.TYPE).Value
	newNode.Name = lexer.eat(tok.IDENTIFIER).Value
	lexer.eat(tok.EQUAL)

	for {

		switch lexer.cur().Type {
		case tok.IDENTIFIER:
			newNode.Body = append(newNode.Body, Node{Type: IDENTIFIER, Name: lexer.cur().Value})
		case tok.NUMBER:
			newNode.Body = append(newNode.Body, Node{Type: NUMBER, Value: lexer.cur().Value})
		case tok.STRING:
			newNode.Body = append(newNode.Body, Node{Type: STRING, Value: lexer.cur().Value})
		case tok.SELECTOR:
			newNode.Body = append(newNode.Body, Node{Type: SELECTOR, Value: lexer.cur().Value})
		default:
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   perr.UnexpectedToken(lexer.cur().Value, "식별자, 숫자, 문자열, 선택자"),
			})
		}

		lexer.advance()

		if lexer.cur().Type == tok.SEMICOLON {
			newNode.End = lexer.cur().End
			lexer.advance()
			break
		}
		if lexer.cur().Type == tok.EOF {
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   perr.UnexpectedToken(lexer.cur().Value, "';'를 예상함"),
			})
		}
	}

	return
}

func (lexer *Lexer) typeSpecifier() (Node, error) {

	if lexer.next(1).Type == tok.IDENTIFIER &&
		lexer.next(2).Type == tok.EQUAL {
		newNode, err := lexer.variableAssignment()
		if err != nil {
			return newNode, err
		}
		return newNode, nil
	}

	return Node{}, &perr.TokenError{
		Begin: lexer.cur().Begin,
		End:   lexer.cur().End,
		Msg:   "알 수 없는 구조의 코드",
	}
}

func (lexer *Lexer) Lexicalize() (Node, error) {

	for lexer.idx < lexer.length && lexer.cur().Type != tok.EOF {
		if lexer.cur().Type == tok.TYPE {
			newDefinition, err := lexer.functionDefinition()
			if err != nil {
				return Node{}, err
			}
			lexer.ast.Body = append(lexer.ast.Body, newDefinition)
		}
	}

	return lexer.ast, nil
}

func New(toks []tok.Token) Lexer {
	lex := Lexer{}
	lex.tokens = toks
	lex.ast = Node{}
	lex.ast.Type = ROOT
	lex.length = len(toks)
	lex.idx = 0
	return lex
}
