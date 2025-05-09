/// https://github.com/muesli/reflow/pull/75

package skip_test

import (
	"bytes"
	"testing"

	"github.com/ploMP4/kyma/internal/skip"
)

func TestNewWriter(t *testing.T) {
	for _, tt := range tests() {
		t.Run(tt.name, func(t *testing.T) {
			f := skip.NewWriter(tt.width, tt.prefix)

			_, err := f.Write([]byte(tt.give))
			if err != nil {
				t.Error(err)
			}

			got := f.String()
			if got != tt.want {
				t.Errorf("Expected:\n\n`%s`\n\nActual Output:\n\n`%s`", tt.want, got)
			}
		})
	}
}

func TestNewWriterPipe(t *testing.T) {
	for _, tt := range tests() {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := skip.NewWriterPipe(&buf, tt.width, tt.prefix)

			_, err := f.Write([]byte(tt.give))
			if err != nil {
				t.Error(err)
			}

			if f.String() != "" {
				t.Errorf("Expected w.String() to return empty string, got `%s`", f.String())
			}
			if len(f.Bytes()) > 0 {
				t.Errorf("Expected w.Bytes() to return empty slice, got `%v`", f.Bytes())
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("Expected:\n\n`%s`\n\nActual Output:\n\n`%s`", tt.want, got)
			}
		})
	}
}

func TestString(t *testing.T) {
	for _, tt := range nonPrefixTests() {
		t.Run(tt.name, func(t *testing.T) {
			got := skip.String(tt.give, tt.width)
			if got != tt.want {
				t.Errorf("Expected:\n\n`%s`\n\nActual Output:\n\n`%s`", tt.want, got)
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			skip.String("\x1B[38;2;249;38;114mhello你好\x1B[0m", 5)
		}
	})
}

func TestStringWithPrefix(t *testing.T) {
	for _, tt := range tests() {
		t.Run(tt.name, func(t *testing.T) {
			got := skip.StringWithPrefix(tt.give, tt.width, tt.prefix)
			if got != tt.want {
				t.Errorf("Expected:\n\n`%s`\n\nActual Output:\n\n`%s`", tt.want, got)
			}
		})
	}
}

func BenchmarkStringWithPrefix(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			skip.StringWithPrefix("\x1B[38;2;249;38;114mhello你好\x1B[0m", 5, "…")
		}
	})
}

func TestBytes(t *testing.T) {
	for _, tt := range nonPrefixTests() {
		t.Run(tt.name, func(t *testing.T) {
			got := string(skip.Bytes([]byte(tt.give), tt.width))
			if got != tt.want {
				t.Errorf("Expected:\n\n`%s`\n\nActual Output:\n\n`%s`", tt.want, got)
			}
		})
	}
}

func TestBytesWithPrefix(t *testing.T) {
	for _, tt := range tests() {
		t.Run(tt.name, func(t *testing.T) {
			got := string(skip.BytesWithPrefix([]byte(tt.give), tt.width, []byte(tt.prefix)))
			if got != tt.want {
				t.Errorf("Expected:\n\n`%s`\n\nActual Output:\n\n`%s`", tt.want, got)
			}
		})
	}
}

type test struct {
	name   string
	width  uint
	prefix string
	give   string
	want   string
}

func tests() []test {
	return []test{
		{
			name:  "no-op",
			width: 0,
			give:  "foo",
			want:  "foo",
		},
		{
			name:   "no-op with prefix",
			width:  0,
			prefix: "…",
			give:   "foo",
			want:   "foo",
		},
		{
			name:  "no-op with ansi",
			width: 0,
			give:  "\x1B[7mfoo",
			want:  "\x1B[7mfoo",
		},
		{
			name:   "no-op with ansi and prefix",
			width:  0,
			prefix: "…",
			give:   "\x1B[7mfoo",
			want:   "\x1B[7mfoo",
		},
		{
			name:  "basic skip",
			width: 3,
			give:  "foobar",
			want:  "bar",
		},
		{
			name:   "basic skip with prefix",
			width:  3,
			prefix: "…",
			give:   "foobar",
			want:   "…ar",
		},
		{
			// corner case: prefix is honored even if it hides the only
			// remaining visible rune
			name:   "width minus 1 with prefix",
			width:  5,
			prefix: "…",
			give:   "foobar",
			want:   "…",
		},
		{
			name:  "same width",
			width: 3,
			give:  "foo",
			want:  "",
		},
		{
			name:   "same width with prefix",
			width:  3,
			prefix: "…",
			give:   "foo",
			want:   "",
		},
		{
			name:  "spaces only",
			width: 2,
			give:  "    ",
			want:  "  ",
		},
		{
			name:   "spaces only with prefix",
			width:  2,
			prefix: "…",
			give:   "    ",
			want:   "… ",
		},
		{
			name:  "double-width runes",
			width: 7,
			give:  "hello你好",
			want:  "好",
		},
		{
			name:  "double-width rune chopped and replaced by space",
			width: 6,
			give:  "hello你好",
			want:  " 好",
		},
		{
			name:   "double-width rune chopped and replaced by prefix",
			width:  6,
			prefix: "…",
			give:   "hello你好",
			want:   "…好",
		},
		{
			name:   "double-width rune replaced by prefix and space",
			width:  5,
			prefix: "…",
			give:   "hello你好",
			want:   "… 好",
		},
		{
			name:  "double-width rune chopped and replaced by space, with ansi",
			width: 6,
			give:  "\x1B[38;2;249;38;114mhello你好\x1B[0m",
			want:  "\x1B[38;2;249;38;114m 好\x1B[0m",
		},
		{
			name:   "double-width rune chopped and replaced by prefix, with ansi",
			width:  6,
			prefix: "…",
			give:   "\x1B[38;2;249;38;114mhello你好\x1B[0m",
			want:   "\x1B[38;2;249;38;114m…好\x1B[0m",
		},
		{
			name:   "double-width rune replaced by prefix and space, with ansi",
			width:  5,
			prefix: "…",
			give:   "\x1B[38;2;249;38;114mhello你好\x1B[0m",
			want:   "\x1B[38;2;249;38;114m… 好\x1B[0m",
		},
	}
}

func nonPrefixTests() []test {
	var ts []test
	for _, tt := range tests() {
		if tt.prefix == "" {
			ts = append(ts, tt)
		}
	}
	return ts
}
