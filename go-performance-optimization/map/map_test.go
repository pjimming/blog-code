package _map

import "testing"

func NoPreAlloc(size int) {
	data := make(map[int]int)
	for i := 0; i < size; i++ {
		data[i] = i
	}
}

func PreAlloc(size int) {
	data := make(map[int]int, size)
	for i := 0; i < size; i++ {
		data[i] = i
	}
}

func BenchmarkNoPreAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NoPreAlloc(1000)
	}
}

func BenchmarkPreAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PreAlloc(1000)
	}
}
