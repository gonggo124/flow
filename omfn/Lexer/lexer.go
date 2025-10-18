package lexer

import (
	perr "flow/omfn/ParseErrors"
	tok "flow/omfn/Tokenizer"
	"fmt"
)

type NodeType int

const (
	ROOT                  NodeType = iota
	FUNCTION_DEFINITION            // field: Body, Value as ReturnType
	FUNCTION_DECLARATION           // field: Body, Value as ReturnType
	VARIABLE_DECLARATION           // field: Body, Name, Value as VariableType
	CALL_EXPRESSION                // field: Body
	BINARY_EXPRESSION              // field: Body
	ASSIGNMENT_EXPRESSION          // field: Body
	LITERAL_EXPRESSION             // field: Body
	OPERATOR                       // field: Value
	NUMBER                         // field: Value
	SELECTOR                       // field: Value
	STRING                         // field: Value
	IDENTIFIER                     // field: Name
	PARAM_LIST                     // field: Body
	PARAM                          // field: Name, Value as ParamType
	COMPOUND_STATEMENT             // field: Body
	ARGUMENT_LIST                  // field: Body
	RAWLINE                        // field: Value
	RETURN_EXPRESSION              // field: Body
)

type Node struct {
	Type     NodeType
	Body     []Node
	Name     string
	Value    string
	Begin    int
	End      int
	DataType string // TypeChecker가 사용.
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
		panic(&perr.TokenError{
			Begin: lexer.cur().Begin,
			End:   lexer.cur().End,
			Msg:   "토큰을 예상했으나 EOF",
		})
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

		newExpr, err := lexer.expression([]tok.TokenType{tok.COMMA, tok.RPAREN})
		if err != nil {
			panic(err)
		}
		newNode.Body = append(newNode.Body, newExpr)

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

	newNode.End = lexer.eat(tok.RPAREN).End

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

		if lexer.cur().Type == tok.IDENTIFIER && lexer.next(1).Type == tok.LPAREN {
			newCallNode, err := lexer.callExpression()
			if err != nil {
				panic(err)
			}
			lexer.eat(tok.SEMICOLON)
			newNode.Body = append(newNode.Body, newCallNode)
		}

		if lexer.cur().Type == tok.IDENTIFIER && lexer.next(1).Type == tok.EQUAL {
			newAssignNode, err := lexer.variableAssignment()
			if err != nil {
				panic(err)
			}
			newNode.Body = append(newNode.Body, newAssignNode)
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
			newNode.Body = append(newNode.Body, Node{
				Type: RAWLINE,
				Value: lexer.cur().Value,
				Begin: lexer.cur().Begin,
				End: lexer.cur().End,
			})
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
	newNode.Value = lexer.cur().Value // 함수 반환 타입
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

	if lexer.cur().Type == tok.SEMICOLON {
		lexer.eat(tok.SEMICOLON)
		newNode.Type = FUNCTION_DECLARATION
		return
	}

	compoundNode, err := lexer.compoundStatement()
	if err != nil {
		panic(err)
	}
	newNode.Body = append(newNode.Body, compoundNode)

	return
}

func (lexer *Lexer) expression(end_toks []tok.TokenType) (returnNode Node, returnErr error) {
	newNode := Node{}
	// newNode.Type =

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	check_operand := func(t Node) Node {
		if t.Type != NUMBER && t.Type != SELECTOR && t.Type != BINARY_EXPRESSION &&
			t.Type != CALL_EXPRESSION && t.Type != IDENTIFIER {
			panic(&perr.TokenError{
				Begin: t.Begin,
				End:   t.End,
				Msg:   "숫자, 선택자, 변수, 함수를 예상했으나 '" + t.Value + "'을 찾음",
			})
		}
		return t
	}
	check_op := func(t Node) Node {
		if t.Type != OPERATOR {
			panic(&perr.TokenError{
				Begin: t.Begin,
				End:   t.End,
				Msg:   "연산자를 예상했으나 '" + t.Value + "'을 찾음",
			})
		}
		return t
	}

	parse := func(exprs []Node) Node {
		idx := 0
		turn := true
		for len(exprs) > 2 {
			if turn && len(exprs)-idx < 3 {
				turn = false
				idx = 0
			}
			fmt.Println(exprs, turn, idx)
			op1 := check_operand(exprs[idx])
			opr := check_op(exprs[idx+1])
			op2 := check_operand(exprs[idx+2])
			if turn && opr.Value != "*" && opr.Value != "/" {
				idx += 2
			} else {
				if idx+3 >= len(exprs) {
					exprs = exprs[:idx]
				} else {
					exprs = append(exprs[:idx], exprs[idx+3:]...)
				}
				exprs = append(exprs[:idx+1], exprs[idx:]...)
				exprs[idx] = Node{
					Type: BINARY_EXPRESSION,
					Body: []Node{op1, opr, op2},
				}
				idx = 0
			}

		}

		return exprs[0]
	}

	var tok_to_node func(t tok.Token) Node
	tok_to_node = func(t tok.Token) Node {
		switch t.Type {
		case tok.LPAREN:
			exprss := make([]Node, 0)
			lexer.advance() // '(' 무한 재귀 방지
			if lexer.cur().Type == tok.RPAREN {
				panic(&perr.TokenError{
					Begin: t.Begin,
					End:   t.End,
					Msg:   "숫자, 선택자, 변수, 함수를 예상했으나 '" + lexer.cur().Value + "'을 찾음",
				})
			}
			for lexer.cur().Type != tok.RPAREN {
				exprss = append(exprss, tok_to_node(lexer.cur()))
				lexer.advance()
			}
			resultNode := parse(exprss)
			resultNode.Begin = t.Begin
			resultNode.End = exprss[len(exprss)-1].End
			return resultNode
		case tok.OPERATOR:
			return Node{
				Type:  OPERATOR,
				Value: lexer.cur().Value,
				Begin: t.Begin,
				End:   t.End,
			}
		case tok.NUMBER:
			return Node{
				Type:  NUMBER,
				Value: lexer.cur().Value,
				Begin: t.Begin,
				End:   t.End,
			}
		case tok.SELECTOR:
			return Node{
				Type:  SELECTOR,
				Value: lexer.cur().Value,
				Begin: t.Begin,
				End:   t.End,
			}
		case tok.IDENTIFIER:
			if lexer.next(1).Type == tok.LPAREN {
				newCallExpr, err := lexer.callExpression()
				if err != nil {
					panic(err)
				}
				lexer.idx--
				return newCallExpr
			}
			return Node{
				Type:  IDENTIFIER,
				Name:  lexer.cur().Value,
				Begin: t.Begin,
				End:   t.End,
			}
		}
		panic(&perr.TokenError{
			Begin: t.Begin,
			End:   t.End,
			Msg:   perr.UnexpectedToken(tok.TokTypeToString(t.Type), "숫자, 선택자, 변수, 함수, 연산자"),
		})
	}

	is_end := func(t tok.Token) bool {
		for _, ttt := range end_toks {
			if ttt == t.Type {
				return true
			}
		}
		return false
	}

	exprs := make([]Node, 0)

	for !is_end(lexer.cur()) {
		if lexer.cur().Type == tok.EOF {
			panic(&perr.TokenError{
				Begin: lexer.cur().Begin,
				End:   lexer.cur().End,
				Msg:   perr.UnexpectedToken(lexer.cur().Value, "';'를 예상함"),
			})
		}
		exprs = append(exprs, tok_to_node(lexer.cur()))
		lexer.advance()
	}

	begin := exprs[0].Begin
	end := exprs[len(exprs)-1].End
	newNode = parse(exprs)
	newNode.Begin = begin
	newNode.End = end

	return
}

func (lexer *Lexer) variableDeclaration() (returnNode Node, returnErr error) {
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

	newExpr, err := lexer.expression([]tok.TokenType{tok.SEMICOLON})
	if err != nil {
		panic(err)
	}
	newNode.Body = append(newNode.Body, newExpr)

	newNode.End = lexer.eat(tok.SEMICOLON).End

	return
}

func (lexer *Lexer) variableAssignment() (returnNode Node, returnErr error) {
	newNode := Node{}
	newNode.Type = ASSIGNMENT_EXPRESSION

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
	newNode.Name = lexer.eat(tok.IDENTIFIER).Value
	lexer.eat(tok.EQUAL)

	newExpr, err := lexer.expression([]tok.TokenType{tok.SEMICOLON})
	if err != nil {
		panic(err)
	}
	newNode.Body = append(newNode.Body, newExpr)

	newNode.End = lexer.eat(tok.SEMICOLON).End

	return
}

func (lexer *Lexer) typeSpecifier() (returnNode Node, returnErr error) {
	newNode := Node{}

	defer func() {
		returnNode = newNode
		returnErr = nil
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnErr = e // paniced error를 반환
			}
		}
	}()

	if lexer.next(1).Type == tok.IDENTIFIER &&
		lexer.next(2).Type == tok.EQUAL {
		newDec, err := lexer.variableDeclaration()
		if err != nil {
			panic(err)
		}
		newNode = newDec
		return
	}

	panic(&perr.TokenError{
		Begin: lexer.cur().Begin,
		End:   lexer.cur().End,
		Msg:   "알 수 없는 구조의 코드",
	})
}

func (lexer *Lexer) Lexicalize() (Node, error) {

	for lexer.idx < lexer.length && lexer.cur().Type != tok.EOF {
		switch lexer.cur().Type {
		case tok.TYPE:
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
