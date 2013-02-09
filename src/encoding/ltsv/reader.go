// based on encoding/csv
package ltsv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// A Reader reads records from a LTSV-encoded file.
//
// As returned by NewReader, a Reader expects input LTSV-encoded file.
// The exported fields can be changed to customize the details before the
// first call to Read or ReadAll.
//
// Delimiter is the field delimiter.  It defaults to '\t'.
//
// Comment, if not 0, is the comment character. Lines beginning with the
// Comment character are ignored.
type Reader struct {
	Delimiter rune // Field delimiter (set to '\t' by NewReader)
	Comment   rune // Comment character for start of line
	line      int
	r         *bufio.Reader
	label     bytes.Buffer
	field     bytes.Buffer
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Delimiter: '\t',
		r:         bufio.NewReader(r),
	}
}

// Read reads one record from r.  The record is a slice of strings with each
// string representing one field.
func (r *Reader) Read() (map[string]string, error) {
	r.line++
	record := make(map[string]string)
	for {
		// If we are support comments and it is the comment character
		// then skip to the end of line.
		r1, err := r.readRune()
		if r.Comment != 0 && r1 == r.Comment {
			for {
				r1, err := r.readRune()
				if err != nil {
					return nil, err
				} else if r1 == '\n' {
					break
				}
			}
			continue
		}
		r.r.UnreadRune()

		label, end, err := r.parseLabel()
		if err != nil {
			return nil, err
		}
		if label == "" {
			if end {
				if len(record) != 0 {
					return record, nil
				}
			}
			continue // skip empty label
		}

		field, end, err := r.parseField()
		if err != nil {
			return nil, err
		}

		record[label] = field
		if end {
			return record, nil
		}
	}
	panic("unreachable")
}

func (r *Reader) parseLabel() (string, bool, error) {
	r.label.Reset()
	for {
		r1, err := r.readRune()
		if err != nil {
			return "", false, err
		} else if r1 == ':' {
			return strings.TrimSpace(r.label.String()), false, nil
		} else if r1 == '\n' {
			return "", true, nil
		} else if r1 == '\t' {
			return "", false, nil // no label
		} else if unicode.IsControl(r1) || !unicode.IsPrint(r1) {
			return "", false, errors.New(fmt.Sprintf("line %d: invalid rune at label", r.line))
		}
		r.label.WriteRune(r1)
	}
	panic("unreachable")
}

func (r *Reader) parseField() (string, bool, error) {
	r.field.Reset()
	for {
		r1, err := r.readRune()
		if err != nil {
			if err == io.EOF {
				return r.field.String(), true, nil
			}
			return "", false, err
		} else if r1 == '\t' {
			return r.field.String(), false, nil
		} else if r1 == '\n' {
			return r.field.String(), true, nil
		}
		r.field.WriteRune(r1)
	}
	panic("unreachable")
}

func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.r.ReadRune()
	if r1 == '\r' {
		r1, _, err = r.r.ReadRune()
		if err == nil {
			if r1 != '\n' {
				r.r.UnreadRune()
				r1 = '\r'
			}
		}
	}
	return r1, err
}

func (r *Reader) ReadAll() ([]map[string]string, error) {
	records := make([]map[string]string, 0)
	for {
		record, err := r.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	panic("unreachable")
}
