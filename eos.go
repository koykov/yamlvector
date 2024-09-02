package yamlvector

import "math"

var eos [math.MaxUint8]bool

func init() {
	eos['\n'] = true
	eos['\r'] = true
}

func eol(src []byte, offset int) int {
	n := len(src)
	for i := offset; i < len(src); i++ {
		if eos[src[i]] {
			return i
		}
	}
	return n
}
