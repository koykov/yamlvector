package yamlvector

import "unsafe"

func bin(src []byte, offset, size int) uint64 {
	n := len(src)
	if offset+size >= n {
		return 0
	}
	_ = src[n-1]
	binSrc := src[offset : offset+size]
	b := *(*uint64)(unsafe.Pointer(&binSrc[0]))
	return b
}
