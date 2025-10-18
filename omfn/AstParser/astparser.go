package astparser

import (
	"bufio"
	"errors"
	utils "flow/Utils"
	lexer "flow/omfn/Lexer"
	perr "flow/omfn/ParseErrors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Variable struct {
	Type       string
	Offset     string
	IsRegister bool
}

type Parser struct {
	varOffsetTable map[string]Variable
	varOffset      int
}

func reverseSlice(s []lexer.Node) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (p *Parser) Parse(ast lexer.Node) (returnErr error) {
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
	p.readTree(&ast)
	return
}

func (p *Parser) getOffset() string {
	cur := p.varOffset
	p.varOffset++
	return strconv.Itoa(cur)
}

func (p *Parser) loadId(target lexer.Node, tt string) (string, error) {
	var writer strings.Builder
	tVar, ok := p.varOffsetTable[target.Name]
	if !ok {
		return "", errors.New("존재하지 않는 변수입니다")
	}

	operand := "on passengers "
	suffix := "if entity @s[tag=_flow_internal.stack.bit,type=interaction] "
	offset, err := strconv.Atoi(tVar.Offset)
	if err != nil {
		panic(err)
	}
	if offset < 0 {
		offset = -offset
		operand = "on vehicle "
		suffix = ""
	}
	writer.WriteString("execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle ")
	writer.WriteString(strings.Repeat(operand, offset))
	writer.WriteString(suffix)
	writer.WriteString("run scoreboard players operation " + tt + " = @s _flow_internal.stack")
	writer.WriteString("\n")

	return writer.String(), nil
}

func (p *Parser) endFunc() string {
	var writer strings.Builder
	writer.WriteString("execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function " + os.Getenv("INT_NS") + ":mem/stack/stackptr/attach\n")
	writer.WriteString("execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function " + os.Getenv("INT_NS") + ":mem/stack/ret\n")
	writer.WriteString("execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function " + os.Getenv("INT_NS") + ":mem/stack/cut\n")
	return writer.String()
}

func (p *Parser) callExprssion(stuff *lexer.Node) string {
	var writer strings.Builder
	reverseSlice(stuff.Body[0].Body)
	for _, arg := range stuff.Body[0].Body {
		if arg.Type == lexer.NUMBER {
			writer.WriteString("scoreboard players set #sa0 _flow_internal.register " + arg.Value + "\n")
			writer.WriteString("function " + os.Getenv("INT_NS") + ":mem/stack/push\n")
		}
		if arg.Type == lexer.IDENTIFIER {
			result, err := p.loadId(arg, "#sa0 _flow_internal.register")
			if err != nil {
				panic(&perr.TokenError{
					Begin: stuff.Begin,
					End:   stuff.End,
					Msg:   "존재하지 않는 변수입니다",
				})
			}
			writer.WriteString(result)
			writer.WriteString("function " + os.Getenv("INT_NS") + ":mem/stack/push\n")
		}
	}
	writer.WriteString("function " + os.Getenv("MAIN_NS") + ":" + stuff.Name + "\n")
	if len(stuff.Body[0].Body) != 0 {
		writer.WriteString("execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle ")
		writer.WriteString(strings.Repeat("on vehicle ", len(stuff.Body[0].Body)))
		writer.WriteString("run function " + os.Getenv("INT_NS") + ":mem/stack/stackptr/attach\n")
		writer.WriteString("function " + os.Getenv("INT_NS") + ":mem/stack/cut\n")
	}
	return writer.String()
}

func (p *Parser) binaryExpression(stuff *lexer.Node, cnt int, resultR string) string {
	var writer strings.Builder

	cntReg := func() string {
		current := cnt
		cnt++
		return "#r" + strconv.Itoa(current)
	}

	operand1 := stuff.Body[0]
	operator := stuff.Body[1]
	operand2 := stuff.Body[2]

	op1Reg := cntReg()
	op2Reg := cntReg()

	switch operand1.Type {
	case lexer.BINARY_EXPRESSION:
		writer.WriteString(p.binaryExpression(&operand1, cnt, op1Reg))
	case lexer.NUMBER:
		writer.WriteString("scoreboard players set " + op1Reg + " _flow_internal.register " + operand1.Value + "\n")
	case lexer.IDENTIFIER:
		loadedResult, err := p.loadId(operand1,op1Reg+" _flow_internal.register")
		if err != nil {
			panic(&perr.TokenError{
				Begin: operand1.Begin,
				End:   operand1.End,
				Msg:   "존재하지 않는 변수입니다",
			})
		}
		writer.WriteString(loadedResult)
	}
	switch operand2.Type {
	case lexer.BINARY_EXPRESSION:
		writer.WriteString(p.binaryExpression(&operand2, cnt, op2Reg))
	case lexer.NUMBER:
		writer.WriteString("scoreboard players set " + op2Reg + " _flow_internal.register " + operand2.Value + "\n")
	case lexer.IDENTIFIER:
		loadedResult, err := p.loadId(operand2,op2Reg+" _flow_internal.register")
		if err != nil {
			panic(&perr.TokenError{
				Begin: operand2.Begin,
				End:   operand2.End,
				Msg:   "존재하지 않는 변수입니다",
			})
		}
		writer.WriteString(loadedResult)
	}

	if resultR == "" {
		writer.WriteString("scoreboard players operation " + op1Reg + " _flow_internal.register " + operator.Value + "= " + op2Reg + " _flow_internal.register\n")
	} else {
		writer.WriteString("scoreboard players operation " + resultR + " _flow_internal.register = " + op1Reg + " _flow_internal.register\n")
		writer.WriteString("scoreboard players operation " + resultR + " _flow_internal.register " + operator.Value + "= " + op2Reg + " _flow_internal.register\n")
	}

	return writer.String()
}

func (p *Parser) variableDeclaration(stuff *lexer.Node) string {
	var writer strings.Builder
	if stuff.Value != "int" && stuff.Value != "selector" {
		panic(&perr.TokenError{
			Begin: stuff.Begin,
			End:   stuff.End,
			Msg:   "변수의 타입은 'int'와 'selector'만 지원합니다",
		})
	}
	if stuff.Body[0].Type == lexer.BINARY_EXPRESSION {
		writer.WriteString(p.binaryExpression(&stuff.Body[0], 0, "")) // #r0
		writer.WriteString("scoreboard players operation #sa0 _flow_internal.register = #r0 _flow_internal.register\n")
	} else {
		writer.WriteString("scoreboard players set #sa0 _flow_internal.register " + stuff.Body[0].Value + "\n")
	}
	writer.WriteString("function " + os.Getenv("INT_NS") + ":mem/stack/push\n")
	p.varOffsetTable[stuff.Name] = Variable{Type: stuff.Value, Offset: p.getOffset()}
	return writer.String()
}

func (p *Parser) variableAssignment(stuff *lexer.Node) string {
	var writer strings.Builder
	tVar, ok := p.varOffsetTable[stuff.Name]
	if !ok {
		panic(&perr.TokenError{
			Begin: stuff.Begin,
			End:   stuff.End,
			Msg:   "변수를 찾을 수 없습니다",
		})
	}

	operand := "on passengers "
	suffix := "if entity @s[tag=_flow_internal.stack.bit,type=interaction] "
	offset, err := strconv.Atoi(tVar.Offset)
	if err != nil {
		panic(err)
	}
	if offset < 0 {
		offset = -offset
		operand = "on vehicle "
		suffix = ""
	}

	if stuff.Body[0].Type == lexer.BINARY_EXPRESSION {
		writer.WriteString(p.binaryExpression(&stuff.Body[0], 0, "")) // #r0
	} else {
		writer.WriteString("scoreboard players set #r0 _flow_internal.register " + stuff.Body[0].Value + "\n")
	}

	writer.WriteString("execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle ")
	writer.WriteString(strings.Repeat(operand, offset))
	writer.WriteString(suffix)
	writer.WriteString("run scoreboard players operation @s _flow_internal.stack = #r0 _flow_internal.register")
	writer.WriteString("\n")

	return writer.String()
}

func (p *Parser) returnExprssion(stuff *lexer.Node) string {

	var writer strings.Builder

	if len(stuff.Body) == 0 {
		writer.WriteString(p.endFunc())
		writer.WriteString("return 1\n")
		return writer.String()
	}
	if stuff.Body[0].Type == lexer.NUMBER {
		writer.WriteString(p.endFunc())
		writer.WriteString("scoreboard players set #return _flow_internal.register " + stuff.Body[0].Value + "\n")
		writer.WriteString("return 1\n")
	}
	if stuff.Body[0].Type == lexer.IDENTIFIER {
		result, err := p.loadId(stuff.Body[0], "#return _flow_internal.register")
		if err != nil {
			panic(&perr.TokenError{
				Begin: stuff.Begin,
				End:   stuff.End,
				Msg:   "존재하지 않는 변수입니다",
			})
		}
		writer.WriteString(result)
		writer.WriteString(p.endFunc())
		writer.WriteString("return 1\n")
	}

	return writer.String()
}

func (p *Parser) functionDefinition(ast *lexer.Node) {
	filename := filepath.Join(os.Getenv("MAIN_PATH"), "function", ast.Name) + ".mcfunction"
	new_file := utils.Must(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm))
	defer new_file.Close()

	writer := bufio.NewWriter(new_file)

	writer.WriteString("#> COMPILED BY FLOW\n")

	p.varOffsetTable = make(map[string]Variable)
	p.varOffset = 1

	for idx, param := range ast.Body[0].Body {
		p.varOffsetTable[param.Name] = Variable{Type: param.Value, Offset: strconv.Itoa(-(idx + 1))}
	}

	writer.WriteString("execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run tag @s add _flow_internal.stack.old_baseptr\n")

	writer.WriteString("scoreboard players set #sa0 _flow_internal.register 0\n")
	writer.WriteString("function " + os.Getenv("INT_NS") + ":mem/stack/push\n")

	writer.WriteString("execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function " + os.Getenv("INT_NS") + ":mem/stack/baseptr/attach\n")

	last_return := false
	for _, stuff := range ast.Body[1].Body {
		last_return = false
		if stuff.Type == lexer.CALL_EXPRESSION {
			writer.WriteString(p.callExprssion(&stuff))
		}
		if stuff.Type == lexer.VARIABLE_DECLARATION {
			writer.WriteString(p.variableDeclaration(&stuff))
		}
		if stuff.Type == lexer.ASSIGNMENT_EXPRESSION {
			writer.WriteString(p.variableAssignment(&stuff))
		}
		if stuff.Type == lexer.RETURN_EXPRESSION {
			writer.WriteString(p.returnExprssion(&stuff))
		}
		if stuff.Type == lexer.RAWLINE {
			writer.WriteString(stuff.Value+"\n")
		}
	}

	if !last_return {
		writer.WriteString(p.endFunc())
	}

	writer.Flush() // 버퍼 비우기
}

func (p *Parser) readTree(ast *lexer.Node) {

	if ast.Type == lexer.FUNCTION_DEFINITION {
		p.functionDefinition(ast)
		p.varOffsetTable = make(map[string]Variable)
		p.varOffset = 1
	}
	for i := range ast.Body {

		p.readTree(&ast.Body[i])
	}
}

func New() *Parser {
	newParser := Parser{}
	return &newParser
}
