package helpers

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRangeSearch(t *testing.T) {
	rng := rand.New(rand.NewSource(73))
	bounds := [][2]rune{{'a', 'q'}, {0, 0}, {255, 255}, {-1, 1}, {-2147483648, 2147483647}, {2147483647, 2147483647}, {-2147483648, -2147483648}, {10, -10}, {128, 255}}
	for n := 0; n <= 257; n++ {
		bytes := make([]byte, n+16)
		runes := make([]rune, n+16)
		for i := range bytes {
			bytes[i] = byte(rng.Uint32())
			runes[i] = rune(rng.Uint32())
		}

		for offset := 0; offset < 16; offset++ {
			for _, bound := range bounds {
				input := runes[offset : offset+n]
				if got, want := IndexOfAnyInRange(input, bound[0], bound[1]), indexRuneRangeScalar(input, bound[0], bound[1]); got != want {
					t.Fatalf("rune n=%d offset=%d bounds=%v got=%d want=%d", n, offset, bound, got, want)
				}

				s := string(bytes)[offset : offset+n]
				if got, want := IndexByteInRange(s, byte(bound[0]), byte(bound[1])), indexByteRangeScalar(s, byte(bound[0]), byte(bound[1])); got != want {
					t.Fatalf("byte n=%d offset=%d bounds=%v got=%d want=%d", n, offset, bound, got, want)
				}
			}
		}
	}
}

func TestRangeEveryPosition(t *testing.T) {
	for n := 1; n <= 129; n++ {
		for pos := 0; pos < n; pos++ {
			data := make([]byte, n)
			runes := make([]rune, n)
			data[pos] = 200
			runes[pos] = -100
			if got := IndexByteInRange(string(data), 128, 255); got != pos {
				t.Fatalf("byte n=%d pos=%d got=%d", n, pos, got)
			}

			if got := IndexOfAnyInRange(runes, -101, -99); got != pos {
				t.Fatalf("rune n=%d pos=%d got=%d", n, pos, got)
			}
		}
	}
}

func FuzzRangeSearch(f *testing.F) {
	f.Add([]byte{0, 255, 128, 1, 0, 0, 0, 128}, int32(-1), int32(255))
	f.Add(make([]byte, 129), int32(0), int32(0))
	f.Fuzz(func(t *testing.T, data []byte, first, last int32) {
		s := string(data)
		if got, want := IndexByteInRange(s, byte(first), byte(last)), indexByteRangeScalar(s, byte(first), byte(last)); got != want {
			t.Fatalf("byte got=%d want=%d", got, want)
		}

		r := make([]rune, len(data)/4)
		for i := range r {
			j := i * 4
			r[i] = rune(uint32(data[j]) | uint32(data[j+1])<<8 | uint32(data[j+2])<<16 | uint32(data[j+3])<<24)
		}

		if got, want := IndexOfAnyInRange(r, first, last), indexRuneRangeScalar(r, first, last); got != want {
			t.Fatalf("rune got=%d want=%d", got, want)
		}
	})
}

func TestRangeChunkBoundaries(t *testing.T) {
	for _, n := range []int{16383, 16384, 16385, 65535, 65536, 65537, 131073} {
		for _, pos := range []int{-1, n - 1, n / 2} {
			data := make([]byte, n)
			runes := make([]rune, n)
			if pos >= 0 {
				data[pos] = 1
				runes[pos] = 1
			}

			if got := IndexByteInRange(string(data), 1, 1); got != pos {
				t.Fatalf("byte n=%d pos=%d got=%d", n, pos, got)
			}

			if got := IndexOfAnyInRange(runes, 1, 1); got != pos {
				t.Fatalf("rune n=%d pos=%d got=%d", n, pos, got)
			}
		}
	}
}

func TestRangeConcurrent(t *testing.T) {
	data := make([]rune, 65537)
	data[len(data)-1] = 1
	s := string(data)
	for worker := 0; worker < 16; worker++ {
		t.Run(fmt.Sprint(worker), func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 100; i++ {
				if got := IndexOfAnyInRange(data, 1, 1); got != len(data)-1 {
					t.Fatal(got)
				}

				if got := IndexByteInRange(s, 1, 1); got != len(s)-1 {
					t.Fatal(got)
				}
			}
		})
	}
}

func TestRangeNoAllocations(t *testing.T) {
	data := make([]rune, 131073)
	s := string(data)
	allocs := testing.AllocsPerRun(100, func() {
		if IndexOfAnyInRange(data, 1, 1) != -1 || IndexByteInRange(s, 1, 1) != -1 {
			t.Fatal("unexpected match")
		}
	})
	if allocs != 0 {
		t.Fatalf("allocations=%v", allocs)
	}
}
