package yamlvector

import "math"

var eos [math.MaxUint8]bool

func init() {
	eos['\n'] = true
	eos['\r'] = true
}

func eol(src []byte, offset int) int {
	n := len(src)
	for i := offset; i < n; i++ {
		if eos[src[i]] {
			return i
		}
	}
	return n
}

func eot(src []byte, offset int) int {
	n := len(src)
	for i := offset; i < n; i++ {
		// todo check end-of-text block
	}
	return n
}
