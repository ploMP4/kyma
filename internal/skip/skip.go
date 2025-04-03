/// https://github.com/muesli/reflow/pull/75

package skip

import (
	"bytes"
	"io"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
)

// Writer drops some number of leading columns from a line of text, while
// leaving any ansi sequences intact.
type Writer struct {
	width  uint
	prefix string

	ansiWriter ansi.Writer
	buf        bytes.Buffer
	ansi       bool
}

// NewWriter returns a new writer that drops the given number of columns from a
// line of text, using the given prefix. Any visible content after the skipped
// portion will be preceded by the given prefix. The prefix is often used to
// provide some visual indication of when content has been scrolled.
func NewWriter(width uint, prefix string) *Writer {
	w := &Writer{
		width:  width,
		prefix: prefix,
	}
	w.ansiWriter.Forward = &w.buf
	return w
}

// NewWriterPipe returns a new writer that forwards the result to the given
// writer instead of its internal buffer.
func NewWriterPipe(forward io.Writer, width uint, prefix string) *Writer {
	return &Writer{
		width:  width,
		prefix: prefix,
		ansiWriter: ansi.Writer{
			Forward: forward,
		},
	}
}

// Bytes drops the specified number of printed columns from the given byte
// slice, leaving any ansi sequences intact.
func Bytes(b []byte, width uint) []byte {
	return BytesWithPrefix(b, width, nil)
}

// BytesWithPrefix drops the specified number of printed columns from the given
// byte slice, leaving any any sequences intact. Any visible content after the
// skipped portion will be preceded by the given prefix. The prefix is often
// used to provide some visual indication of when content has been scrolled.
func BytesWithPrefix(b []byte, width uint, prefix []byte) []byte {
	w := NewWriter(width, string(prefix))
	_, _ = w.Write(b)

	return w.Bytes()
}

// String drops the specified number of printed columns from the given string,
// leaving any ansi sequences intact.
func String(s string, width uint) string {
	return StringWithPrefix(s, width, "")
}

// StringWithPrefix drops the specified number of printed columns from the
// given string, leaving any ansi sequences intact. Any visible content after
// the skipped portion will be preceded by the given prefix. The prefix is
// often used to provide some visual indication of when content has been
// scrolled.
func StringWithPrefix(s string, width uint, prefix string) string {
	w := NewWriter(width, prefix)
	_, _ = w.Write([]byte(s))

	return w.String()
}

func (w *Writer) Write(b []byte) (int, error) {
	width := w.width
	if width > 0 {
		width += uint(ansi.PrintableRuneWidth(w.prefix))
	}

	var currentWidth uint
	for _, r := range string(b) {
		if r == ansi.Marker {
			// ANSI escape sequence
			w.ansi = true
		} else if w.ansi {
			if ansi.IsTerminator(r) {
				w.ansi = false
			}
		} else if currentWidth < width {
			rw := uint(runewidth.RuneWidth(r))
			if len(w.prefix) > 0 && currentWidth+rw >= width {
				_, err := w.ansiWriter.Write([]byte(w.prefix))
				if err != nil {
					return 0, err
				}
			}

			if currentWidth+rw > width {
				// double-width rune across the skip boundary.
				// Add spaces to preserve alignment.
				for currentWidth < width {
					_, _ = w.ansiWriter.Write([]byte(" "))
					currentWidth++
				}
			}
			currentWidth += rw

			continue
		}

		_, err := w.ansiWriter.Write([]byte(string(r)))
		if err != nil {
			return 0, err
		}
	}

	return len(b), nil
}

// Bytes returns the result as a byte slice.
func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}

// String returns the result as a string
func (w *Writer) String() string {
	return w.buf.String()
}
