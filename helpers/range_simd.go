//go:build (amd64 || arm64) && !purego && !race

package helpers

import "unsafe"

func IndexOfAnyInRange(in []rune, first, last rune) int {
	if len(in) > 0 && in[0] >= first && in[0] <= last {
		return 0
	}

	return indexRuneRange(in, first, last)
}

func indexRuneRange(in []rune, first, last rune) int {
	if len(in) < 8 {
		return indexRuneRangeScalar(in, first, last)
	}

	for offset := 0; offset < len(in); {
		n := min(len(in)-offset, 16384)
		if i := indexRuneRangeSIMD(&in[offset], n, first, last); i >= 0 {
			return offset + i
		}

		offset += n
	}

	return -1
}

func IndexByteInRange(in string, first, last byte) int {
	if len(in) > 0 && in[0] >= first && in[0] <= last {
		return 0
	}

	return indexByteRange(in, first, last)
}

func indexByteRange(in string, first, last byte) int {
	if len(in) < 16 {
		return indexByteRangeScalar(in, first, last)
	}

	for offset := 0; offset < len(in); {
		n := min(len(in)-offset, 65536)
		if i := indexByteRangeSIMD(unsafe.StringData(in[offset:]), n, uint32(first), uint32(last)); i >= 0 {
			return offset + i
		}

		offset += n
	}

	return -1
}

//go:noescape
func indexRuneRangeSIMD(in *rune, n int, first, last rune) int

//go:noescape
func indexByteRangeSIMD(in *byte, n int, first, last uint32) int
