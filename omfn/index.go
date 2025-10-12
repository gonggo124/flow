package omfn

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	utils "flow/Utils"
	lexer "flow/omfn/Lexer"
	perr "flow/omfn/ParseErrors"
	tokenizer "flow/omfn/Tokenizer"
	typechecker "flow/omfn/TypeChecker"
)

func printTree(ast lexer.Node, depth int) {
	fmt.Print(strings.Repeat("    ", depth))
	fmt.Printf("[%d]", ast.Type)
	fmt.Print(":" + ast.DataType)
	fmt.Printf(" %d~%d\n", ast.Begin, ast.End)

	if ast.Name != "" {
		fmt.Print(strings.Repeat("    ", depth))
		fmt.Print("    Name: '", ast.Name, "'\n")
	}
	if ast.Value != "" {
		fmt.Print(strings.Repeat("    ", depth))
		fmt.Print("    Value: '", ast.Value, "'\n")
	}
	for _, node := range ast.Body {
		printTree(node, depth+1)
	}
}

func removeLineComments(src string) string {
	var result []string

	reader := strings.NewReader(src)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") {
			// 줄 전체가 주석이면 무시
			continue
		}
		result = append(result, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return strings.Join(result, "\n")
}

func Parse(target_path string) {

	fmt.Println(target_path)
	data := utils.Must(os.ReadFile(target_path))
	code := removeLineComments(string(data))

	tokens := tokenizer.Tokenize(code)
	for _, tok := range tokens {
		fmt.Println(tok)
	}

	newLexer := lexer.New(tokens)
	ast, err := newLexer.Lexicalize()
	if err != nil {
		switch typedErr := err.(type) {
		case *perr.TokenError:
			panic(fmt.Sprintf("%s at %d~%d", typedErr, typedErr.Begin, typedErr.End))
		default:
			panic(typedErr)
		}
	}

	err = typechecker.Annotate(&ast)
	if err != nil {
		switch typedErr := err.(type) {
		case *perr.TokenError:
			panic(fmt.Sprintf("%s at %d~%d", typedErr, typedErr.Begin, typedErr.End))
		default:
			panic(typedErr)
		}
	}

	printTree(ast, 0)
	// err = astparser.Parse(ast)
	// if err != nil {
	// 	switch typedErr := err.(type) {
	// 	case *perr.TokenError:
	// 		panic(fmt.Sprintf("%s at %d~%d", typedErr, typedErr.Begin, typedErr.End))
	// 	default:
	// 		panic(typedErr)
	// 	}
	// }
}
