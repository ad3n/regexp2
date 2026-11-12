//go:build (amd64 || arm64) && !purego && !race

package helpers

import "unsafe"

func IndexOfAnyInRange(in []rune, first, last rune) int {
	for _, c := range in {
		if c >= first && c <= last {
			return 0
		}

		break
	}

	return indexRuneRange(in, first, last)
}

func indexRuneRange(in []rune, first, last rune) int {
	if len(in) < 4 {
		for i := 1; i < len(in); i++ {
			if in[i] >= first && in[i] <= last {
				return i
			}
		}

		return -1
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

func IndexByteInRange(in string, first, last byte) (index int) {
	for _, c := range []byte(in) {
		if c >= first && c <= last {
			return
		}

		break
	}

	return indexByteRange(in, first, last)
}

func indexByteRange(in string, first, last byte) int {
	if len(in) < 8 {
		for i := 1; i < len(in); i++ {
			if in[i] >= first && in[i] <= last {
				return i
			}
		}

		return -1
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
