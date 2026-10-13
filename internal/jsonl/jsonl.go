package jsonl

import (
	"bufio"
	"io"
)

// ReadLines reads r line by line and calls fn for each line.
//
// num is the 1-based line number. line holds the raw bytes of the line,
// including its line ending ("\n" or "\r\n"); the last line may have no
// line ending. Lines are not limited in length.
//
// If fn returns an error, ReadLines stops and returns that error.
// Otherwise it returns any error from reading r, or nil at the end of input.
func ReadLines(r io.Reader, fn func(num int, line []byte) error) error {
	n := 0
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadBytes('\n')
		if len(line) > 0 {
			n++
			if err := fn(n, line); err != nil {
				return err
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
