//go:build (!amd64 && !arm64) || purego || race

package helpers

func IndexOfAnyInRange(in []rune, first, last rune) int {
	return indexRuneRangeScalar(in, first, last)
}

func IndexByteInRange(in string, first, last byte) int {
	return indexByteRangeScalar(in, first, last)
}
