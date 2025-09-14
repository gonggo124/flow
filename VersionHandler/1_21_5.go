package cmdvh

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	utils "github.com/gonggo124/objective-mcfunction/Utils"
	anonyfunc "github.com/gonggo124/objective-mcfunction/omfn/Anonyfunc"
	fnreader "github.com/gonggo124/objective-mcfunction/omfn/Reader"
)

// 이제 안 씀 ㅅㄱ ㅋㅋㅋ
// var v1_25_5align_rule = regexp.MustCompile(`^[xyz]{1,3}$`)
// var v1_25_5selector_rule = regexp.MustCompile(`^@([arpnes])(?:\[([^,]*=(?:{.*}|[^,]*)(?:,[^,]*=(?:{.*}|[^,]*))*)?\])?$`)
// var v1_25_5position_slot_rule = regexp.MustCompile(`^[~\^]?(?:-?\d*\.?\d+)?$`)

// 이것도 안 씀
// func boolToInt(a bool) int {
// 	if a {
// 		return 1
// 	} else {
// 		return 0
// 	}
// }

// 이것도 안 쓸 거임 ㅅㄱ
// func handleSelector(selector string) string {
// 	matches := v1_25_5selector_rule.FindStringSubmatch(selector)
// 	if matches == nil {
// 		return "WHAT IS THIS"
// 	}
// 	selectorType := matches[1]
// 	args := matches[2]
// 	keys := make([]string, 0)
// 	vals := make([]string, 0)
// 	if args != "" {
// 		pairs := strings.Split(args, ",")
// 		for _, pair := range pairs {
// 			key_value := strings.SplitN(pair, "=", 2)
// 			key := strings.TrimSpace(key_value[0])
// 			val := strings.TrimSpace(key_value[1])
// 			keys = append(keys, key)
// 			vals = append(vals, val)
// 		}
// 		ASSEMBLE := ""
// 		for i := 0; i < len(pairs); i++ {
// 			if ASSEMBLE != "" {
// 				ASSEMBLE += ","
// 			}
// 			ASSEMBLE += keys[i] + "=" + vals[i]
// 		}
// 		return "@" + selectorType + "[" + ASSEMBLE + "]"
// 	} else {
// 		return "@" + selectorType
// 	}
// }

// func handleset(end rune) {

// }

// gpt 레전드네 ㄹㅇ 와;;
func v1_21_5_split_cmd(s string) []string {
	var tokens []string
	var cur strings.Builder

	bracketDepth := 0   // '[' ']'의 깊이 (selector 안의 콤마 무시하기 위해)
	inQuote := false    // 큰따옴표 내부 플래그
	escapeNext := false // 백슬래시 이스케이프 처리 (간단히 지원)

	for _, r := range s {
		if escapeNext {
			// 이스케이프된 문자는 무조건 현재 토큰에 추가
			cur.WriteRune(r)
			escapeNext = false
			continue
		}

		switch r {
		case '\\':
			// 다음 문자를 이스케이프 처리
			escapeNext = true
			// NOTE: 백슬래시 자체를 토큰에 넣고 싶지 않다면 생략
			// cur.WriteRune(r)
		case '"':
			// 따옴표 토글 (따옴표도 토큰에 포함)
			inQuote = !inQuote
			cur.WriteRune(r)
		case '[':
			bracketDepth++
			cur.WriteRune(r)
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
			cur.WriteRune(r)
		default:
			// 공백 처리: 공백이면서 대괄호/따옴표 안이 아니면 토큰 경계
			if (r == ' ' || r == '\t' || r == '\n') && bracketDepth == 0 && !inQuote {
				if cur.Len() > 0 {
					tokens = append(tokens, cur.String())
					cur.Reset()
				}
				// 연속 공백은 무시
			} else {
				cur.WriteRune(r)
			}
		}
	}

	// 마지막 토큰 푸시
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}

	return tokens
}

// func v1_21_5_title(scanner *fnreader.Reader, line string, line_num uint64) string {
// 	for _, c := range v1_21_5_split_cmd(line) {
// 		fmt.Println(c)
// 	}
// 	return "title " + line
// }

func v1_21_5_execute(scanner *fnreader.Reader, line string, line_num int) string {
	utils.SetLine(line_num)
	code := "execute "
	args := v1_21_5_split_cmd(line)

	// args_num := 0

	length := len(args)
	idx := 0

	push := func() {
		code += args[idx] + " "
		idx++
	}

	for idx < length {
		switch args[idx] {
		case "align":
			push()
			push()
			continue
		case "as":
			push()
			push()
			continue
		case "at":
			push()
			push()
			continue
		case "anchored":
			push()
			push()
			continue
		case "facing":
			push()
			push()
			push()
			push()
			continue
		case "if":
			push()
			switch args[idx] {
			case "biome":
				push()
				push()
				push()
				push()
				push()
			case "block":
				push()
				push()
				push()
				push()
				push()
			case "blocks":
				push()
				push()
				push()
				push()
				push()
				push()
				push()
				push()
				push()
				push()
				push()
			case "data":
				push()
				switch args[idx] {
				case "entity":
					push()
					push()
					push()
				case "block":
					push()
					push()
					push()
					push()
					push()
				case "storage":
					push()
					push()
					push()
				}
			case "dimension":
				push()
				push()
			case "entity":
				push()
				push()
			case "function":
				push()
				push()
			case "items":
				push()
				switch args[idx] {
				case "entity":
					push()
					push()
					push()
					push()
				case "block":
					push()
					push()
					push()
					push()
					push()
					push()
				}
			case "loaded":
				push()
				push()
				push()
				push()
			case "predicate":
				push()
				push()
			case "score":
				push()
				push()
				push()
				switch args[idx] {
				case "matches":
					push()
					push()
				default:
					push()
					push()
					push()
				}
			}
			continue
		case "in":
			push()
			push()
			continue
		case "on":
			push()
			push()
			continue
		case "positioned":
			push()
			if args[idx] == "as" || args[idx] == "over" {
				push()
				push()
				continue
			}
			push()
			push()
			push()
			continue
		case "rotated":
			push()
			if args[idx] == "as" {
				push()
				push()
				continue
			}
			push()
			push()
		case "run":
			push()
			// 아오;; 뻑ㄱㄱ

			cut_line := fnreader.Line{Text: strings.Join(args[idx:], " "), Number: line_num}

			code += ParseCmd(scanner, &cut_line)

			idx = length
			continue
		case "store":
			push()
			push()
			switch args[idx] {
			case "block":
				push()
				push()
				push()
				push()
				push()
				push()
				push()
			case "bossbar":
				push()
				push()
				push()
			case "entity", "storage":
				push()
				push()
				push()
				push()
				push()
			case "score":
				push()
				push()
				push()
			}
		case "summon":
			push()
			push()
		case "unless":
			push()
			switch args[idx] {
			case "biome":
				push()
				push()
				push()
				push()
				push()
			case "block":
				push()
				push()
				push()
				push()
				push()
			case "blocks":
				push()
				push()
				push()
				push()
				push()
				push()
				push()
				push()
				push()
				push()
				push()
			case "data":
				push()
				switch args[idx] {
				case "entity":
					push()
					push()
					push()
				case "block":
					push()
					push()
					push()
					push()
					push()
				case "storage":
					push()
					push()
					push()
				}
			case "dimension":
				push()
				push()
			case "entity":
				push()
				push()
			case "function":
				push()
				push()
			case "items":
				push()
				switch args[idx] {
				case "entity":
					push()
					push()
					push()
					push()
				case "block":
					push()
					push()
					push()
					push()
					push()
					push()
				}
			case "loaded":
				push()
				push()
				push()
				push()
			case "predicate":
				push()
				push()
			case "score":
				push()
				push()
				push()
				switch args[idx] {
				case "matches":
					push()
					push()
				default:
					push()
					push()
					push()
				}
			}
			continue
		default:
			fmt.Printf("line %d: 엥 이게 뭐람\n", line_num)
			continue
		}
	}
	return code
}

func v1_21_5_bind(scanner *fnreader.Reader, line string, line_num int) string {
	utils.SetLine(line_num)
	code := "execute as "

	args := v1_21_5_split_cmd(line)

	// length := len(args)
	idx := 0

	push := func() string {
		cur := args[idx]
		code += args[idx] + " "
		idx++
		return cur
	}

	if strings.HasPrefix(push(), "@s") {
		code = "function "
	} else {
		code += "at @s run function "
	}
	magic := strings.Split(args[idx], ":")
	code += os.Getenv("INTERNAL_NAMESPACE") + ":" + strings.ReplaceAll(filepath.Join("bind", magic[0], magic[1]), "\\", "/")

	return code
}

func v1_21_5_do(scanner *fnreader.Reader, line string, line_num int) string {
	utils.SetLine(line_num)
	code := "execute as "

	args := v1_21_5_split_cmd(line)

	// length := len(args)
	idx := 0

	if strings.HasPrefix(args[idx], "@") {
		// do @s.something style
		splited := strings.SplitN(args[idx], ".", 2)
		code += splited[0]
		code += " at @s run function " + os.Getenv("INTERNAL_NAMESPACE") + ":" + strings.ReplaceAll(filepath.Join("do", splited[1]), "\\", "/")
	} else {
		// do something style
		code = "function "
		code += os.Getenv("INTERNAL_NAMESPACE") + ":" + strings.ReplaceAll(filepath.Join("do", args[idx]), "\\", "/")
	}

	return code
}

func v1_21_5_fn(scanner *fnreader.Reader, line string, line_num int) string {
	utils.SetLine(line_num)
	// code := ""

	// inner_code := ""

	brace_stack := 0
	var inner_code strings.Builder

	args := v1_21_5_split_cmd(line)
	idx := 0

	contains := func(s string, sub_s rune) bool {
		escaped := false
		for _, chr := range s {
			if escaped {
				escaped = false
				continue
			}
			if chr == '\\' {
				escaped = true
			}
			if chr == sub_s {
				return true
			}
		}
		return false
	}

	length := len(args)
	for {
		// switch args[idx] {
		// case "{":
		// 	brace_stack += 1
		// 	inner_code.WriteString("{")
		// case "}":
		// 	brace_stack -= 1
		// 	inner_code.WriteString("}")
		// case "\\":
		// 	idx++
		// 	if idx < length && (args[idx] != "{" && args[idx] != "}") {
		// 		inner_code.WriteString("\\")
		// 	}
		// 	inner_code.WriteString(args[idx] + " ")
		// default:
		// 	inner_code.WriteString(args[idx] + " ")
		// }
		if strings.Contains(args[idx], "{") {
			brace_stack += 1
		}
		if contains(args[idx], '}') {
			brace_stack -= 1
		}
		inner_code.WriteString(args[idx] + " ")

		idx++
		if brace_stack <= 0 {
			break
		} else if idx >= length && scanner.Scan() {
			inner_code.WriteRune('\n')
			idx = 0
			current_line := scanner.Text()
			utils.SetLine(current_line.Number)
			args = v1_21_5_split_cmd(current_line.Text)
			length = len(args)
		} else if idx >= length {
			utils.Panic("중괄호가 닫히지 않았습니다")
		}
	}
	new_code := inner_code.String()
	new_code = strings.TrimSpace(new_code)
	new_code = strings.TrimPrefix(new_code, "{")
	new_code = strings.TrimSuffix(new_code, "}")
	new_code = strings.TrimSpace(new_code)

	var writer strings.Builder

	scanner2 := bufio.NewScanner(strings.NewReader(new_code))

	var reader fnreader.Reader
	reader.SetScanner(scanner2)

	for reader.Scan() {
		line := reader.Text()
		line.Text = strings.TrimSpace(line.Text)
		writer.WriteString(ParseCmd(&reader, line) + "\n")
	}

	// anonyfunc.New("")
	return "function " + anonyfunc.New(writer.String())
}
