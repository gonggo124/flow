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

func Parse(ast lexer.Node) (returnErr error) {
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
	readTree(ast)
	return
}

type Variable struct {
	Type       string
	Offset     string
	IsRegister bool
}

func reverseSlice(s []lexer.Node) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func readTree(ast lexer.Node) {

	if ast.Type == lexer.FUNCTION_DEFINITION {
		filename := filepath.Join(os.Getenv("MAIN_PATH"), "function", ast.Name) + ".mcfunction"
		new_file := utils.Must(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm))
		defer new_file.Close()

		writer := bufio.NewWriter(new_file)

		endFunc := func() {
			writer.WriteString("execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function " + os.Getenv("INT_NS") + ":mem/stack/stackptr/attach\n")
			writer.WriteString("execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run function " + os.Getenv("INT_NS") + ":mem/stack/ret\n")
			writer.WriteString("execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function " + os.Getenv("INT_NS") + ":mem/stack/cut\n")
		}

		writer.WriteString("#> COMPILED BY FLOW\n")

		// writer.WriteString("# @Return: " + ast.Body[0].Value + "\n")

		varOffsetTable := make(map[string]Variable)
		varOffset := 1
		getOffset := func() string {
			cur := varOffset
			varOffset++
			return strconv.Itoa(cur)
		}

		loadId := func(target lexer.Node, tt string) error {
			tVar, ok := varOffsetTable[target.Name]
			if !ok {
				return errors.New("존재하지 않는 변수입니다")
			}

			operand := "on passengers "
			suffix := "if entity @s[tag=_flow_internal.stack.bit,type=marker] "
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

			return nil
		}

		for idx, param := range ast.Body[1].Body {
			varOffsetTable[param.Name] = Variable{Type: param.Value, Offset: strconv.Itoa(-(idx + 1))}
		}

		writer.WriteString("execute as 6a56ec26-fbbd-4b1c-a7bf-59d89fd54460 on vehicle run tag @s add _flow_internal.stack.old_baseptr\n")

		writer.WriteString("scoreboard players set #sa0 _flow_internal.register 0\n")
		writer.WriteString("function " + os.Getenv("INT_NS") + ":mem/stack/push\n")

		writer.WriteString("execute as de8d7920-b907-4853-b3a2-c73cb0d5a84d on vehicle run function " + os.Getenv("INT_NS") + ":mem/stack/baseptr/attach\n")

		last_return := false
		for _, stuff := range ast.Body[2].Body {
			last_return = false
			if stuff.Type == lexer.CALL_EXPRESSION {
				reverseSlice(stuff.Body[0].Body)
				for _, arg := range stuff.Body[0].Body {
					if arg.Type == lexer.NUMBER {
						writer.WriteString("scoreboard players set #sa0 _flow_internal.register " + arg.Value + "\n")
						writer.WriteString("function " + os.Getenv("INT_NS") + ":mem/stack/push\n")
					}
					if arg.Type == lexer.IDENTIFIER {
						err := loadId(arg, "#sa0 _flow_internal.register")
						if err != nil {
							panic(&perr.TokenError{
								Begin: stuff.Begin,
								End:   stuff.End,
								Msg:   "존재하지 않는 변수입니다",
							})
						}
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
			}
			if stuff.Type == lexer.VARIABLE_DECLARATION {
				if stuff.Value != "int" {
					panic(&perr.TokenError{
						Begin: stuff.Begin,
						End:   stuff.End,
						Msg:   "변수의 타입은 'int'만 지원합니다",
					})
				}
				writer.WriteString("scoreboard players set #sa0 _flow_internal.register " + stuff.Body[0].Value + "\n")
				writer.WriteString("function " + os.Getenv("INT_NS") + ":mem/stack/push\n")
				varOffsetTable[stuff.Name] = Variable{Type: stuff.Value, Offset: getOffset()}
			}
			if stuff.Type == lexer.RETURN_EXPRESSION {
				last_return = true

				if len(stuff.Body) == 0 {
					endFunc()
					writer.WriteString("return 1\n")
					continue
				}
				if stuff.Body[0].Type == lexer.NUMBER {
					endFunc()
					writer.WriteString("scoreboard players set #return _flow_internal.register " + stuff.Body[0].Value + "\n")
					writer.WriteString("return 1\n")
				}
				if stuff.Body[0].Type == lexer.IDENTIFIER {
					err := loadId(stuff.Body[0], "#return _flow_internal.register")
					if err != nil {
						panic(&perr.TokenError{
							Begin: stuff.Begin,
							End:   stuff.End,
							Msg:   "존재하지 않는 변수입니다",
						})
					}
					endFunc()
					writer.WriteString("return 1\n")
				}

			}
		}

		if !last_return {
			endFunc()
		}

		writer.Flush() // 버퍼 비우기
		return
	}
	for _, node := range ast.Body {
		readTree(node)
	}
}
