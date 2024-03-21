package slice

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func NoPreAlloc(size int) {
	data := make([]int, 0)
	for i := 0; i < size; i++ {
		data = append(data, i)
	}
}

func PreAlloc(size int) {
	data := make([]int, 0, size)
	for i := 0; i < size; i++ {
		data = append(data, i)
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

func generateWithCap(n int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, r.Int())
	}
	return nums
}

func printMem(t *testing.T) {
	t.Helper()
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	t.Logf("%.2f MB", float64(rtm.Alloc)/1024./1024.)
}

func testLastChars(t *testing.T, f func([]int) []int) {
	t.Helper()
	ans := make([][]int, 0)
	for k := 0; k < 100; k++ {
		origin := generateWithCap(128 * 1024) // 1M
		ans = append(ans, f(origin))
		runtime.GC() // 显性垃圾回收
	}
	printMem(t)
	_ = ans
}

func GetLastBySlice(origin []int) []int {
	return origin[len(origin)-2:]
}

func GetLastByCopy(origin []int) []int {
	ret := make([]int, 2)
	copy(ret, origin[len(origin)-2:])
	return ret
}

func TestLastCharsBySlice(t *testing.T) {
	testLastChars(t, GetLastBySlice)
}

func TestLastCharsByCopy(t *testing.T) {
	testLastChars(t, GetLastByCopy)
}
