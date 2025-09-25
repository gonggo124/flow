package cmdvh

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	utils "flow/Utils"
	fnreader "flow/omfn/Reader"
)

// 이름 → 함수 매핑
var cmdMap = make(map[string]func(*fnreader.Reader, string, int) string)

func init() {
	cmdMap["v1_21_5_execute"] = v1_21_5_execute
	cmdMap["v1_21_5_bind"] = v1_21_5_bind
	cmdMap["v1_21_5_do"] = v1_21_5_do
	cmdMap["v1_21_5_fn"] = v1_21_5_fn
	cmdMap["v1_21_5_return"] = v1_21_5_return
	cmdMap["v1_21_5_while"] = v1_21_5_while
}

var version = "v1_21_5"

// 버전별 함수 리스트 반환
func SetVersion(v string) {
	version = "v" + strings.ReplaceAll(v, ".", "_")
}

type numberRange struct {
	begin int
	end   int
}

func handleRepeater(scanner *fnreader.Reader, line string, line_num int) string {
	numbers := make(map[string]numberRange)
	idx := 0
	length := len(line)

	rLine := []rune(line)

	inner_contents := ""
	idx++ // 맨 앞 '<' 무시
	for rLine[idx] != '>' {
		if idx >= length {
			utils.Panic("'<' 후 '>'를 찾을 수 없습니다")
		}
		inner_contents += string(rLine[idx])
		idx++
	}

	ranges := strings.Split(inner_contents, ",")
	for _, r := range ranges {
		key_val := strings.SplitN(r, "=", 2)
		val1, err := strconv.Atoi(strings.Split(key_val[1], "..")[0])
		if err != nil {
			utils.Panic("범위는 정수 단위 숫자여야 합니다")
		}
		val2, err := strconv.Atoi(strings.Split(key_val[1], "..")[1])
		if err != nil {
			utils.Panic("범위는 정수 단위 숫자여야 합니다")
		}
		numbers[key_val[0]] = numberRange{val1, val2}
	}

	idx++

	// original
	oCmd := strings.TrimSpace(string(rLine[idx:]))

	// target
	tCmd := ""

	if strings.HasPrefix(oCmd, "{") {
		// rune 다른데서 쓰지 마셈. 아래 루프에서 다른 값으로 바꿀 수도 있으니까
		rCmd := []rune(oCmd)

		length := len(rCmd)
		idx = 0
		brace_stack := 0
		escaped := false

		var icantnameit strings.Builder

		for {
			if length != 0 {
				if rCmd[idx] == '{' && !escaped {
					brace_stack += 1
				} else if rCmd[idx] == '}' && !escaped {
					brace_stack -= 1
					if brace_stack < 0 {
						utils.Panic("비정삭적인 중괄호 블록 발견")
					}
				} else if rCmd[idx] == '\\' && !escaped {
					escaped = true
				}
			}
			icantnameit.WriteRune(rCmd[idx])
			idx++
			if brace_stack == 0 {
				tCmd += icantnameit.String()
				break
			}
			if idx >= length && scanner.Scan() {
				tCmd += strings.TrimSpace(icantnameit.String()) + "\n"
				icantnameit.Reset()
				nLine := scanner.Text()
				utils.SetLine(nLine.Number)
				rCmd = []rune(nLine.Text)
				length = len(rCmd)
				idx = 0
				continue
			}
			if idx >= length {
				utils.Panic("중괄호가 닫히지 않았습니다")
				break
			}
		}
	} else {
		tCmd = oCmd
	}

	tCmd = strings.TrimPrefix(tCmd, "{")
	tCmd = strings.TrimSuffix(tCmd, "}")
	tCmd = strings.TrimSpace(tCmd)

	// new
	var newCmds strings.Builder
	maxx := make(map[string]bool)
	for len(maxx) < len(numbers) {
		// newCmd
		nC := tCmd
		for name, num_range := range numbers {
			nC = strings.ReplaceAll(
				strings.ReplaceAll(nC, "%"+name+"%", strconv.Itoa(num_range.begin)),
				"\\%", "%")
			if num_range.begin < num_range.end {
				num_range.begin += 1
			} else {
				maxx[name] = true
			}
			numbers[name] = num_range
		}
		newCmds.WriteString(nC + "\n")
	}

	var newCompiledCmd strings.Builder
	lScanner := bufio.NewScanner(strings.NewReader(newCmds.String()))
	var lReader fnreader.Reader
	lReader.SetScanner(lScanner)
	fmt.Println(len(lReader.Lines), lReader.Line)

	for lReader.Scan() {
		lReaderLine := lReader.Text()
		utils.SetLine(lReaderLine.Number%len(lReader.Lines) + line_num)
		newCompiledCmd.WriteString(
			ParseCmd(&lReader, lReaderLine) +
				"\n")
	}

	return newCompiledCmd.String()
}

func ParseCmd(scanner *fnreader.Reader, line *fnreader.Line) string {
	if line.Text == "" {
		return ""
	}
	text := line.Text
	last, _ := utf8.DecodeLastRuneInString(line.Text)
	if last == '\\' {
		text = strings.TrimSuffix(text, "\\")
		for last == '\\' && scanner.Scan() {
			curr := scanner.Text().Text
			text += strings.TrimSuffix(strings.TrimSpace(curr), "\\")
			last, _ = utf8.DecodeLastRuneInString(curr)
		}
	}
	if strings.HasPrefix(text, "<") {
		// Repeater
		return handleRepeater(scanner, text, line.Number)
	}
	keywords := strings.Split(text, " ")
	cmd, ok := cmdMap[version+"_"+keywords[0]]
	if !ok {
		fmt.Println("text:", text)
		return text
	}
	result := cmd(scanner, strings.Join(keywords[1:], " "), line.Number)
	fmt.Println(result)
	return result
}
