package syntax

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const meta = `\.+*?()|[]{}^$# `

func Escape(input string) string {
	for i, r := range input {
		if unicode.IsPrint(r) && !strings.ContainsRune(meta, r) && (r != utf8.RuneError || strings.HasPrefix(input[i:], "\uFFFD")) {
			continue
		}

		var buf strings.Builder
		buf.Grow(len(input))
		buf.WriteString(input[:i])
		for _, ch := range input[i:] {
			escape(&buf, ch, false)
		}

		return buf.String()
	}

	return input
}

func escape(b *strings.Builder, r rune, force bool) {
	if unicode.IsPrint(r) {
		if strings.ContainsRune(meta, r) || force {
			b.WriteRune('\\')
		}
		b.WriteRune(r)
		return
	}

	switch r {
	case '\a':
		// \a is an identity escape in non-Unicode ECMAScript patterns.
		b.WriteString(`\x07`)
	case '\f':
		b.WriteString(`\f`)
	case '\n':
		b.WriteString(`\n`)
	case '\r':
		b.WriteString(`\r`)
	case '\t':
		b.WriteString(`\t`)
	case '\v':
		b.WriteString(`\v`)
	default:
		if r < 0x100 {
			b.WriteString(`\x`)
			s := strconv.FormatInt(int64(r), 16)
			if len(s) == 1 {
				b.WriteRune('0')
			}
			b.WriteString(s)
			break
		}
		b.WriteString(`\u`)
		b.WriteString(strconv.FormatInt(int64(r), 16))
	}
}

func Unescape(input string) (string, error) {
	idx := strings.IndexRune(input, '\\')
	// no slashes means no unescape needed
	if idx == -1 {
		return input, nil
	}

	buf := bytes.NewBufferString(input[:idx])
	// get the runes for the rest of the string -- we're going full parser scan on this

	p := parser{}
	p.setPattern(input[idx+1:])
	for {
		if p.rightMost() {
			return "", p.getErr(ErrIllegalEndEscape)
		}
		r, err := p.scanCharEscape(false)
		if err != nil {
			return "", err
		}
		buf.WriteRune(r)
		// are we done?
		if p.rightMost() {
			return buf.String(), nil
		}

		r = p.moveRightGetChar()
		for r != '\\' {
			buf.WriteRune(r)
			if p.rightMost() {
				// we're done, no more slashes
				return buf.String(), nil
			}
			// keep scanning until we get another slash
			r = p.moveRightGetChar()
		}
	}
}
