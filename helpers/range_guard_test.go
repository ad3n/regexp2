//go:build darwin || linux

package helpers

import (
	"os"
	"syscall"
	"testing"
	"unsafe"
)

func TestRangeGuardPages(t *testing.T) {
	page := os.Getpagesize()
	memory, err := syscall.Mmap(-1, 0, 3*page, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := syscall.Munmap(memory); err != nil {
			t.Error(err)
		}
	}()

	for i := page; i < 2*page; i++ {
		memory[i] = 0x7f
	}

	for _, region := range [][]byte{memory[:page], memory[2*page:]} {
		if err := syscall.Mprotect(region, syscall.PROT_NONE); err != nil {
			t.Fatal(err)
		}
	}

	if err := syscall.Mprotect(memory[page:2*page], syscall.PROT_READ); err != nil {
		t.Fatal(err)
	}

	for n := 0; n <= 257; n++ {
		for _, start := range []int{page, 2*page - n} {
			s := unsafe.String(&memory[start], n)
			for _, bounds := range [][2]byte{{0, 126}, {127, 127}, {128, 255}, {255, 0}} {
				if got, want := IndexByteInRange(s, bounds[0], bounds[1]), indexByteRangeScalar(s, bounds[0], bounds[1]); got != want {
					t.Fatalf("byte n=%d got=%d want=%d", n, got, want)
				}
			}
		}

		for _, start := range []int{page, 2*page - 4*n} {
			r := unsafe.Slice((*rune)(unsafe.Pointer(&memory[start])), n)
			for _, bounds := range [][2]rune{{0, 126}, {0x7f7f7f7f, 0x7f7f7f7f}, {-2147483648, -1}, {255, 0}} {
				if got, want := IndexOfAnyInRange(r, bounds[0], bounds[1]), indexRuneRangeScalar(r, bounds[0], bounds[1]); got != want {
					t.Fatalf("rune n=%d got=%d want=%d", n, got, want)
				}
			}
		}
	}
}
