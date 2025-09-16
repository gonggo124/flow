package cmdvh

import (
	"fmt"
	"strings"
	"unicode/utf8"

	fnreader "flux/omfn/Reader"
)

// 이름 → 함수 매핑
var cmdMap = make(map[string]func(*fnreader.Reader, string, int) string)

func init() {
	cmdMap["v1_21_5_execute"] = v1_21_5_execute
	cmdMap["v1_21_5_bind"] = v1_21_5_bind
	cmdMap["v1_21_5_do"] = v1_21_5_do
	cmdMap["v1_21_5_fn"] = v1_21_5_fn
}

var version = "v1_21_5"

// 버전별 함수 리스트 반환
func SetVersion(v string) {
	version = "v" + strings.ReplaceAll(v, ".", "_")
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
