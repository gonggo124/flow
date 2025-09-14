package anonyfunc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	utils "github.com/gonggo124/objective-mcfunction/Utils"
)

var anonyfunc_list = make([]string, 0)

func New(code string) string {
	anonyfunc_list = append(anonyfunc_list, code)
	index := len(anonyfunc_list) - 1
	return fmt.Sprintf("%s:anony/%d", os.Getenv("INTERNAL_NAMESPACE"), index)
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
