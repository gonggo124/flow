package omfn

import (
	"bufio"
	"os"

	utils "flow/Utils"
	cmdvh "flow/VersionHandler"
	fnreader "flow/omfn/Reader"
)

func Parse(target_path string) string {
	file := utils.Must(os.Open(target_path))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var reader fnreader.Reader
	reader.SetScanner(scanner)

	utils.SetCurrentFile(target_path)

	result := ""

	for reader.Scan() {
		line := reader.Text()
		result += cmdvh.ParseCmd(&reader, line) + "\n"
	}

	return result
}
