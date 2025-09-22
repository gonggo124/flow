package anonyfunc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	utils "flow/Utils"
)

var anonyfunc_list = make([]string, 0)

// 등록
func New(code string) string {
	anonyfunc_list = append(anonyfunc_list, code)
	index := len(anonyfunc_list) - 1
	return fmt.Sprintf("%s:anony/%d", os.Getenv("INTERNAL_NAMESPACE"), index)
}

// 자리만 맡아둠
func Reserve() int {
	anonyfunc_list = append(anonyfunc_list, "")
	index := len(anonyfunc_list) - 1
	return index
}

func Update(idx int, code string) {
	anonyfunc_list[idx] = code
}

// func Load() []string {
// 	return anonyfunc_list
// }

func CreateFiles() {
	internal := filepath.Join(os.Getenv("INTERNAL_PATH"), "function", "anony")
	for index, anonycode := range anonyfunc_list {
		anonyfunc_file_path := filepath.Join(internal, fmt.Sprintf("%d.mcfunction", index))
		new_anonyfunc_file := utils.Must(os.Create(anonyfunc_file_path))
		defer new_anonyfunc_file.Close()
		writer := bufio.NewWriter(new_anonyfunc_file)

		writer.WriteString(anonycode)

		writer.Flush()
	}
}
