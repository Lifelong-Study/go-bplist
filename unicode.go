package go_bplist

import (
	"encoding/binary"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

const (
	replacementChar = '\uFFFD'     // Unicode replacement character
	maxRune         = '\U0010FFFF' // Maximum valid Unicode code point.
)

const (
	// 0xd800-0xdc00 encodes the high 10 bits of a pair.
	// 0xdc00-0xe000 encodes the low 10 bits of a pair.
	// the value is those 20 bits plus 0x10000.
	surr1 = 0xd800
	surr2 = 0xdc00
	surr3 = 0xe000

	// surrSelf = 0x10000
)

func UInt16ArrayToString(data []byte) string {
	return DecodeUTF16ToString(UInt8ToUInt16(data))
}

func UInt8ToUInt16(data []byte) []uint16 {
	out := make([]uint16, 0)

	for i := 0; i < len(data); i += 2 {
		buf := make([]byte, 0)
		buf = append(buf, data[i+0])
		buf = append(buf, data[i+1])
		out = append(out, binary.BigEndian.Uint16(buf))
	}

	return out
}

func DecodeUTF16ToString(s []uint16) string {
	n := 0
	for i := 0; i < len(s); i++ {
		switch r := s[i]; {
		case r < surr1, surr3 <= r:
			// normal rune
			n += utf8.RuneLen(rune(r))
		case surr1 <= r && r < surr2 && i+1 < len(s) &&
			surr2 <= s[i+1] && s[i+1] < surr3:
			// valid surrogate sequence
			n += utf8.RuneLen(utf16.DecodeRune(rune(r), rune(s[i+1])))
			i++
		default:
			// invalid surrogate sequence
			n += utf8.RuneLen(replacementChar)
		}
	}
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < len(s); i++ {
		switch r := s[i]; {
		case r < surr1, surr3 <= r:
			// normal rune
			b.WriteRune(rune(r))
		case surr1 <= r && r < surr2 && i+1 < len(s) &&
			surr2 <= s[i+1] && s[i+1] < surr3:
			// valid surrogate sequence
			b.WriteRune(utf16.DecodeRune(rune(r), rune(s[i+1])))
			i++
		default:
			// invalid surrogate sequence
			b.WriteRune(replacementChar)
		}
	}
	return b.String()
}
