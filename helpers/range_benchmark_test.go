package helpers

import (
	"fmt"
	"strings"
	"testing"
)

var rangeResult int

func BenchmarkRangeScan(b *testing.B) {
	for _, n := range []int{8, 16, 32, 256, 32768} {
		for _, pos := range []string{"miss", "first", "last"} {
			data := []byte(strings.Repeat("z", n))
			if pos == "first" {
				data[0] = 'b'
			}

			if pos == "last" {
				data[n-1] = 'b'
			}

			s := string(data)
			r := []rune(s)
			for _, impl := range []struct {
				name string
				f    func() int
			}{
				{"byte_simd", func() int { return IndexByteInRange(s, 'a', 'q') }},
				{"byte_scalar", func() int { return indexByteRangeScalar(s, 'a', 'q') }},
				{"rune_simd", func() int { return IndexOfAnyInRange(r, 'a', 'q') }},
				{"rune_scalar", func() int { return indexRuneRangeScalar(r, 'a', 'q') }},
			} {
				b.Run(fmt.Sprintf("%d/%s/%s", n, pos, impl.name), func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						rangeResult = impl.f()
					}
				})
			}
		}
	}
}
