// 유틸 함수 모음
package utils

import (
	"strconv"
	"strings"
)

// 귀찮은 에러체크는 이제 안녕~
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// 파일 확장자 구하기
func FileEx(name string) string {
	namearr := strings.Split(name, ".")
	return namearr[len(namearr)-1]
}

// 파일 확장자 빼고 이름만
func FileName(name string) string {
	namearr := strings.Split(name, ".")
	return namearr[0]
}

var ct uint // 카운터 변수
func Count() uint {
	current := ct
	ct++
	return current
}

var panic_line int = -1
var current_file_path string

func SetLine(line int) {
	panic_line = line
}

func SetCurrentFile(curpath string) {
	current_file_path = curpath
}

func GetLine() int {
	return panic_line
}

func GetCurrentFile() string {
	return current_file_path
}

func Panic(msg any) {
	switch val := msg.(type) {
	case error:
		panic(val.Error() + " at " + current_file_path + ":" + strconv.Itoa(panic_line))
	case string:
		panic(val + " at " + current_file_path + ":" + strconv.Itoa(panic_line))
	default:
		panic("unsupported type in utils.Panic()")
	}
}
