package _struct

import "testing"

func EmptyStructMap(n int) {
	m := make(map[int]struct{})
	for i := 0; i < n; i++ {
		m[i] = struct{}{}
	}
}

func BoolMap(n int) {
	m := make(map[int]bool)
	for i := 0; i < n; i++ {
		m[i] = false
	}
}

func BenchmarkEmptyStructMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EmptyStructMap(1000)
	}
}

func BenchmarkBoolMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BoolMap(1000)
	}
}
