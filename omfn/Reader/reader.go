package fnreader

import (
	"bufio"
)

type Reader struct {
	Lines []string
	Line  int
}

type Line struct {
	Text   string
	Number int
}

func (r Reader) Scan() bool {
	return r.Line < len(r.Lines)
}

func (r *Reader) Text() *Line {
	new_line := Line{r.Lines[r.Line], r.Line + 1}
	r.Line += 1
	return &new_line
}

func (r *Reader) SetScanner(s *bufio.Scanner) {
	for s.Scan() {
		r.Lines = append(r.Lines, s.Text())
	}
}
