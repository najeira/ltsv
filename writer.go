// based on encoding/csv
package ltsv

import (
	"bufio"
	"io"
)

// A Writer writes records to a LTSV encoded file.
//
// As returned by NewWriter, a Writer writes records terminated by a
// newline and uses '\t' as the field delimiter.  The exported fields can be
// changed to customize the details before the first call to Write or WriteAll.
//
// Delimiter is the field delimiter.
//
// If UseCRLF is true, the Writer ends each record with \r\n instead of \n.
type Writer struct {
	Delimiter rune // Label delimiter (set to to '\t' by NewWriter)
	UseCRLF   bool // True to use \r\n as the line terminator
	w         *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Delimiter: '\t',
		w:         bufio.NewWriter(w),
	}
}

// Writer writes a single CSV record to w along with any necessary quoting.
// A record is a slice of strings with each string being one field.
func (w *Writer) Write(record map[string]string) error {
	var err error
	num := 0
	for label, field := range record {
		if num >= 1 {
			if _, err = w.w.WriteRune('\t'); err != nil {
				return err
			}
		}
		if _, err = w.w.WriteString(label); err != nil {
			return err
		}
		if _, err = w.w.WriteRune(':'); err != nil {
			return err
		}
		if _, err = w.w.WriteString(field); err != nil {
			return err
		}
		num++
	}
	var line_end string
	if w.UseCRLF {
		line_end = "\r\n"
	} else {
		line_end = "\n"
	}
	if _, err = w.w.WriteString(line_end); err != nil {
		return err
	}
	return nil
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() {
	w.w.Flush()
}

// WriteAll writes multiple LTSV records to w using Write and then calls Flush.
func (w *Writer) WriteAll(records []map[string]string) error {
	var err error
	for _, record := range records {
		if err = w.Write(record); err != nil {
			break
		}
	}
	w.Flush()
	return err
}
