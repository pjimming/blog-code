package string

import (
	"bytes"
	"strings"
	"testing"
)

func Plus(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s += str
	}
	return s
}

func StrBuilder(n int, str string) string {
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(str)
	}
	return builder.String()
}

func ByteBuffer(n int, str string) string {
	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		buf.WriteString(str)
	}
	return buf.String()
}

func PreStrBuilder(n int, str string) string {
	var builder strings.Builder
	builder.Grow(n * len(str))
	for i := 0; i < n; i++ {
		builder.WriteString(str)
	}
	return builder.String()
}

func PreByteBuffer(n int, str string) string {
	buf := new(bytes.Buffer)
	buf.Grow(n * len(str))
	for i := 0; i < n; i++ {
		buf.WriteString(str)
	}
	return buf.String()
}

func BenchmarkPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Plus(100, "abc")
	}
}

func BenchmarkStrBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StrBuilder(100, "abc")
	}
}

func BenchmarkByteBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ByteBuffer(100, "abc")
	}
}

func BenchmarkPreStrBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PreStrBuilder(100, "abc")
	}
}

func BenchmarkPreByteBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PreByteBuffer(100, "abc")
	}
}
