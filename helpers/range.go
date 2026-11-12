package helpers

func indexRuneRangeScalar(in []rune, first, last rune) int {
	for i, c := range in {
		if c >= first && c <= last {
			return i
		}
	}

	return -1
}

func indexByteRangeScalar(in string, first, last byte) int {
	for i := 0; i < len(in); i++ {
		if in[i] >= first && in[i] <= last {
			return i
		}
	}

	return -1
}
