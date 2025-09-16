package omce

import (
	"fmt"
	"os"
	"strings"

	ut "flux/Utils"
	lexer "flux/omce/Lexer"
	node "flux/omce/Node"
	tokenizer "flux/omce/Tokenizer"
)

type MethodSet struct {
	Name      string
	Namespace string
	Includes  []string
	Methods   map[string]string
}

// type CustomEntity struct {
// 	Name         string
// 	Extender     string
// 	Nbt          map[string]any
// 	Init_command string
// 	Init_params  []string
// 	Methods      map[string]string
// }

func parseTree(ms *MethodSet, ast node.Node) {
	switch ast.Type {
	case node.INCLUDE_DECLARATION:
		for _, include := range ast.Body {
			ms.Includes = append(ms.Includes, include.Name)
		}
	case node.METHOD_DECLARATION:
		ms.Methods[ast.Id.Name] = ast.Body[0].Name
	}
	for _, node := range ast.Body {
		parseTree(ms, node)
	}
}

func printTree(ast node.Node, depth int) {
	fmt.Print(strings.Repeat("    ", depth))
	fmt.Print("[" + ast.String() + "]\n")
	if ast.Name != "" {
		fmt.Print(strings.Repeat("    ", depth))
		fmt.Print("    Name: '", ast.Name, "'\n")
	}
	if ast.Id != nil {
		fmt.Print(strings.Repeat("    ", depth))
		fmt.Print("    Id: '", ast.Id.Name, "'\n")
	}
	if ast.Code != "" {
		fmt.Print(strings.Repeat("    ", depth))
		fmt.Print("    Code: '", strings.ReplaceAll(ast.Code, "\n", "\\n"), "'\n")
	}
	if len(ast.Nbt) != 0 {
		fmt.Print(strings.Repeat("    ", depth))
		fmt.Print("    Nbt: map[string]any ...\n")
	}
	for _, node := range ast.Body {
		printTree(node, depth+1)
	}
}

func Parse(omce_file_path string, omce_namespace string, omce_name string) MethodSet {
	ut.SetCurrentFile(omce_file_path)
	ut.SetLine(0)

	code := string(ut.Must(os.ReadFile(omce_file_path)))
	var new_method_set MethodSet
	new_method_set.Name = omce_name
	new_method_set.Namespace = omce_namespace
	new_method_set.Methods = make(map[string]string)

	tokens := tokenizer.Tokenize(code)

	// for _, tok := range tokens {
	// 	fmt.Println(tok)
	// }

	ast := lexer.Lexicalize(tokens)
	// printTree(ast, 0)
	parseTree(&new_method_set, ast)

	return new_method_set
}
