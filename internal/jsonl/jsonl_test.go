package jsonl

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
)

func TestReadLines(t *testing.T) {
	long := strings.Repeat("a", 200000)
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty", input: "", want: nil},
		{name: "one line", input: "a\n", want: []string{"a\n"}},
		{name: "no final newline", input: "a\nb", want: []string{"a\n", "b"}},
		{name: "CRLF", input: "a\r\nb\r\n", want: []string{"a\r\n", "b\r\n"}},
		{name: "empty lines", input: "\n\n", want: []string{"\n", "\n"}},
		{name: "long line", input: long + "\n" + "short\n", want: []string{long + "\n", "short\n"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			err := ReadLines(strings.NewReader(tt.input), func(num int, line []byte) error {
				if num != len(got)+1 {
					t.Errorf("line number = %v, want %v", num, len(got)+1)
				}
				got = append(got, string(line))
				return nil
			})
			if err != nil {
				t.Errorf("error = %v, want nil", err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("ReadLines = %q, want %q", shortAll(got), shortAll(tt.want))
			}
		})
	}
}

func TestReadLinesStopsOnFnError(t *testing.T) {
	errStop := errors.New("stop")
	input := "a\nb\nc\n"

	calls := 0
	err := ReadLines(strings.NewReader(input), func(num int, line []byte) error {
		calls++
		if num == 2 {
			return errStop
		}
		return nil
	})
	if calls != 2 {
		t.Errorf("fn called %d times, want 2", calls)
	}
	if !errors.Is(err, errStop) {
		t.Errorf("ReadLines error = %v, want %v", err, errStop)
	}
}

func TestReadLinesReadError(t *testing.T) {
	errDisk := errors.New("disk read error")
	mr := io.MultiReader(strings.NewReader("a\nb"), iotest.ErrReader(errDisk))
	var got []string
	want := []string{"a\n", "b"}
	err := ReadLines(mr, func(num int, line []byte) error {
		if num != len(got)+1 {
			t.Errorf("line number = %v, want %v", num, len(got)+1)
		}
		got = append(got, string(line))
		return nil
	})
	if !errors.Is(err, errDisk) {
		t.Errorf("ReadLines error = %v, want %v", err, errDisk)
	}
	if !slices.Equal(got, want) {
		t.Errorf("ReadLines = %q, want %q", got, want)
	}
}

func short(s string) string {
	if len(s) <= 40 {
		return s
	}
	return fmt.Sprintf("%s...(%d bytes)", s[:20], len(s))
}

func shortAll(lines []string) []string {
	var out []string
	for _, s := range lines {
		out = append(out, short(s))
	}
	return out
}
