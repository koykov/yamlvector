package yamlvector

import (
	"strconv"

	"github.com/koykov/byteconv"
)

// isNumeric reports whether b can be parsed as a number.
func isNumeric(s []byte) bool {
	if len(s) == 0 || !firstByte[s[0]] {
		return false
	}
	switch s[0] {
	case 'I', 'i', 'N', 'n':
		// fallthrough to ParseFloat
	default:
		for i := 1; i < len(s); i++ {
			if !bodyByte[s[i]] {
				return false
			}
		}
	}
	_, err := strconv.ParseFloat(byteconv.B2S(s), 64)
	return err == nil
}

var firstByte = func() (t [256]bool) {
	for _, c := range "+-.0123456789IiNn" {
		t[c] = true
	}
	return
}()

var bodyByte = func() (t [256]bool) {
	for _, c := range "0123456789eE+-." {
		t[c] = true
	}
	return
}()
